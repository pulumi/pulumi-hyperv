// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vm

import (
	"context"
	"fmt"

	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/memory"
	"github.com/microsoft/wmi/pkg/virtualization/core/processor"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance" // Updated import path
	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// The following statements are not required. They are type assertions to indicate to Go that Vm implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors at this location.

// They would normally be included in the vmController.go file, but they're located here for instructive purposes.
var _ = (infer.CustomResource[VmInputs, VmOutputs])((*Vm)(nil))
var _ = (infer.CustomUpdate[VmInputs, VmOutputs])((*Vm)(nil))
var _ = (infer.CustomDelete[VmOutputs])((*Vm)(nil))

func (c *Vm) Connect(ctx context.Context) (*vmms.VMMS, *service.VirtualSystemManagementService, error) {
	// Create the VMMS client.
	config := infer.GetConfig[common.Config](ctx)
	var whost *host.WmiHost
	if config.Host != "" {
		whost = host.NewWmiHost(config.Host)
	} else {
		whost = host.NewWmiLocalHost()
	}

	vmmsClient, err := vmms.NewVMMS(whost)
	if err != nil {
		return nil, nil, err
	}
	vsms := vmmsClient.GetVirtualSystemManagementService()
	return vmmsClient, vsms, nil
}

// This is the Get Metadata method.
func (c *Vm) Read(ctx context.Context, id string, inputs VmInputs, preview bool) (VmOutputs, error) {
	// This is a no-op. We don't need to do anything here.
	return VmOutputs{}, nil
}

// This is the Create method. This will be run on every Vm resource creation.
func (c *Vm) Create(ctx context.Context, name string, input VmInputs, preview bool) (string, VmOutputs, error) {
	logger := provider.GetLogger(ctx)
	state := VmOutputs{VmInputs: input}
	id, err := resource.NewUniqueHex(name, 8, 0)
	if err != nil {
		return id, state, err
	}

	// If in preview, don't run the command.
	if preview {
		return id, state, nil
	}
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		return id, state, err
	}
	setting, err := virtualsystem.GetVirtualSystemSettingData(vmmsClient.GetVirtualizationConn().WMIHost, *input.VmName)
	if err != nil {
		return id, state, err
	}
	defer setting.Close()
	logger.Debug("Create VMSettings")

	err = setting.SetHyperVGeneration(virtualsystem.HyperVGeneration_V2)
	if err != nil {
		return id, state, fmt.Errorf("Failed [%+v]", err)
	}

	memorySettings, err := memory.GetDefaultMemorySettingData(vmmsClient.GetVirtualizationConn().WMIHost)
	if err != nil {
		return id, state, fmt.Errorf("Failed [%+v]", err)
	}
	defer memorySettings.Close()
	var memorySizeMB uint64 = 1024 // Default value
	if input.MemorySize != nil {
		memorySizeMB = uint64(*input.MemorySize)
	}
	memorySettings.SetSizeMB(memorySizeMB)

	processorSettings, err := processor.GetDefaultProcessorSettingData(vmmsClient.GetVirtualizationConn().WMIHost)
	if err != nil {
		return id, state, fmt.Errorf("Failed [%+v]", err)
	}
	var cpuCount uint64 = 1 // Default value
	if input.ProcessorCount != nil {
		cpuCount = uint64(*input.ProcessorCount)
	}
	processorSettings.SetCPUCount(cpuCount)

	vm, err := vsms.CreateVirtualMachine(setting, memorySettings, processorSettings)
	if err != nil {
		return id, state, fmt.Errorf("Failed [%+v]", err)
	}
	logger.Debug("Created VM")
	defer func() {
		if vm != nil {
			vm.Close()
		}
	}()

	return id, state, nil
}

// WireDependencies controls how secrets and unknowns flow through a resource.
//
//	var _ = (infer.ExplicitDependencies[VmInputs, VmOutputs])((*Vm)(nil))
//	func (r *Vm) WireDependencies(f infer.FieldSelector, args *VmInputs, state *VmOutputs) { .. }
//
// Because we want every output to depend on every input, we can leave the default behavior.

// The Update method will be run on every update.
func (c *Vm) Update(ctx context.Context, id string, olds VmOutputs, news VmInputs, preview bool) (VmOutputs, error) {
	state := VmOutputs{VmInputs: news}
	// If in preview, don't run the command.
	if preview {
		return state, nil
	}
	// Use Create command if Update is unspecified.
	cmd := news.Create
	if news.Update != nil {
		cmd = news.Update
	}
	// If neither are specified, do nothing.
	if cmd == nil {
		return state, nil
	}
	return state, nil
}

// The Delete method will run when the resource is deleted.
func (c *Vm) Delete(ctx context.Context, id string, props VmOutputs) error {
	if props.Delete == nil {
		return nil
	}
	_, vsms, err := c.Connect(ctx)
	if err != nil {
		return err
	}

	vm, err := vsms.GetVirtualMachineByName(*props.VmName)
	if err != nil {
		return err
	}

	defer vm.Close()
	err = vm.Start()
	if err != nil {
		return fmt.Errorf("Failed [%+v]", err)
	}

	err = vm.Stop(true)
	if err != nil {
		return fmt.Errorf("Failed [%+v]", err)
	}
	err = vsms.DeleteVirtualMachine(vm)
	if err != nil {
		return fmt.Errorf("Failed [%+v]", err)
	}
	return nil
}

// VirtualMachineState represents the state of a virtual machine.
type VirtualMachineState uint16

