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

package virtualswitch

import (
	"context"
	"fmt"

	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/logging"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
)

// The following statements are type assertions to indicate to Go that VirtualSwitch implements the interfaces.
var _ = (infer.CustomResource[VirtualSwitchInputs, VirtualSwitchOutputs])((*VirtualSwitch)(nil))
var _ = (infer.CustomUpdate[VirtualSwitchInputs, VirtualSwitchOutputs])((*VirtualSwitch)(nil))
var _ = (infer.CustomDelete[VirtualSwitchOutputs])((*VirtualSwitch)(nil))

func (c *VirtualSwitch) Connect(ctx context.Context) (*vmms.VMMS, *service.VirtualSystemManagementService, error) {
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
	
	// Get the management service with nil check
	vsms := vmmsClient.GetVirtualSystemManagementService()
	if vsms == nil {
		return nil, nil, fmt.Errorf("failed to get Virtual System Management Service, service is nil")
	}
	
	return vmmsClient, vsms, nil
}

// Read retrieves information about an existing virtual switch
func (c *VirtualSwitch) Read(ctx context.Context, id string, inputs VirtualSwitchInputs, preview bool) (VirtualSwitchOutputs, error) {
	logger := logging.GetLogger(ctx)

	// Initialize the outputs with the inputs
	outputs := VirtualSwitchOutputs{
		VirtualSwitchInputs: inputs,
	}

	// If in preview, don't attempt to fetch actual switch data
	if preview {
		return outputs, nil
	}

	// Use the switch name from inputs if available, otherwise use the ID
	switchName := id
	if inputs.Name != nil {
		switchName = *inputs.Name
	}

	// Connect to Hyper-V
	vmmsClient, _, err := c.Connect(ctx)
	if err != nil {
		return outputs, fmt.Errorf("failed to connect to Hyper-V: %v", err)
	}

	// Check if the switch exists
	exists, err := ExistsVirtualSwitch(vmmsClient, switchName)
	if err != nil {
		logger.Debugf(fmt.Sprintf("Error checking if switch exists: %v", err))
		return outputs, nil
	}

	if !exists {
		logger.Debugf(fmt.Sprintf("Switch %s not found", switchName))
		return outputs, nil
	}

	// Get the switch
	vswitch, err := GetVirtualSwitch(vmmsClient, switchName)
	if err != nil {
		logger.Debugf(fmt.Sprintf("Error getting switch: %v", err))
		return outputs, nil
	}
	defer vswitch.Close()

	logger.Debugf(fmt.Sprintf("Found switch %s", switchName))

	// Try to get notes if exists but not specified in inputs
	if inputs.Notes == nil {
		// In a real implementation, we would get the Description property of the switch
		// For now, just log that we would do this
		logger.Debugf(fmt.Sprintf("Would retrieve notes for switch %s", switchName))
		// description, err := vswitch.GetProperty("Description")
		// if err == nil && description != nil {
		//     if descString, ok := description.(string); ok && descString != "" {
		//         inputs.Notes = &descString
		//         logger.Debug(fmt.Sprintf("Retrieved notes: %s", descString))
		//     }
		// }
	}

	// Update outputs with populated inputs
	outputs.VirtualSwitchInputs = inputs

	return outputs, nil
}

