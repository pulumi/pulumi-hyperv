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
	"fmt"

	"github.com/microsoft/wmi"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
)

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
func CreateKeyProtector(v *vmms.VMMS) (*wmi.Result, error) {
	return v.HgsConn().CreateInstance("MSFT_HgsKeyProtector", nil)
}

// DefineSystem defines a virtual machine system.
func DefineSystem(v *vmms.VMMS, systemSettings *wmi.Result, resourceSettings []*wmi.Result) (*wmi.Result, error) {
	// Get the WMI text representation of the system settings
	systemText, err := systemSettings.GetText()
	if err != nil {
		return nil, fmt.Errorf("failed to get system settings text: %w", err)
	}

	// Convert resource settings to an array of strings
	rsStrings := make([]string, len(resourceSettings))
	for i, rs := range resourceSettings {
		text, err := rs.GetText()
		if err != nil {
			return nil, fmt.Errorf("failed to get resource setting text: %w", err)
		}
		rsStrings[i] = text
	}

	params := map[string]interface{}{
		"SystemSettings":   systemText,
		"ResourceSettings": rsStrings,
	}

	result, err := v.VirtualMachineManagementService().InvokeMethod("DefineSystem", params)
	if err != nil {
		return nil, fmt.Errorf("failed to define system: %w", err)
	}

	if err := v.ValidateOutput(result); err != nil {
		return nil, err
	}

	resultingSystemPath, err := result.GetString("ResultingSystem")
	if err != nil {
		return nil, fmt.Errorf("failed to get resulting system path: %w", err)
	}

	return v.VirtualizationConn().Get(resultingSystemPath)
}

// ExistsVirtualMachine checks if a virtual machine with the given name exists.
func ExistsVirtualMachine(v *vmms.VMMS, name string) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM Msvm_ComputerSystem WHERE Caption = 'Virtual Machine' AND ElementName = '%s'", name)
	vms, err := v.VirtualizationConn().Query(query)
	if err != nil {
		return false, fmt.Errorf("failed to query virtual machines: %w", err)
	}

	return len(vms) > 0, nil
}

// DestroyVirtualMachine destroys a virtual machine with the given name.
func DestroyVirtualMachine(v *vmms.VMMS, name string) error {
	if name == "" {
		return fmt.Errorf("virtual machine name cannot be empty")
	}

	vm, err := GetVirtualMachine(v, name)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"AffectedSystem": vm.Path(),
	}

	result, err := v.VirtualMachineManagementService().InvokeMethod("DestroySystem", params)
	if err != nil {
		return fmt.Errorf("failed to destroy virtual machine: %w", err)
	}

	return v.ValidateOutput(result)
}

// GetVirtualMachine gets a virtual machine by name.
func GetVirtualMachine(v *vmms.VMMS, name string) (*wmi.Result, error) {
	query := fmt.Sprintf("SELECT * FROM Msvm_ComputerSystem WHERE Caption = 'Virtual Machine' AND ElementName = '%s'", name)
	vms, err := v.VirtualizationConn().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query virtual machines: %w", err)
	}

	if len(vms) == 0 {
		return nil, fmt.Errorf("unable to find the Virtual Machine %s", name)
	}

	return vms[0], nil
}

// ModifySystemSettings modifies system settings.
func ModifySystemSettings(v *vmms.VMMS, systemSettings *wmi.Result) error {
	// Get the WMI text representation of the system settings
	systemText, err := systemSettings.GetText()
	if err != nil {
		return fmt.Errorf("failed to get system settings text: %w", err)
	}

	params := map[string]interface{}{
		"SystemSettings": systemText,
	}

	result, err := v.VirtualMachineManagementService().InvokeMethod("ModifySystemSettings", params)
	if err != nil {
		return fmt.Errorf("failed to modify system settings: %w", err)
	}

	return v.ValidateOutput(result)
}

// BaseInputs is the common set of inputs for all local commands.
type BaseInputs struct {
	VmName         *string `pulumi:"vmname,optional"`
	ProcessorCount *int    `pulumi:"processorCount,optional"`
	MemorySize     *int    `pulumi:"memorySize,optional"`
}

// Implementing Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *BaseInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.VmName, "Name of the Virtual Machine")
	a.Describe(&c.ProcessorCount, "Number of processors to allocate to the Virtual Machine. Defaults to 1.")
	a.Describe(&c.MemorySize, "Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.")
}

type BaseOutputs struct{}

// Implementing Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *BaseOutputs) Annotate(a infer.Annotator) {
}
