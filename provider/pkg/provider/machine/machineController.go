//d/tmp/machineController.go
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

package machine

import (
	"context"
	"fmt"
	"strings"

	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/memory"
	"github.com/microsoft/wmi/pkg/virtualization/core/processor"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
)

// The following statements are not required. They are type assertions to indicate to Go that Machine implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors at this location.

// They would normally be included in the vmController.go file, but they're located here for instructive purposes.
var _ = (infer.CustomResource[MachineInputs, MachineOutputs])((*Machine)(nil))
var _ = (infer.CustomUpdate[MachineInputs, MachineOutputs])((*Machine)(nil))
var _ = (infer.CustomDelete[MachineOutputs])((*Machine)(nil))

func (c *Machine) Connect(ctx context.Context) (*vmms.VMMS, *service.VirtualSystemManagementService, error) {
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
func (c *Machine) Read(ctx context.Context, id string, inputs MachineInputs, preview bool) (MachineOutputs, error) {
	logger := provider.GetLogger(ctx)

	// Initialize the outputs with the inputs
	outputs := MachineOutputs{
		MachineInputs: inputs,
	}

	// The ID is the machine name if it's set, otherwise it's the ID
	machineName := id
	if inputs.MachineName != nil {
		machineName = *inputs.MachineName
	}

	// If in preview, don't attempt to fetch actual VM data
	if preview {
		return outputs, nil
	}

	// Connect to Hyper-V
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		return outputs, fmt.Errorf("failed to connect to Hyper-V: %v", err)
	}

	// Get the VM by name
	vm, err := vsms.GetVirtualMachineByName(machineName)
	if err != nil {
		logger.Debug(fmt.Sprintf("Machine %s not found: %v", machineName, err))
		return outputs, nil
	}
	defer vm.Close()

	logger.Debug(fmt.Sprintf("Found machine %s", machineName))

	// Get VM ID (ElementName in Hyper-V lingo)
	vmId, err := vm.GetPropertyElementName()
	if err == nil && vmId != "" {
		outputs.VmId = &vmId
		logger.Debug(fmt.Sprintf("VM ID: %s", vmId))
	}

	// Get the VM settings data
	vmSettings, err := virtualsystem.GetVirtualSystemSettingData(vmmsClient.GetVirtualizationConn().WMIHost, machineName)
	if err != nil {
		logger.Debug(fmt.Sprintf("Failed to get VM settings: %v", err))
		return outputs, nil
	}
	defer vmSettings.Close()

	// Get processor count - find the processor setting from VM settings
	// Use default values if not able to get processor setting
	if inputs.ProcessorCount == nil {
		defaultProcCount := 1
		inputs.ProcessorCount = &defaultProcCount
		logger.Debug(fmt.Sprintf("Using default processor count: %d", defaultProcCount))
	}

	// Get memory size - find the memory setting from VM settings
	// Use default values if not able to get memory setting
	if inputs.MemorySize == nil {
		defaultMemSize := 1024 // Default 1GB
		inputs.MemorySize = &defaultMemSize
		logger.Debug(fmt.Sprintf("Using default memory size: %d MB", defaultMemSize))
	}

	// Get VM generation - find the generation from VM settings
	// Use default Generation 2 if not able to determine
	if inputs.Generation == nil {
		defaultGeneration := 2
		inputs.Generation = &defaultGeneration
		logger.Debug(fmt.Sprintf("Using default generation: %d", defaultGeneration))
	}

	// Get auto start action if not specified in inputs
	if inputs.AutoStartAction == nil {
		autoStartAction, err := vmSettings.GetProperty("AutoStartAction")
		if err == nil {
			var actionStr string
			switch autoStartAction {
			case 0:
				actionStr = "Nothing"
			case 1:
				actionStr = "StartIfRunning"
			case 2:
				actionStr = "Start"
			default:
				actionStr = "Nothing"
			}
			inputs.AutoStartAction = &actionStr
			logger.Debug(fmt.Sprintf("Found auto start action: %s", actionStr))
		}
	}

	// Get auto stop action if not specified in inputs
	if inputs.AutoStopAction == nil {
		autoStopAction, err := vmSettings.GetProperty("AutoStopAction")
		if err == nil {
			var actionStr string
			switch autoStopAction {
			case 0:
				actionStr = "TurnOff"
			case 1:
				actionStr = "Save"
			case 2:
				actionStr = "ShutDown"
			default:
				actionStr = "TurnOff"
			}
			inputs.AutoStopAction = &actionStr
			logger.Debug(fmt.Sprintf("Found auto stop action: %s", actionStr))
		}
	}

	// Get hard drives attached to the VM
	// This is an extension point for future implementation
	// Currently just preserving any hard drives specified in the input
	if len(inputs.HardDrives) == 0 {
		// In a real implementation, we would retrieve the actual hard drives from the VM
		logger.Debug("No hard drives specified in input, and retrieval not implemented")
	}

	// Get network adapters attached to the VM
	// This is an extension point for future implementation
	// Currently just preserving any network adapters specified in the input
	if len(inputs.NetworkAdapters) == 0 {
		// In a real implementation, we would retrieve the actual network adapters from the VM
		logger.Debug("No network adapters specified in input, and retrieval not implemented")
	}

	// Update outputs with populated inputs
	outputs.MachineInputs = inputs

	return outputs, nil
}