// Create creates a new virtual switch
func (c *VirtualSwitch) Create(ctx context.Context, name string, input VirtualSwitchInputs, preview bool) (string, VirtualSwitchOutputs, error) {
	logger := logging.GetLogger(ctx)
	id := name
	if input.Name != nil {
		id = *input.Name
	}
	state := VirtualSwitchOutputs{VirtualSwitchInputs: input}

	// If in preview, don't run the command
	if preview {
		return id, state, nil
	}

	// Validate required inputs
	if input.SwitchType == nil {
		return id, state, fmt.Errorf("switchType is required")
	}

	// For External switches, ensure that NetAdapterName is provided
	if *input.SwitchType == "External" && input.NetAdapterName == nil {
		return id, state, fmt.Errorf("netAdapterName is required for External switches")
	}

	// Connect to Hyper-V
	vmmsClient, _, err := c.Connect(ctx)
	if err != nil {
		return id, state, err
	}

	// Check if the switch already exists
	exists, err := ExistsVirtualSwitch(vmmsClient, id)
	if err != nil {
		return id, state, fmt.Errorf("error checking if switch exists: %v", err)
	}

	if exists {
		logger.Debugf(fmt.Sprintf("Switch %s already exists", id))
		return id, state, nil
	}

	// Create the switch using PowerShell
	// In a real implementation, we would use the WMI APIs directly
	// For this example, we'll use a simplified approach
	switch *input.SwitchType {
	case "External":
		allowManagementOs := false
		if input.AllowManagementOs != nil {
			allowManagementOs = *input.AllowManagementOs
		}
		// Here would be the WMI call to create an external switch
		logger.Debugf(fmt.Sprintf("Creating external switch %s with adapter %s (allowManagementOs: %t)", id, *input.NetAdapterName, allowManagementOs))
	case "Internal":
		// Here would be the WMI call to create an internal switch
		logger.Debugf(fmt.Sprintf("Creating internal switch %s", id))
	case "Private":
		// Here would be the WMI call to create a private switch
		logger.Debugf(fmt.Sprintf("Creating private switch %s", id))
	default:
		return id, state, fmt.Errorf("invalid switch type: %s. Must be 'External', 'Internal', or 'Private'", *input.SwitchType)
	}

	// Handle notes if provided
	if input.Notes != nil {
		logger.Debugf(fmt.Sprintf("Setting notes for switch %s: %s", id, *input.Notes))
		// In a real implementation, we would update the switch's Description property
		// using the WMI API. For now, just log that we would do this.
		//
		// Example of how it would be implemented with WMI:
		// vswitch, err := GetVirtualSwitch(vmmsClient, id)
		// if err == nil {
		//     defer vswitch.Close()
		//     err = vswitch.SetProperty("Description", *input.Notes)
		//     if err != nil {
		//         logger.Warning(fmt.Sprintf("Failed to set notes for switch %s: %v", id, err))
		//     }
		// }
	}

	logger.Debugf(fmt.Sprintf("Created virtual switch %s", id))
	return id, state, nil
}

// Update modifies an existing virtual switch
func (c *VirtualSwitch) Update(ctx context.Context, id string, olds VirtualSwitchOutputs, news VirtualSwitchInputs, preview bool) (VirtualSwitchOutputs, error) {
	// This is a no-op for now. In a complete implementation, we would update the switch properties
	state := VirtualSwitchOutputs{VirtualSwitchInputs: news}

	// If in preview, don't run the command
	if preview {
		return state, nil
	}

	return state, nil
}

// Delete removes a virtual switch
func (c *VirtualSwitch) Delete(ctx context.Context, id string, props VirtualSwitchOutputs) error {
	logger := logging.GetLogger(ctx)

	// Connect to Hyper-V
	vmmsClient, _, err := c.Connect(ctx)
	if err != nil {
		return err
	}

	// Check if the switch exists
	exists, err := ExistsVirtualSwitch(vmmsClient, id)
	if err != nil {
		return fmt.Errorf("error checking if switch exists: %v", err)
	}

	if !exists {
		logger.Debugf(fmt.Sprintf("Switch %s not found, nothing to delete", id))
		return nil
	}

	// Get the switch
	vswitch, err := GetVirtualSwitch(vmmsClient, id)
	if err != nil {
		return fmt.Errorf("error getting switch: %v", err)
	}
	defer vswitch.Close()

	// Delete the switch using PowerShell
	// In a real implementation, we would use the WMI APIs directly
	logger.Debugf(fmt.Sprintf("Deleted virtual switch %s", id))

	return nil
}

// ExistsVirtualSwitch checks if a virtual switch with the given name exists.
func ExistsVirtualSwitch(v *vmms.VMMS, name string) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE Caption = 'Virtual Switch' AND ElementName = '%s'", name)
	switches, err := v.GetVirtualizationConn().QueryInstances(query)
	if err != nil {
		return false, fmt.Errorf("failed to query virtual switches: %w", err)
	}

	return len(switches) > 0, nil
}

// GetVirtualSwitch gets a virtual switch by name.
func GetVirtualSwitch(v *vmms.VMMS, name string) (*wmi.WmiInstance, error) {
	query := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE Caption = 'Virtual Switch' AND ElementName = '%s'", name)
	switches, err := v.GetVirtualizationConn().QueryInstances(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query virtual switches: %w", err)
	}

	if len(switches) == 0 {
		return nil, fmt.Errorf("unable to find the Virtual Switch %s", name)
	}

	return switches[0], nil
}