const (
	// VirtualMachineStateUnknown indicates the state of the virtual machine could not be determined.
	VirtualMachineStateUnknown VirtualMachineState = 0
	// VirtualMachineStateOther indicates the virtual machine is in an other state.
	VirtualMachineStateOther VirtualMachineState = 1
	// VirtualMachineStateRunning indicates the virtual machine is running.
	VirtualMachineStateRunning VirtualMachineState = 2
	// VirtualMachineStateOff indicates the virtual machine is turned off.
	VirtualMachineStateOff VirtualMachineState = 3
	// VirtualMachineStateShuttingDown indicates the virtual machine is in the process of turning off.
	VirtualMachineStateShuttingDown VirtualMachineState = 4
	// VirtualMachineStateNotApplicable indicates the virtual machine does not support being started or turned off.
	VirtualMachineStateNotApplicable VirtualMachineState = 5
	// VirtualMachineStateEnabledButOffline indicates the virtual machine might be completing commands, and it will drop any new requests.
	VirtualMachineStateEnabledButOffline VirtualMachineState = 6
	// VirtualMachineStateInTest indicates the virtual machine is in a test state.
	VirtualMachineStateInTest VirtualMachineState = 7
	// VirtualMachineStateDeferred indicates the virtual machine might be completing commands, but it will queue any new requests.
	VirtualMachineStateDeferred VirtualMachineState = 8
	// VirtualMachineStateQuiesce indicates the virtual machine is running but in a restricted mode.
	// The behavior of the virtual machine is similar to the Running state, but it processes only a restricted set of commands.
	// All other requests are queued.
	VirtualMachineStateQuiesce VirtualMachineState = 9
	// VirtualMachineStateStarting indicates the virtual machine is in the process of starting. New requests are queued.
	VirtualMachineStateStarting VirtualMachineState = 10
)

// String returns the string representation of the VirtualMachineState.
func (s VirtualMachineState) String() string {
	switch s {
	case VirtualMachineStateUnknown:
		return "Unknown"
	case VirtualMachineStateOther:
		return "Other"
	case VirtualMachineStateRunning:
		return "Running"
	case VirtualMachineStateOff:
		return "Off"
	case VirtualMachineStateShuttingDown:
		return "ShuttingDown"
	case VirtualMachineStateNotApplicable:
		return "NotApplicable"
	case VirtualMachineStateEnabledButOffline:
		return "EnabledButOffline"
	case VirtualMachineStateInTest:
		return "InTest"
	case VirtualMachineStateDeferred:
		return "Deferred"
	case VirtualMachineStateQuiesce:
		return "Quiesce"
	case VirtualMachineStateStarting:
		return "Starting"
	default:
		return "Unknown"
	}
}

// CreateKeyProtector creates a key protector object.
func CreateKeyProtector(v *vmms.VMMS) (*wmi.WmiInstance, error) {
	return nil, fmt.Errorf("CreateKeyProtector not implemented")
	// return v.HgsConn().CreateInstance("MSFT_HgsKeyProtector", nil)
}

// DefineSystem defines a virtual machine system.
func DefineSystem(v *vmms.VMMS, systemSettings *wmi.WmiInstance, resourceSettings []*wmi.WmiInstance) (*wmi.WmiInstance, error) {
	return nil, fmt.Errorf("DefineSystem not implemented")
	// // Get the WMI text representation of the system settings
	// systemText, err := systemSettings.GetText()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get system settings text: %w", err)
	// }

	// // Convert resource settings to an array of strings
	// rsStrings := make([]string, len(resourceSettings))
	// for i, rs := range resourceSettings {
	// 	text, err := rs.GetText()
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get resource setting text: %w", err)
	// 	}
	// 	rsStrings[i] = text
	// }

	// params := map[string]interface{}{
	// 	"SystemSettings":   systemText,
	// 	"ResourceSettings": rsStrings,
	// }

	// result, err := v.VirtualMachineManagementService().InvokeMethod("DefineSystem", params)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to define system: %w", err)
	// }

	// if err := v.ValidateOutput(result); err != nil {
	// 	return nil, err
	// }

	// resultingSystemPath, err := result.GetString("ResultingSystem")
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get resulting system path: %w", err)
	// }

	// return v.VirtualizationConn().Get(resultingSystemPath)
}

// ExistsVirtualMachine checks if a virtual machine with the given name exists.
func ExistsVirtualMachine(v *vmms.VMMS, name string) (bool, error) {
	return false, fmt.Errorf("ExistsVirtualMachine not implemented")
	// query := fmt.Sprintf("SELECT * FROM Msvm_ComputerSystem WHERE Caption = 'Virtual Machine' AND ElementName = '%s'", name)
	// vms, err := v.VirtualizationConn().Query(query)
	// if err != nil {
	// 	return false, fmt.Errorf("failed to query virtual machines: %w", err)
	// }

	// return len(vms) > 0, nil
}

// ModifySystemSettings modifies system settings.
func ModifySystemSettings(v *vmms.VMMS, systemSettings *wmi.WmiInstance) error {
	return fmt.Errorf("ModifySystemSettings not implemented")
	// // Get the WMI text representation of the system settings
	// systemText, err := systemSettings.GetText()
	// if err != nil {
	// 	return fmt.Errorf("failed to get system settings text: %w", err)
	// }

	// params := map[string]interface{}{
	// 	"SystemSettings": systemText,
	// }

	// result, err := v.GetVirtualMachineManagementService().InvokeMethod("ModifySystemSettings", params)
	// if err != nil {
	// 	return fmt.Errorf("failed to modify system settings: %w", err)
	// }

	// return v.ValidateOutput(result)
}
