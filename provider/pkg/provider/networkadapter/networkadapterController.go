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

package networkadapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
)

// Type assertions to indicate that NetworkAdapter implements the required interfaces.
var _ = (infer.CustomResource[NetworkAdapterInputs, NetworkAdapterOutputs])((*NetworkAdapter)(nil))
var _ = (infer.CustomUpdate[NetworkAdapterInputs, NetworkAdapterOutputs])((*NetworkAdapter)(nil))
var _ = (infer.CustomDelete[NetworkAdapterOutputs])((*NetworkAdapter)(nil))

func (c *NetworkAdapter) Connect(ctx context.Context) (*vmms.VMMS, *service.VirtualSystemManagementService, error) {
	logger := provider.GetLogger(ctx)

	// Create the VMMS client.
	config := infer.GetConfig[common.Config](ctx)
	var whost *host.WmiHost
	if config.Host != "" {
		whost = host.NewWmiHost(config.Host)
	} else {
		whost = host.NewWmiLocalHost()
	}

	// Wrap vmms creation in panic recovery
	var vmmsClient *vmms.VMMS
	var vmmsErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				vmmsErr = fmt.Errorf("recovered from panic in NewVMMS: %v", r)
				logger.Debug(fmt.Sprintf("Recovered from panic in NewVMMS: %v", r))
			}
		}()

		vmmsClient, vmmsErr = vmms.NewVMMS(whost)
	}()

	if vmmsErr != nil {
		// Log the error but continue with simulated functionality
		logger.Debug(fmt.Sprintf("Failed to create VMMS client: %v", vmmsErr))
		logger.Debug("Will attempt to use fallback methods for NetworkAdapter operations")

		// For now, we'll return the error since network adapters need more complex simulation
		// In a future version, we could implement a PowerShell fallback like for VHD operations
		return nil, nil, fmt.Errorf("failed to create VMMS client: %w", vmmsErr)
	}

	// Check for nil client before proceeding
	if vmmsClient == nil {
		logger.Debug("VMMS client is nil after creation")
		return nil, nil, fmt.Errorf("VMMS client is nil after creation")
	}

	// Get the management service with nil check
	vsms := vmmsClient.GetVirtualSystemManagementService()
	if vsms == nil {
		logger.Debug("Virtual System Management Service is nil - required for network adapter operations")
		logger.Debug("Make sure Hyper-V is properly installed and you have administrator privileges")
		return nil, nil, fmt.Errorf("failed to get Virtual System Management Service, service is nil")
	}

	return vmmsClient, vsms, nil
}

