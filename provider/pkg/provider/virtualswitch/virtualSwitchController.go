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
	"os/exec"
	"strconv"
	"strings"

	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv/provider/pkg/provider/logging"
	"github.com/pulumi/pulumi-hyperv/provider/pkg/provider/util"
	"github.com/pulumi/pulumi-hyperv/provider/pkg/provider/vmms"
)

// The following statements are type assertions to indicate to Go that VirtualSwitch implements the interfaces.
var _ = (infer.CustomResource[VirtualSwitchInputs, VirtualSwitchOutputs])((*VirtualSwitch)(nil))
var _ = (infer.CustomUpdate[VirtualSwitchInputs, VirtualSwitchOutputs])((*VirtualSwitch)(nil))
var _ = (infer.CustomDelete[VirtualSwitchOutputs])((*VirtualSwitch)(nil))

func (c *VirtualSwitch) Connect(ctx context.Context) (*vmms.VMMS, *service.VirtualSystemManagementService, error) {
	logger := logging.GetLogger(ctx)

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
				logger.Warnf("Recovered from panic in NewVMMS: %v", r)
			}
		}()

		vmmsClient, vmmsErr = vmms.NewVMMS(whost)
	}()

	if vmmsErr != nil {
		// Log the error but continue with simulated functionality
		logger.Warnf("Failed to create VMMS client: %v", vmmsErr)
		logger.Infof("Will attempt to use PowerShell fallback methods for VirtualSwitch operations")

		// We'll return nil client and nil vsms and handle it in the resource methods
		// This allows us to fall back to PowerShell commands when WMI services are unavailable
		return nil, nil, nil
	}

	// Check for nil client before proceeding
	if vmmsClient == nil {
		logger.Warnf("VMMS client is nil after creation")
		logger.Infof("Will attempt to use PowerShell fallback methods for VirtualSwitch operations")
		return nil, nil, nil
	}

	// Get the management service with nil check
	vsms := vmmsClient.GetVirtualSystemManagementService()
	if vsms == nil {
		logger.Warnf("Virtual System Management Service is nil - attempting PowerShell fallback")
		logger.Infof("Make sure Hyper-V is properly installed and you have administrator privileges")
		return vmmsClient, nil, nil
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

	// Only try to close if we got a real WMI instance, not when using PowerShell fallback
	if vswitch != nil {
		defer vswitch.Close()
	}

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

	// Create the switch - using PowerShell since it works in all scenarios
	// This is more reliable than WMI which might be unavailable on client Windows
	// Build the PowerShell command for creating the switch
	var cmdArgs []string

	// Base command
	cmdArgs = append(cmdArgs, "New-VMSwitch", "-Name", fmt.Sprintf("\"%s\"", id))

	switch *input.SwitchType {
	case "External":
		// Ensure we have a network adapter name
		if input.NetAdapterName == nil {
			return id, state, fmt.Errorf("netAdapterName is required for External switches")
		}
		cmdArgs = append(cmdArgs, "-NetAdapterName", fmt.Sprintf("\"%s\"", *input.NetAdapterName))

		// Add AllowManagementOS if specified
		if input.AllowManagementOs != nil && *input.AllowManagementOs {
			cmdArgs = append(cmdArgs, "-AllowManagementOS", "$true")
		}

		logger.Debugf(fmt.Sprintf("Creating external switch %s with adapter %s", id, *input.NetAdapterName))

	case "Internal":
		cmdArgs = append(cmdArgs, "-SwitchType", "Internal")
		logger.Debugf(fmt.Sprintf("Creating internal switch %s", id))

	case "Private":
		cmdArgs = append(cmdArgs, "-SwitchType", "Private")
		logger.Debugf(fmt.Sprintf("Creating private switch %s", id))

	default:
		return id, state, fmt.Errorf("invalid switch type: %s. Must be 'External', 'Internal', or 'Private'", *input.SwitchType)
	}

	// Add notes if provided
	if input.Notes != nil {
		cmdArgs = append(cmdArgs, "-Notes", fmt.Sprintf("\"%s\"", *input.Notes))
		logger.Debugf(fmt.Sprintf("Setting notes for switch %s: %s", id, *input.Notes))
	}

	// Find PowerShell executable
	powershellExe, err := util.FindPowerShellExe()
	if err != nil {
		return id, state, err
	}

	// Execute the PowerShell command with panic recovery
	var cmd *exec.Cmd
	var output []byte
	var cmdErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				cmdErr = fmt.Errorf("recovered from panic in PowerShell execution: %v", r)
			}
		}()

		cmd = exec.Command(powershellExe, cmdArgs...)
		output, cmdErr = cmd.CombinedOutput()
	}()

	if cmdErr != nil {
		return id, state, fmt.Errorf("failed to create switch using PowerShell: %v, output: %s", cmdErr, string(output))
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
		logger.Warnf("Failed to connect to Hyper-V: %v. Will try PowerShell as fallback.", err)
		return c.DeleteWithPowerShell(ctx, id)
	}

	// If vmmsClient is nil, that means we couldn't connect to WMI properly
	// Fall back to PowerShell
	if vmmsClient == nil {
		logger.Warnf("VMMS client is nil. Will try PowerShell as fallback.")
		return c.DeleteWithPowerShell(ctx, id)
	}

	// Try to use WMI first
	// Check if the switch exists
	exists, err := ExistsVirtualSwitch(vmmsClient, id)
	if err != nil {
		logger.Warnf("Failed to check if switch exists via WMI: %v. Will try PowerShell as fallback.", err)
		return c.DeleteWithPowerShell(ctx, id)
	}

	if !exists {
		logger.Debugf(fmt.Sprintf("Switch %s not found, nothing to delete", id))
		return nil
	}

	// At this point, we know the switch exists and we have a valid VMMS client
	// We could try to use the WMI method, but for consistency and reliability
	// we'll use PowerShell for the actual deletion
	return c.DeleteWithPowerShell(ctx, id)
}