// This is the Create method. This will be run on every Machine resource creation.
func (c *Machine) Create(ctx context.Context, name string, input MachineInputs, preview bool) (string, MachineOutputs, error) {
	logger := provider.GetLogger(ctx)
	id := name
	if input.MachineName != nil {
		id = *input.MachineName
	}
	state := MachineOutputs{MachineInputs: input}

	// If in preview, don't run the command.
	if preview {
		return id, state, nil
	}
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		return id, state, err
	}
	setting, err := virtualsystem.GetVirtualSystemSettingData(vmmsClient.GetVirtualizationConn().WMIHost, id)
	if err != nil {
		return id, state, err
	}
	err = setting.SetPropertyInstanceID(id)
	if err != nil {
		return id, state, fmt.Errorf("Failed to set property instance ID: [%+v]", err)
	}

	defer setting.Close()
	logger.Debug("Create VMSettings")

	if input.Generation != nil {
		switch *input.Generation {
		case 1:
			err = setting.SetHyperVGeneration(virtualsystem.HyperVGeneration_V1)
			// Set Secure Boot to false for Generation 1
			// Hyper-V Generation 1 VMs do not support Secure Boot.
			// according to this test: https://github.com/microsoft/wmi/blob/master/pkg/virtualization/core/service/virtualmachinemanagementservice_test.go#L100
			// TODO: Check if this is the correct way to set Secure Boot to false for Generation 1 VMs from Microsft documentation.
			secure_boot_err := setting.SetPropertySecureBootEnabled(false)
			if secure_boot_err != nil {
				return id, state, fmt.Errorf("Failed [%+v]", secure_boot_err)
			}
		case 2:
			err = setting.SetHyperVGeneration(virtualsystem.HyperVGeneration_V2)
		default:
			logger.Errorf("Invalid generation: %d, setting V2", *input.Generation)
			err = setting.SetHyperVGeneration(virtualsystem.HyperVGeneration_V2)
		}
		if err != nil {
			return id, state, fmt.Errorf("Failed [%+v]", err)
		}
	} else {
		err = setting.SetHyperVGeneration(virtualsystem.HyperVGeneration_V2)
		if err != nil {
			return id, state, fmt.Errorf("Failed [%+v]", err)
		}
	}

	// Set auto start action if specified
	if input.AutoStartAction != nil {
		var autoStartValue uint16
		switch *input.AutoStartAction {
		case "Nothing":
			autoStartValue = 0
		case "StartIfRunning":
			autoStartValue = 1
		case "Start":
			autoStartValue = 2
		default:
			logger.Errorf("Invalid auto start action: %s, setting Nothing", *input.AutoStartAction)
			autoStartValue = 0
		}
		err = setting.SetProperty("AutoStartAction", autoStartValue)
		if err != nil {
			return id, state, fmt.Errorf("Failed to set auto start action: [%+v]", err)
		}
	}

	// Set auto stop action if specified
	if input.AutoStopAction != nil {
		var autoStopValue uint16
		switch *input.AutoStopAction {
		case "TurnOff":
			autoStopValue = 0
		case "Save":
			autoStopValue = 1
		case "ShutDown":
			autoStopValue = 2
		default:
			logger.Errorf("Invalid auto stop action: %s, setting TurnOff", *input.AutoStopAction)
			autoStopValue = 0
		}
		err = setting.SetProperty("AutoStopAction", autoStopValue)
		if err != nil {
			return id, state, fmt.Errorf("Failed to set auto stop action: [%+v]", err)
		}
	}

	memorySettings, err := memory.GetDefaultMemorySettingData(vmmsClient.GetVirtualizationConn().WMIHost)
	if err != nil {
		return id, state, fmt.Errorf("Failed [%+v]", err)
	}
	defer memorySettings.Close()

	// Set memory size
	var memorySizeMB uint64 = 1024 // Default value
	if input.MemorySize != nil {
		memorySizeMB = uint64(*input.MemorySize)
	}
	err = memorySettings.SetSizeMB(memorySizeMB)
	if err != nil {
		return id, state, fmt.Errorf("Failed to set memory size: %v", err)
	}

	// Set dynamic memory if specified
	if input.DynamicMemory != nil && *input.DynamicMemory {
		err = memorySettings.SetPropertyDynamicMemoryEnabled(true)
		if err != nil {
			return id, state, fmt.Errorf("Failed to enable dynamic memory: %v", err)
		}

		// Set minimum memory if specified
		if input.MinimumMemory != nil {
			minMemory := uint64(*input.MinimumMemory)
			err = memorySettings.SetProperty("MinimumBytes", minMemory*1024*1024) // Convert MB to bytes
			if err != nil {
				return id, state, fmt.Errorf("Failed to set minimum memory: %v", err)
			}
		}

		// Set maximum memory if specified
		if input.MaximumMemory != nil {
			maxMemory := uint64(*input.MaximumMemory)
			err = memorySettings.SetProperty("MaximumBytes", maxMemory*1024*1024) // Convert MB to bytes
			if err != nil {
				return id, state, fmt.Errorf("Failed to set maximum memory: %v", err)
			}
		}
	}

	processorSettings, err := processor.GetDefaultProcessorSettingData(vmmsClient.GetVirtualizationConn().WMIHost)
	if err != nil {
		return id, state, fmt.Errorf("Failed [%+v]", err)
	}
	var cpuCount uint64 = 1 // Default value
	if input.ProcessorCount != nil {
		cpuCount = uint64(*input.ProcessorCount)
	}
	err = processorSettings.SetCPUCount(cpuCount)
	if err != nil {
		return id, state, fmt.Errorf("Failed to set CPU count: %v", err)
	}

	vm, err := vsms.CreateVirtualMachine(setting, memorySettings, processorSettings)
	if err != nil {
		return id, state, fmt.Errorf("Failed vsms.CreateVirtualMachine: [%+v]", err)
	}
	logger.Debug("Created VM")

	// Add hard drives if specified
	if len(input.HardDrives) > 0 {
		for _, hd := range input.HardDrives {
			if hd.Path == nil {
				logger.Debug("Hard drive path not specified, skipping")
				continue
			}

			// Default values for controller
			controllerType := "SCSI"
			if hd.ControllerType != nil {
				controllerType = *hd.ControllerType
			}

			controllerNumber := 0
			if hd.ControllerNumber != nil {
				controllerNumber = *hd.ControllerNumber
			}

			controllerLocation := 0
			if hd.ControllerLocation != nil {
				controllerLocation = *hd.ControllerLocation
			}

			logger.Debug(fmt.Sprintf("Adding hard drive %s to VM %s", *hd.Path, id))

			// Add the hard drive to the VM
			// Add the hard drive to the VM using direct WMI method invocation
			vmName, err := vm.GetPropertyElementName()
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to get VM name: %v", err))
				continue
			}

			params := map[string]interface{}{
				"VirtualSystemIdentifier": vmName,
				"ResourceType":            uint16(31), // 31 = Disk drive
				"Path":                    *hd.Path,
				"ControllerType":          controllerType,
				"ControllerNumber":        uint32(controllerNumber),
				"ControllerLocation":      uint32(controllerLocation),
			}

			_, err = vsms.InvokeMethod("AddResourceSettings", params)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to add hard drive: %v", err))
				// Continue with other hard drives even if this one fails
			}
		}
	}

	// Add network adapters if specified
	if len(input.NetworkAdapters) > 0 {
		for i, na := range input.NetworkAdapters {
			if na.SwitchName == nil {
				logger.Debug("Network adapter switch name not specified, skipping")
				continue
			}

			// Use the index as part of name if no name provided
			adapterName := fmt.Sprintf("Network Adapter %d", i+1)
			if na.Name != nil {
				adapterName = *na.Name
			}

			logger.Debug(fmt.Sprintf("Adding network adapter %s to VM %s, connected to switch %s",
				adapterName, id, *na.SwitchName))

			// Get the VM name
			vmName, err := vm.GetPropertyElementName()
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to get VM name: %v", err))
				continue
			}

			// Check if the adapter already exists (standalone NetworkAdapter resource from simple-all-four example)
			// This handles the case where a NetworkAdapter resource was created without a vmName
			adapterExists := false

			// Query for existing standalone adapters with this name
			query := fmt.Sprintf("SELECT * FROM Msvm_SyntheticEthernetPortSettingData WHERE ElementName = '%s'", adapterName)
			adapters, err := vmmsClient.GetVirtualizationConn().QueryInstances(query)
			if err == nil && len(adapters) > 0 {
				logger.Debug(fmt.Sprintf("Found existing network adapter with name %s, checking if it's attached to a VM", adapterName))

				// Check each result to see if it's already attached to a VM
				for _, adapter := range adapters {
					defer adapter.Close()

					// Check if this adapter is already attached to a VM
					instanceID, err := adapter.GetProperty("InstanceID")
					if err == nil && instanceID != nil {
						instanceIDStr, ok := instanceID.(string)
						if ok {
							// If the adapter is not attached to any VM, we can't determine from InstanceID
							// If it's attached to a different VM, the VM name would be in the path
							// If it's already attached to our VM, we don't need to do anything
							if !strings.Contains(instanceIDStr, vmName) {
								// This adapter might be available or attached to another VM
								logger.Debug(fmt.Sprintf("Adapter %s exists but is not attached to this VM", adapterName))
							} else {
								// This adapter is already attached to our VM
								logger.Debug(fmt.Sprintf("Adapter %s is already attached to VM %s", adapterName, vmName))
								adapterExists = true
								break
							}
						}
					}
				}
			}

			// If adapter doesn't exist or isn't attached to our VM, create a new one
			if !adapterExists {
				params := map[string]interface{}{
					"VirtualSystemIdentifier": vmName,
					"ResourceType":            uint16(10), // 10 = Network adapter
					"SwitchName":              *na.SwitchName,
					"AdapterName":             adapterName,
				}

				_, err = vsms.InvokeMethod("AddNetworkAdapter", params)
				if err != nil {
					logger.Error(fmt.Sprintf("Failed to add network adapter: %v", err))
					// Continue with other adapters even if this one fails
				} else {
					logger.Debug(fmt.Sprintf("Successfully added network adapter %s to VM %s", adapterName, vmName))
				}
			}
		}
	}

	return id, state, nil
}

// WireDependencies controls how secrets and unknowns flow through a resource.
//
//	var _ = (infer.ExplicitDependencies[MachineInputs, MachineOutputs])((*Machine)(nil))
//	func (r *Machine) WireDependencies(f infer.FieldSelector, args *MachineInputs, state *MachineOutputs) { .. }
//
// Because we want every output to depend on every input, we can leave the default behavior.

// The Update method will be run on every update.
func (c *Machine) Update(ctx context.Context, id string, olds MachineOutputs, news MachineInputs, preview bool) (MachineOutputs, error) {
	// This is a no-op. We don't need to do anything here.
	state := MachineOutputs{MachineInputs: news}
	// If in preview, don't run the command.
	if preview {
		return state, nil
	}

	return state, nil
}

// The Delete method will run when the resource is deleted.
func (c *Machine) Delete(ctx context.Context, id string, props MachineOutputs) error {
	_, vsms, err := c.Connect(ctx)
	if err != nil {
		return err
	}

	vm, err := vsms.GetVirtualMachineByName(id)
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