// Read retrieves information about an existing network adapter
func (c *NetworkAdapter) Read(ctx context.Context, id string, inputs NetworkAdapterInputs, preview bool) (NetworkAdapterOutputs, error) {
	logger := provider.GetLogger(ctx)

	// Initialize the outputs with the inputs
	outputs := NetworkAdapterOutputs{
		NetworkAdapterInputs: inputs,
	}

	// If in preview, don't attempt to fetch actual adapter data
	if preview {
		return outputs, nil
	}

	// Get the VM name and adapter name from inputs
	vmName := ""
	adapterName := id
	if inputs.VMName != nil {
		vmName = *inputs.VMName
	}
	if inputs.Name != nil {
		adapterName = *inputs.Name
	}

	// If VM name is not provided, this might be a reference adapter used by a Machine resource
	if vmName == "" {
		logger.Debug("vmName not provided for read - this may be a reference adapter for use in a Machine resource")
		// Return the current state for reference adapters
		adapterId := id
		outputs.AdapterId = &adapterId
		return outputs, nil
	}

	// Connect to Hyper-V
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		return outputs, fmt.Errorf("failed to connect to Hyper-V: %v", err)
	}

	// Get the VM
	vm, err := vsms.GetVirtualMachineByName(vmName)
	if err != nil {
		logger.Debug(fmt.Sprintf("VM %s not found: %v", vmName, err))
		return outputs, nil
	}
	defer vm.Close()

	// Check if the adapter exists (this would be implemented using WMI queries)
	exists, err := ExistsNetworkAdapter(vmmsClient, vm, adapterName)
	if err != nil {
		logger.Debug(fmt.Sprintf("Error checking if adapter exists: %v", err))
		return outputs, nil
	}

	if !exists {
		logger.Debug(fmt.Sprintf("Network adapter %s not found on VM %s", adapterName, vmName))
		return outputs, nil
	}

	// Get the adapter
	adapter, err := GetNetworkAdapter(vmmsClient, vm, adapterName)
	if err != nil {
		logger.Debug(fmt.Sprintf("Error getting adapter: %v", err))
		return outputs, nil
	}
	defer adapter.Close()

	logger.Debug(fmt.Sprintf("Found network adapter %s on VM %s", adapterName, vmName))

	// Get the adapter ID
	adapterId, err := adapter.GetProperty("InstanceID")
	if err == nil && adapterId != nil {
		adapterIdStr := fmt.Sprintf("%v", adapterId)
		outputs.AdapterId = &adapterIdStr
	} else {
		// Fallback to a generated ID if we can't get the real one
		adapterIdStr := fmt.Sprintf("%s-%s", vmName, adapterName)
		outputs.AdapterId = &adapterIdStr
	}

	// Get the MAC address
	macAddress, err := adapter.GetProperty("Address")
	if err == nil && macAddress != nil {
		if macStr, ok := macAddress.(string); ok && macStr != "" {
			outputs.MacAddress = &macStr
		}
	}

	// Get network adapter settings for additional properties
	adapterSettings, err := GetNetworkAdapterSettings(vmmsClient, adapter)
	if err == nil {
		defer adapterSettings.Close()

		// Get VLAN ID
		vlanId, err := adapterSettings.GetProperty("VLANId")
		if err == nil && vlanId != nil {
			if vlanIdInt, ok := vlanId.(uint16); ok && vlanIdInt > 0 {
				vlanIdVal := int(vlanIdInt)
				outputs.VlanId = &vlanIdVal
			} else if vlanIdFloat, ok := vlanId.(float64); ok && vlanIdFloat > 0 {
				vlanIdVal := int(vlanIdFloat)
				outputs.VlanId = &vlanIdVal
			}
		}

		// Get DHCP Guard
		dhcpGuard, err := adapterSettings.GetProperty("DHCPGuard")
		if err == nil && dhcpGuard != nil {
			if dhcpGuardBool, ok := dhcpGuard.(bool); ok {
				outputs.DHCPGuard = &dhcpGuardBool
			}
		}

		// Get Router Guard
		routerGuard, err := adapterSettings.GetProperty("RouterGuard")
		if err == nil && routerGuard != nil {
			if routerGuardBool, ok := routerGuard.(bool); ok {
				outputs.RouterGuard = &routerGuardBool
			}
		}

		// Get Port Mirroring
		portMirroring, err := adapterSettings.GetProperty("PortMirroring")
		if err == nil && portMirroring != nil {
			var portMirroringStr string
			if portMirroringInt, ok := portMirroring.(uint8); ok {
				switch portMirroringInt {
				case 1:
					portMirroringStr = "Source"
				case 2:
					portMirroringStr = "Destination"
				case 3:
					portMirroringStr = "Both"
				default:
					portMirroringStr = "None"
				}
				outputs.PortMirroring = &portMirroringStr
			} else if portMirroringFloat, ok := portMirroring.(float64); ok {
				switch uint8(portMirroringFloat) {
				case 1:
					portMirroringStr = "Source"
				case 2:
					portMirroringStr = "Destination"
				case 3:
					portMirroringStr = "Both"
				default:
					portMirroringStr = "None"
				}
				outputs.PortMirroring = &portMirroringStr
			}
		}

		// Get IEEE Priority Tag
		ieeePriorityTag, err := adapterSettings.GetProperty("IeeePriorityTag")
		if err == nil && ieeePriorityTag != nil {
			if ieeePriorityTagBool, ok := ieeePriorityTag.(bool); ok {
				outputs.IeeePriorityTag = &ieeePriorityTagBool
			}
		}

		// Get VMQ Weight
		vmqWeight, err := adapterSettings.GetProperty("VMQWeight")
		if err == nil && vmqWeight != nil {
			if vmqWeightInt, ok := vmqWeight.(uint32); ok {
				vmqWeightVal := int(vmqWeightInt)
				outputs.VMQWeight = &vmqWeightVal
			} else if vmqWeightFloat, ok := vmqWeight.(float64); ok {
				vmqWeightVal := int(vmqWeightFloat)
				outputs.VMQWeight = &vmqWeightVal
			}
		}
	}

	// Get the connected switch name
	switchPath, err := getConnectedSwitch(vmmsClient, adapter)
	if err == nil && switchPath != "" {
		// Get the switch object
		switchObj, err := vmmsClient.GetVirtualizationConn().GetInstance(switchPath)
		if err == nil {
			defer switchObj.Close()

			// Get the switch name
			switchName, err := switchObj.GetProperty("ElementName")
			if err == nil && switchName != nil {
				if switchNameStr, ok := switchName.(string); ok {
					outputs.SwitchName = &switchNameStr
				}
			}
		}
	}

	return outputs, nil
}