// DeleteWithPowerShell removes a virtual switch using PowerShell
func (c *VirtualSwitch) DeleteWithPowerShell(ctx context.Context, id string) error {
	logger := logging.GetLogger(ctx)

	// First check if the switch exists
	exists, err := ExistsVirtualSwitchPowerShellFallback(id)
	if err != nil {
		return fmt.Errorf("error checking if switch exists: %v", err)
	}

	if !exists {
		logger.Debugf(fmt.Sprintf("Switch %s not found, nothing to delete", id))
		return nil
	}

	// Find PowerShell executable
	powershellExe, err := util.FindPowerShellExe()
	if err != nil {
		return err
	}

	// Use PowerShell to delete the switch
	var cmd *exec.Cmd
	var output []byte
	var cmdErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				cmdErr = fmt.Errorf("recovered from panic in PowerShell execution: %v", r)
			}
		}()

		cmd = exec.Command(powershellExe, "-Command", fmt.Sprintf("Remove-VMSwitch -Name \"%s\" -Force", id))
		output, cmdErr = cmd.CombinedOutput()
	}()

	if cmdErr != nil {
		return fmt.Errorf("failed to delete switch with PowerShell: %v, output: %s", cmdErr, string(output))
	}

	logger.Debugf(fmt.Sprintf("Deleted virtual switch %s using PowerShell", id))
	return nil
}

// ExistsVirtualSwitch checks if a virtual switch with the given name exists.
func ExistsVirtualSwitch(v *vmms.VMMS, name string) (bool, error) {
	// Check if VMMS client is nil, use PowerShell fallback
	if v == nil {
		return ExistsVirtualSwitchPowerShellFallback(name)
	}

	vConn := v.GetVirtualizationConn()
	if vConn == nil {
		return ExistsVirtualSwitchPowerShellFallback(name)
	}

	// Wrap query in panic recovery to prevent crash
	var switches []*wmi.WmiInstance
	var queryErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				queryErr = fmt.Errorf("recovered from panic in QueryInstances: %v", r)
			}
		}()

		// Add extra safety checks to ensure no nil pointers
		if vConn == nil {
			queryErr = fmt.Errorf("virtualization connection is nil")
			return
		}

		// Check if WMI connection is healthy
		if vConn.WMIHost == nil {
			queryErr = fmt.Errorf("WMI host is nil in connection")
			return
		}

		// This WMI session is already checked when we got the connection

		// We've verified the connection has essential components

		// Build and execute the query
		query := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE Caption = 'Virtual Switch' AND ElementName = '%s'", name)
		switches, queryErr = vConn.QueryInstances(query)
	}()

	if queryErr != nil {
		// If WMI query fails, try PowerShell fallback
		return ExistsVirtualSwitchPowerShellFallback(name)
	}

	return len(switches) > 0, nil
}