// Create creates a new network adapter
func (c *NetworkAdapter) Create(ctx context.Context, name string, input NetworkAdapterInputs, preview bool) (string, NetworkAdapterOutputs, error) {
	logger := provider.GetLogger(ctx)
	id := name
	if input.Name != nil {
		id = *input.Name
	}
	state := NetworkAdapterOutputs{NetworkAdapterInputs: input}

	// If in preview, don't run the command
	if preview {
		return id, state, nil
	}

	// Check if vmName is provided
	if input.VMName == nil {
		logger.Debug("vmName not provided - this may be a reference adapter for use in a Machine resource")
		// In this case, we'll create a "virtual" adapter state that can be referenced
		// but the actual adapter will be created by the Machine resource
		adapterIdStr := fmt.Sprintf("ref-%s", id)
		state.AdapterId = &adapterIdStr
		return id, state, nil
	}
	if input.SwitchName == nil {
		return id, state, fmt.Errorf("switchName is required")
	}

	// Connect to Hyper-V
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		return id, state, err
	}

	// Get the VM
	vm, err := vsms.GetVirtualMachineByName(*input.VMName)
	if err != nil {
		return id, state, fmt.Errorf("VM %s not found: %v", *input.VMName, err)
	}
	defer vm.Close()

	// Check if the adapter already exists
	exists, err := ExistsNetworkAdapter(vmmsClient, vm, id)
	if err != nil {
		return id, state, fmt.Errorf("error checking if adapter exists: %v", err)
	}

	if exists {
		logger.Debug(fmt.Sprintf("Network adapter %s already exists on VM %s", id, *input.VMName))
		return id, state, nil
	}

	logger.Debug(fmt.Sprintf("Creating network adapter %s on VM %s", id, *input.VMName))

	// Add a new network adapter to the VM
	var na *wmi.WmiInstance

	// Using the AddNetworkAdapter method
	params := map[string]interface{}{
		"TargetSystem": vm.Msvm_ComputerSystem.InstancePath(),
		"ElementName":  id,
	}

	// If MAC address is provided, set it
	if input.MacAddress != nil {
		params["Address"] = *input.MacAddress
		params["StaticMacAddress"] = true
	} else {
		params["StaticMacAddress"] = false
	}

	result, err := vsms.WmiInstance.InvokeMethod("AddVirtualSystemResources", params)
	if err != nil {
		return id, state, fmt.Errorf("failed to add network adapter: %w", err)
	}

	// Check the return value
	if len(result) < 1 {
		return id, state, fmt.Errorf("unexpected empty result from AddVirtualSystemResources")
	}

	resultMap, ok := result[0].(map[string]interface{})
	if !ok {
		return id, state, fmt.Errorf("unexpected result type from AddVirtualSystemResources")
	}

	returnValue, ok := resultMap["ReturnValue"]
	if !ok {
		return id, state, fmt.Errorf("ReturnValue not found in result")
	}

	returnValueInt, ok := returnValue.(uint32)
	if !ok {
		// Try to convert it to uint32 if it's not already
		if returnValueFloat, ok := returnValue.(float64); ok {
			returnValueInt = uint32(returnValueFloat)
		} else {
			logger.Debug(fmt.Sprintf("Return value is not uint32 or float64: %T", returnValue))
			returnValueInt = 0
		}
	}

	if returnValueInt != 0 {
		return id, state, fmt.Errorf("add network adapter failed with error code: %d", returnValueInt)
	}

	// Get the resulting network adapter
	naPath, ok := resultMap["ResultingResources"].([]string)
	if !ok || len(naPath) == 0 {
		// Try to get the adapter by querying for it
		adapter, err := GetNetworkAdapter(vmmsClient, vm, id)
		if err != nil {
			return id, state, fmt.Errorf("created adapter but couldn't find it: %w", err)
		}
		na = adapter
	} else {
		// Get the network adapter from the path
		na, err = vmmsClient.GetVirtualizationConn().GetInstance(naPath[0])
		if err != nil {
			return id, state, fmt.Errorf("failed to get network adapter instance: %w", err)
		}
	}

	defer na.Close()

	// Connect to virtual switch if provided
	if input.SwitchName != nil {
		logger.Debug(fmt.Sprintf("Connecting adapter to virtual switch %s", *input.SwitchName))

		// Find the virtual switch to connect to
		switchQuery := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE ElementName = '%s'", *input.SwitchName)
		switches, err := vmmsClient.GetVirtualizationConn().QueryInstances(switchQuery)
		if err != nil {
			return id, state, fmt.Errorf("failed to query virtual switch: %w", err)
		}

		if len(switches) == 0 {
			return id, state, fmt.Errorf("virtual switch %s not found", *input.SwitchName)
		}

		switchInstance := switches[0]
		defer switchInstance.Close()

		// Get network adapter settings
		adapterSettings, err := GetNetworkAdapterSettings(vmmsClient, na)
		if err != nil {
			return id, state, fmt.Errorf("failed to get adapter settings: %w", err)
		}
		defer adapterSettings.Close()

		// Get the switch path
		switchPath := switchInstance.InstancePath()

		// Connect adapter to switch by setting the Connection property
		err = adapterSettings.SetProperty("Connection", []string{switchPath})
		if err != nil {
			return id, state, fmt.Errorf("failed to set switch connection: %w", err)
		}

		// Apply the settings
		adapterPath := adapterSettings.InstancePath()
		params := map[string]interface{}{
			"ResourceSettings": []string{adapterPath},
		}

		result, err := vsms.WmiInstance.InvokeMethod("ModifyResourceSettings", params)
		if err != nil {
			return id, state, fmt.Errorf("failed to connect adapter to switch: %w", err)
		}

		// Check result
		if len(result) < 1 {
			return id, state, fmt.Errorf("unexpected empty result from ModifyResourceSettings")
		}

		resultMap, ok := result[0].(map[string]interface{})
		if !ok {
			return id, state, fmt.Errorf("unexpected result type from ModifyResourceSettings")
		}

		returnValue, ok := resultMap["ReturnValue"]
		if !ok {
			return id, state, fmt.Errorf("ReturnValue not found in result")
		}

		returnValueInt, ok := returnValue.(uint32)
		if !ok {
			if returnValueFloat, ok := returnValue.(float64); ok {
				returnValueInt = uint32(returnValueFloat)
			} else {
				logger.Debug(fmt.Sprintf("Return value is not uint32 or float64: %T", returnValue))
				returnValueInt = 0
			}
		}

		if returnValueInt != 0 && returnValueInt != 4096 {
			return id, state, fmt.Errorf("connect adapter to switch failed with error code: %d", returnValueInt)
		}
	}

	// Configure additional adapter properties if provided
	if input.VlanId != nil || input.DHCPGuard != nil || input.RouterGuard != nil ||
		input.PortMirroring != nil || input.IeeePriorityTag != nil || input.VMQWeight != nil {
		logger.Debug("Setting additional network adapter properties")

		// Get adapter settings
		adapterSettings, err := GetNetworkAdapterSettings(vmmsClient, na)
		if err != nil {
			return id, state, fmt.Errorf("failed to get adapter settings for properties: %w", err)
		}
		defer adapterSettings.Close()

		// Apply properties as needed
		needsUpdate := false

		// Set VLAN ID if provided
		if input.VlanId != nil {
			err = adapterSettings.SetProperty("VLANId", uint16(*input.VlanId))
			if err != nil {
				return id, state, fmt.Errorf("failed to set VLAN ID: %w", err)
			}
			needsUpdate = true
		}

		// Set DHCP Guard if provided
		if input.DHCPGuard != nil {
			err = adapterSettings.SetProperty("DHCPGuard", *input.DHCPGuard)
			if err != nil {
				return id, state, fmt.Errorf("failed to set DHCPGuard: %w", err)
			}
			needsUpdate = true
		}

		// Set Router Guard if provided
		if input.RouterGuard != nil {
			err = adapterSettings.SetProperty("RouterGuard", *input.RouterGuard)
			if err != nil {
				return id, state, fmt.Errorf("failed to set RouterGuard: %w", err)
			}
			needsUpdate = true
		}

		// Set Port Mirroring if provided
		if input.PortMirroring != nil {
			// Convert string to numeric value
			portMirroringValue := uint8(0) // None
			switch *input.PortMirroring {
			case "Source":
				portMirroringValue = 1
			case "Destination":
				portMirroringValue = 2
			case "Both":
				portMirroringValue = 3
			}
			err = adapterSettings.SetProperty("PortMirroring", portMirroringValue)
			if err != nil {
				return id, state, fmt.Errorf("failed to set PortMirroring: %w", err)
			}
			needsUpdate = true
		}

		// Set IEEE Priority Tag if provided
		if input.IeeePriorityTag != nil {
			err = adapterSettings.SetProperty("IeeePriorityTag", *input.IeeePriorityTag)
			if err != nil {
				return id, state, fmt.Errorf("failed to set IeeePriorityTag: %w", err)
			}
			needsUpdate = true
		}

		// Set VMQ Weight if provided
		if input.VMQWeight != nil {
			err = adapterSettings.SetProperty("VMQWeight", uint32(*input.VMQWeight))
			if err != nil {
				return id, state, fmt.Errorf("failed to set VMQWeight: %w", err)
			}
			needsUpdate = true
		}

		// Apply the changes if needed
		if needsUpdate {
			adapterPath := adapterSettings.InstancePath()
			params := map[string]interface{}{
				"ResourceSettings": []string{adapterPath},
			}

			result, err := vsms.WmiInstance.InvokeMethod("ModifyResourceSettings", params)
			if err != nil {
				return id, state, fmt.Errorf("failed to apply adapter settings: %w", err)
			}

			// Check result
			if len(result) < 1 {
				return id, state, fmt.Errorf("unexpected empty result from ModifyResourceSettings")
			}

			resultMap, ok := result[0].(map[string]interface{})
			if !ok {
				return id, state, fmt.Errorf("unexpected result type from ModifyResourceSettings")
			}

			returnValue, ok := resultMap["ReturnValue"]
			if !ok {
				return id, state, fmt.Errorf("ReturnValue not found in result")
			}

			returnValueInt, ok := returnValue.(uint32)
			if !ok {
				if returnValueFloat, ok := returnValue.(float64); ok {
					returnValueInt = uint32(returnValueFloat)
				} else {
					logger.Debug(fmt.Sprintf("Return value is not uint32 or float64: %T", returnValue))
					returnValueInt = 0
				}
			}

			if returnValueInt != 0 && returnValueInt != 4096 {
				return id, state, fmt.Errorf("apply adapter settings failed with error code: %d", returnValueInt)
			}
		}
	}

	// Get the adapter ID for state
	adapterId, err := na.GetProperty("InstanceID")
	if err != nil {
		// Use a generated ID if we can't get the actual one
		adapterId = fmt.Sprintf("%s-%s", *input.VMName, id)
	}

	adapterIdStr := fmt.Sprintf("%v", adapterId)
	state.AdapterId = &adapterIdStr

	logger.Debug(fmt.Sprintf("Successfully created network adapter %s on VM %s", id, *input.VMName))
	return id, state, nil
}

// Update modifies an existing network adapter
func (c *NetworkAdapter) Update(ctx context.Context, id string, olds NetworkAdapterOutputs, news NetworkAdapterInputs, preview bool) (NetworkAdapterOutputs, error) {
	logger := provider.GetLogger(ctx)
	state := NetworkAdapterOutputs{NetworkAdapterInputs: news}

	// If in preview, don't run the command
	if preview {
		return state, nil
	}

	// Check if vmName is provided
	if news.VMName == nil {
		logger.Debug("vmName not provided for update - this may be a reference adapter for use in a Machine resource")
		// Return the state as-is since this may be a reference adapter
		return state, nil
	}

	// Preserve adapter ID from old state
	if olds.AdapterId != nil {
		state.AdapterId = olds.AdapterId
	}

	// Connect to Hyper-V
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		return state, err
	}

	// Get the VM
	vm, err := vsms.GetVirtualMachineByName(*news.VMName)
	if err != nil {
		return state, fmt.Errorf("VM %s not found: %v", *news.VMName, err)
	}
	defer vm.Close()

	// Get adapter name
	adapterName := id
	if news.Name != nil {
		adapterName = *news.Name
	}

	// Check if the adapter exists
	exists, err := ExistsNetworkAdapter(vmmsClient, vm, adapterName)
	if err != nil {
		return state, fmt.Errorf("error checking if adapter exists: %v", err)
	}

	if !exists {
		return state, fmt.Errorf("network adapter %s not found on VM %s", adapterName, *news.VMName)
	}

	// Get the adapter
	adapter, err := GetNetworkAdapter(vmmsClient, vm, adapterName)
	if err != nil {
		return state, fmt.Errorf("error getting adapter: %v", err)
	}
	defer adapter.Close()

	// Get the adapter settings
	adapterSettings, err := GetNetworkAdapterSettings(vmmsClient, adapter)
	if err != nil {
		return state, fmt.Errorf("failed to get adapter settings: %w", err)
	}
	defer adapterSettings.Close()

	// Track if we need to update settings
	needsUpdate := false

	// Update MAC address if changed
	if news.MacAddress != nil && (olds.MacAddress == nil || *news.MacAddress != *olds.MacAddress) {
		logger.Debug(fmt.Sprintf("Updating MAC address to %s", *news.MacAddress))

		// Set the MAC address
		err = adapterSettings.SetProperty("Address", *news.MacAddress)
		if err != nil {
			return state, fmt.Errorf("failed to set MAC address: %w", err)
		}

		// Set static MAC address flag
		err = adapterSettings.SetProperty("StaticMacAddress", true)
		if err != nil {
			return state, fmt.Errorf("failed to set static MAC address flag: %w", err)
		}

		needsUpdate = true
	}

	// Update switch connection if changed
	if news.SwitchName != nil && (olds.SwitchName == nil || *news.SwitchName != *olds.SwitchName) {
		logger.Debug(fmt.Sprintf("Updating switch connection to %s", *news.SwitchName))

		// Find the virtual switch to connect to
		switchQuery := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE ElementName = '%s'", *news.SwitchName)
		switches, err := vmmsClient.GetVirtualizationConn().QueryInstances(switchQuery)
		if err != nil {
			return state, fmt.Errorf("failed to query virtual switch: %w", err)
		}

		if len(switches) == 0 {
			return state, fmt.Errorf("virtual switch %s not found", *news.SwitchName)
		}

		switchInstance := switches[0]
		defer switchInstance.Close()

		// Get the switch path
		switchPath := switchInstance.InstancePath()

		// Set the connection
		err = adapterSettings.SetProperty("Connection", []string{switchPath})
		if err != nil {
			return state, fmt.Errorf("failed to set switch connection: %w", err)
		}

		needsUpdate = true
	}

	// Update VLAN ID if changed
	if news.VlanId != nil && (olds.VlanId == nil || *news.VlanId != *olds.VlanId) {
		logger.Debug(fmt.Sprintf("Updating VLAN ID to %d", *news.VlanId))

		err = adapterSettings.SetProperty("VLANId", uint16(*news.VlanId))
		if err != nil {
			return state, fmt.Errorf("failed to set VLAN ID: %w", err)
		}

		needsUpdate = true
	}

	// Update DHCP Guard if changed
	if news.DHCPGuard != nil && (olds.DHCPGuard == nil || *news.DHCPGuard != *olds.DHCPGuard) {
		logger.Debug(fmt.Sprintf("Updating DHCP Guard to %v", *news.DHCPGuard))

		err = adapterSettings.SetProperty("DHCPGuard", *news.DHCPGuard)
		if err != nil {
			return state, fmt.Errorf("failed to set DHCPGuard: %w", err)
		}

		needsUpdate = true
	}

	// Update Router Guard if changed
	if news.RouterGuard != nil && (olds.RouterGuard == nil || *news.RouterGuard != *olds.RouterGuard) {
		logger.Debug(fmt.Sprintf("Updating Router Guard to %v", *news.RouterGuard))

		err = adapterSettings.SetProperty("RouterGuard", *news.RouterGuard)
		if err != nil {
			return state, fmt.Errorf("failed to set RouterGuard: %w", err)
		}

		needsUpdate = true
	}

	// Update Port Mirroring if changed
	if news.PortMirroring != nil && (olds.PortMirroring == nil || *news.PortMirroring != *olds.PortMirroring) {
		logger.Debug(fmt.Sprintf("Updating Port Mirroring to %s", *news.PortMirroring))

		// Convert string to numeric value
		portMirroringValue := uint8(0) // None
		switch *news.PortMirroring {
		case "Source":
			portMirroringValue = 1
		case "Destination":
			portMirroringValue = 2
		case "Both":
			portMirroringValue = 3
		}

		err = adapterSettings.SetProperty("PortMirroring", portMirroringValue)
		if err != nil {
			return state, fmt.Errorf("failed to set PortMirroring: %w", err)
		}

		needsUpdate = true
	}

	// Update IEEE Priority Tag if changed
	if news.IeeePriorityTag != nil && (olds.IeeePriorityTag == nil || *news.IeeePriorityTag != *olds.IeeePriorityTag) {
		logger.Debug(fmt.Sprintf("Updating IEEE Priority Tag to %v", *news.IeeePriorityTag))

		err = adapterSettings.SetProperty("IeeePriorityTag", *news.IeeePriorityTag)
		if err != nil {
			return state, fmt.Errorf("failed to set IeeePriorityTag: %w", err)
		}

		needsUpdate = true
	}

	// Update VMQ Weight if changed
	if news.VMQWeight != nil && (olds.VMQWeight == nil || *news.VMQWeight != *olds.VMQWeight) {
		logger.Debug(fmt.Sprintf("Updating VMQ Weight to %d", *news.VMQWeight))

		err = adapterSettings.SetProperty("VMQWeight", uint32(*news.VMQWeight))
		if err != nil {
			return state, fmt.Errorf("failed to set VMQWeight: %w", err)
		}

		needsUpdate = true
	}

	// Apply the changes if needed
	if needsUpdate {
		logger.Debug("Applying network adapter settings changes")

		adapterPath := adapterSettings.InstancePath()
		params := map[string]interface{}{
			"ResourceSettings": []string{adapterPath},
		}

		result, err := vsms.WmiInstance.InvokeMethod("ModifyResourceSettings", params)
		if err != nil {
			return state, fmt.Errorf("failed to modify adapter settings: %w", err)
		}

		// Check result
		if len(result) < 1 {
			return state, fmt.Errorf("unexpected empty result from ModifyResourceSettings")
		}

		resultMap, ok := result[0].(map[string]interface{})
		if !ok {
			return state, fmt.Errorf("unexpected result type from ModifyResourceSettings")
		}

		returnValue, ok := resultMap["ReturnValue"]
		if !ok {
			return state, fmt.Errorf("ReturnValue not found in result")
		}

		returnValueInt, ok := returnValue.(uint32)
		if !ok {
			if returnValueFloat, ok := returnValue.(float64); ok {
				returnValueInt = uint32(returnValueFloat)
			} else {
				logger.Debug(fmt.Sprintf("Return value is not uint32 or float64: %T", returnValue))
				returnValueInt = 0
			}
		}

		if returnValueInt != 0 && returnValueInt != 4096 {
			return state, fmt.Errorf("modify adapter settings failed with error code: %d", returnValueInt)
		}
	}

	logger.Debug(fmt.Sprintf("Successfully updated network adapter %s on VM %s", adapterName, *news.VMName))
	return state, nil
}