// ExistsVirtualSwitchPowerShellFallback uses PowerShell to check if a virtual switch exists.
func ExistsVirtualSwitchPowerShellFallback(name string) (bool, error) {
	// Use a PowerShell command to check if the switch exists
	// The -ErrorAction SilentlyContinue prevents errors from being displayed if the switch doesn't exist
	// The Count property will be 0 if no switch exists, and 1 or more if it does
	powershellCmd := fmt.Sprintf("(Get-VMSwitch -Name \"%s\" -ErrorAction SilentlyContinue | Measure-Object).Count", name)

	// Find PowerShell executable
	powershellExe, err := util.FindPowerShellExe()
	if err != nil {
		return false, err
	}

	// Execute the PowerShell command with panic recovery
	var cmd *exec.Cmd
	var output []byte
	var cmdErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				cmdErr = fmt.Errorf("recovered from panic in PowerShell execution: %v", r)
			}
		}()

		cmd = exec.Command(powershellExe, "-Command", powershellCmd)
		output, cmdErr = cmd.CombinedOutput()
	}()

	if cmdErr != nil {
		return false, fmt.Errorf("failed to check switch existence using PowerShell: %v, output: %s", cmdErr, string(output))
	}

	// Parse the output, which should be a number (0 or 1)
	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return false, fmt.Errorf("failed to parse PowerShell output: %v, output: %s", err, string(output))
	}

	return count > 0, nil
}

// GetVirtualSwitch gets a virtual switch by name.
func GetVirtualSwitch(v *vmms.VMMS, name string) (*wmi.WmiInstance, error) {
	// If VMMS client is nil, first check if the switch exists using PowerShell
	if v == nil {
		exists, err := ExistsVirtualSwitchPowerShellFallback(name)
		if err != nil {
			return nil, fmt.Errorf("failed to check if switch exists using PowerShell: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("unable to find the Virtual Switch %s", name)
		}
		// Return nil to indicate PowerShell fallback mode
		return nil, nil
	}

	vConn := v.GetVirtualizationConn()
	if vConn == nil {
		// Same as above, check existence and return nil for PowerShell fallback
		exists, err := ExistsVirtualSwitchPowerShellFallback(name)
		if err != nil {
			return nil, fmt.Errorf("failed to check if switch exists using PowerShell: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("unable to find the Virtual Switch %s", name)
		}
		return nil, nil
	}

	// Wrap query in panic recovery to prevent crash
	var switches []*wmi.WmiInstance
	var queryErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				queryErr = fmt.Errorf("recovered from panic in QueryInstances: %v", r)
			}
		}()

		// Safety checks to prevent nil pointer dereference
		if vConn == nil {
			queryErr = fmt.Errorf("virtualization connection is nil")
			return
		}

		if vConn.WMIHost == nil {
			queryErr = fmt.Errorf("WMI host is nil in connection")
			return
		}

		// This WMI session is already checked when we got the connection

		// Build and execute the query
		query := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE Caption = 'Virtual Switch' AND ElementName = '%s'", name)
		switches, queryErr = vConn.QueryInstances(query)
	}()

	if queryErr != nil {
		// If WMI query fails, try PowerShell fallback
		exists, err := ExistsVirtualSwitchPowerShellFallback(name)
		if err != nil {
			return nil, fmt.Errorf("failed to check if switch exists using PowerShell: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("unable to find the Virtual Switch %s", name)
		}
		return nil, nil
	}

	if len(switches) == 0 {
		return nil, fmt.Errorf("unable to find the Virtual Switch %s", name)
	}

	return switches[0], nil
}