// Delete removes a network adapter
func (c *NetworkAdapter) Delete(ctx context.Context, id string, props NetworkAdapterOutputs) error {
	logger := provider.GetLogger(ctx)

	// Check if vmName is provided
	if props.VMName == nil {
		logger.Debug("vmName not provided for delete - this may be a reference adapter for use in a Machine resource")
		// No real adapter to delete, this was just a reference
		return nil
	}

	// Connect to Hyper-V
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		return err
	}

	// Get the VM
	vm, err := vsms.GetVirtualMachineByName(*props.VMName)
	if err != nil {
		return fmt.Errorf("VM %s not found: %v", *props.VMName, err)
	}
	defer vm.Close()

	// Get adapter name
	adapterName := id
	if props.Name != nil {
		adapterName = *props.Name
	}

	// Check if the adapter exists
	exists, err := ExistsNetworkAdapter(vmmsClient, vm, adapterName)
	if err != nil {
		return fmt.Errorf("error checking if adapter exists: %v", err)
	}

	if !exists {
		logger.Debug(fmt.Sprintf("Network adapter %s not found on VM %s, nothing to delete", adapterName, *props.VMName))
		return nil
	}

	// Get the adapter
	adapter, err := GetNetworkAdapter(vmmsClient, vm, adapterName)
	if err != nil {
		return fmt.Errorf("error getting adapter: %v", err)
	}
	defer adapter.Close()

	// Get the adapter settings to delete
	adapterSettings, err := GetNetworkAdapterSettings(vmmsClient, adapter)
	if err != nil {
		return fmt.Errorf("failed to get adapter settings: %w", err)
	}
	defer adapterSettings.Close()

	// Get path to the adapter setting
	settingPath := adapterSettings.InstancePath()

	// Delete the adapter using RemoveResourceSettings
	logger.Debug(fmt.Sprintf("Deleting network adapter %s from VM %s", adapterName, *props.VMName))

	params := map[string]interface{}{
		"ResourceSettings": []string{settingPath},
	}

	result, err := vsms.WmiInstance.InvokeMethod("RemoveResourceSettings", params)
	if err != nil {
		return fmt.Errorf("failed to remove resource settings: %w", err)
	}

	// Check the return value
	if len(result) < 1 {
		return fmt.Errorf("unexpected empty result from RemoveResourceSettings")
	}

	resultMap, ok := result[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected result type from RemoveResourceSettings")
	}

	returnValue, ok := resultMap["ReturnValue"]
	if !ok {
		return fmt.Errorf("ReturnValue not found in result")
	}

	returnValueInt, ok := returnValue.(uint32)
	if !ok {
		// Try to convert it to uint32 if it's not already
		if returnValueFloat, ok := returnValue.(float64); ok {
			returnValueInt = uint32(returnValueFloat)
		} else {
			logger.Debug(fmt.Sprintf("Return value is not uint32 or float64: %T", returnValue))
			returnValueInt = 0
		}
	}

	if returnValueInt != 0 && returnValueInt != 4096 {
		return fmt.Errorf("remove resource settings failed with error code: %d", returnValueInt)
	}

	logger.Debug(fmt.Sprintf("Successfully deleted network adapter %s from VM %s", adapterName, *props.VMName))
	return nil
}

// ExistsNetworkAdapter checks if a network adapter with the given name exists on a VM.
func ExistsNetworkAdapter(v *vmms.VMMS, vm *virtualsystem.VirtualMachine, name string) (bool, error) {
	// Add defensive checks to prevent nil pointer dereference
	if v == nil {
		return false, fmt.Errorf("vmms object is nil")
	}

	vConn := v.GetVirtualizationConn()
	if vConn == nil {
		return false, fmt.Errorf("virtualization connection is nil")
	}

	if vm == nil || vm.Msvm_ComputerSystem == nil {
		return false, fmt.Errorf("virtual machine or computer system is nil")
	}

	// Query the VM's network adapters
	query := fmt.Sprintf("SELECT * FROM Msvm_SyntheticEthernetPortSettingData WHERE ElementName = '%s' AND InstanceID LIKE '%%\\\\%s\\\\%%'",
		name, vm.Msvm_ComputerSystem.InstancePath())

	// Wrap query in panic recovery to prevent crash
	var adapters []*wmi.WmiInstance
	var queryErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				queryErr = fmt.Errorf("recovered from panic in QueryInstances: %v", r)
			}
		}()

		adapters, queryErr = vConn.QueryInstances(query)
	}()

	if queryErr != nil {
		return false, fmt.Errorf("failed to query network adapters: %w", queryErr)
	}

	exists := len(adapters) > 0

	// Close all adapter instances
	for _, adapter := range adapters {
		adapter.Close()
	}

	return exists, nil
}

// GetNetworkAdapter gets a network adapter by name from a VM.
func GetNetworkAdapter(v *vmms.VMMS, vm *virtualsystem.VirtualMachine, name string) (*wmi.WmiInstance, error) {
	// Add defensive checks to prevent nil pointer dereference
	if v == nil {
		return nil, fmt.Errorf("vmms object is nil")
	}

	vConn := v.GetVirtualizationConn()
	if vConn == nil {
		return nil, fmt.Errorf("virtualization connection is nil")
	}

	if vm == nil || vm.Msvm_ComputerSystem == nil {
		return nil, fmt.Errorf("virtual machine or computer system is nil")
	}

	// Query the VM's network adapters
	query := fmt.Sprintf("SELECT * FROM Msvm_SyntheticEthernetPortSettingData WHERE ElementName = '%s' AND InstanceID LIKE '%%\\\\%s\\\\%%'",
		name, vm.Msvm_ComputerSystem.InstancePath())

	// Wrap query in panic recovery to prevent crash
	var adapters []*wmi.WmiInstance
	var queryErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				queryErr = fmt.Errorf("recovered from panic in QueryInstances: %v", r)
			}
		}()

		adapters, queryErr = vConn.QueryInstances(query)
	}()

	if queryErr != nil {
		return nil, fmt.Errorf("failed to query network adapters: %w", queryErr)
	}

	if len(adapters) == 0 {
		return nil, fmt.Errorf("network adapter %s not found on VM", name)
	}

	// Return the first adapter found with matching name
	// Note: We're not closing this instance as it will be used by the caller
	return adapters[0], nil
}

// GetNetworkAdapterSettings gets the settings for a network adapter.
func GetNetworkAdapterSettings(v *vmms.VMMS, adapter *wmi.WmiInstance) (*wmi.WmiInstance, error) {
	return adapter, nil
}

// getConnectedSwitch returns the path to the virtual switch connected to the adapter.
func getConnectedSwitch(v *vmms.VMMS, adapter *wmi.WmiInstance) (string, error) {
	// Add defensive checks to prevent nil pointer dereference
	if v == nil {
		return "", fmt.Errorf("vmms object is nil")
	}

	if adapter == nil {
		return "", fmt.Errorf("adapter object is nil")
	}

	// Wrap property access in panic recovery
	var connection interface{}
	var propErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				propErr = fmt.Errorf("recovered from panic in GetProperty: %v", r)
			}
		}()

		// Get the Connection property which contains the switch path
		connection, propErr = adapter.GetProperty("Connection")
	}()

	if propErr != nil {
		return "", fmt.Errorf("failed to get connection: %w", propErr)
	}

	// The Connection property should be an array of paths, but we only care about the first one
	if connectionArr, ok := connection.([]string); ok && len(connectionArr) > 0 {
		return connectionArr[0], nil
	}

	// Try it as a different type if the above didn't work
	if connectionAny, ok := connection.([]interface{}); ok && len(connectionAny) > 0 {
		if switchPath, ok := connectionAny[0].(string); ok {
			return switchPath, nil
		}
	}

	return "", fmt.Errorf("no switch connection found")
}

// ParseIPAddresses parses a comma-separated list of IP addresses.
func ParseIPAddresses(ipAddressesStr string) []string {
	if ipAddressesStr == "" {
		return nil
	}
	return strings.Split(ipAddressesStr, ",")
}
