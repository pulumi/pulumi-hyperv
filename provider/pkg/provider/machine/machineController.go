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
	"os/exec"
	"strconv"
	"strings"

	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/memory"
	"github.com/microsoft/wmi/pkg/virtualization/core/processor"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/logging"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/networkadapter"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/util"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
)

// The following statements are not required. They are type assertions to indicate to Go that Machine implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors at this location.

// They would normally be included in the vmController.go file, but they're located here for instructive purposes.
var _ = (infer.CustomResource[MachineInputs, MachineOutputs])((*Machine)(nil))
var _ = (infer.CustomUpdate[MachineInputs, MachineOutputs])((*Machine)(nil))
var _ = (infer.CustomDelete[MachineOutputs])((*Machine)(nil))

func (c *Machine) Connect(ctx context.Context) (*vmms.VMMS, *service.VirtualSystemManagementService, error) {
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
		// Log the error but continue with PowerShell fallback
		logger.Warnf("Failed to create VMMS client: %v", vmmsErr)
		logger.Infof("Will attempt to use PowerShell fallback for machine operations")
		// Return nil, nil to indicate PowerShell fallback should be used
		return nil, nil, nil
	}

	// Check for nil client before proceeding
	if vmmsClient == nil {
		logger.Warnf("VMMS client is nil after creation")
		logger.Infof("Will attempt to use PowerShell fallback for machine operations")
		return nil, nil, nil
	}

	// Get the management service with nil check
	vsms := vmmsClient.GetVirtualSystemManagementService()
	if vsms == nil {
		logger.Warnf("Virtual System Management Service is nil on Windows 10/11, falling back to PowerShell")
		logger.Infof("Windows 10/11 client editions often have more limited WMI service access, using PowerShell fallback")
		// For Windows 10/11, we'll use PowerShell fallback - return client but nil VSMS
		return vmmsClient, nil, nil
	}

	return vmmsClient, vsms, nil
}

// This is the Get Metadata method.
func (c *Machine) Read(ctx context.Context, id string, inputs MachineInputs, preview bool) (MachineOutputs, error) {
	logger := logging.GetLogger(ctx)

	// Initialize the outputs with the inputs
	outputs := MachineOutputs{
		MachineInputs: inputs,
	}

	// The ID is the machine name if it's set, otherwise it's the ID
	machineName := id
	if inputs.MachineName != nil {
		machineName = *inputs.MachineName
	}

	// Always ensure vmId is set even in preview or when VM doesn't exist yet
	EnsureVmId(&outputs, machineName)

	// If in preview, don't attempt to fetch actual VM data
	if preview {
		return outputs, nil
	}

	// Connect to Hyper-V
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		logger.Warnf("Error connecting to Hyper-V: %v", err)
		return outputs, fmt.Errorf("failed to connect to Hyper-V: %v", err)
	}

	// If we have no VMMS client or VSMS, use PowerShell to check if the VM exists
	if vmmsClient == nil || vsms == nil {
		logger.Infof("Using PowerShell fallback to check VM %s", machineName)
		exists, err := checkVMExistsPowerShell(machineName)
		if err != nil {
			logger.Warnf("Error checking VM existence with PowerShell: %v", err)
			return outputs, nil
		}

		if !exists {
			logger.Debugf("Machine %s not found (PowerShell check)", machineName)
			return outputs, nil
		}

		// VM exists - gather information using PowerShell
		return readVMWithPowerShell(ctx, machineName, inputs)
	}

	// Use WMI if available
	vm, err := vsms.GetVirtualMachineByName(machineName)
	if err != nil {
		logger.Debugf("Machine %s not found: %v", machineName, err)
		// Try fallback to PowerShell
		logger.Infof("Falling back to PowerShell to check VM %s", machineName)
		return readVMWithPowerShell(ctx, machineName, inputs)
	}
	defer vm.Close()

	logger.Debugf("Found machine %s", machineName)

	// Get VM ID (ElementName in Hyper-V lingo)
	vmId, err := vm.GetPropertyElementName()
	if err == nil && vmId != "" {
		outputs.VmId = &vmId
		logger.Debugf("VM ID: %s", vmId)
	}

	// Get the VM settings data
	vmSettings, err := virtualsystem.GetVirtualSystemSettingData(vmmsClient.GetVirtualizationConn().WMIHost, machineName)
	if err != nil {
		logger.Debugf("Failed to get VM settings: %v", err)
		return outputs, nil
	}
	defer vmSettings.Close()

	// Get processor count - find the processor setting from VM settings
	// Use default values if not able to get processor setting
	if inputs.ProcessorCount == nil {
		defaultProcCount := 1
		inputs.ProcessorCount = &defaultProcCount
		logger.Debugf("Using default processor count: %d", defaultProcCount)
	}

	// Get memory size - find the memory setting from VM settings
	// Use default values if not able to get memory setting
	if inputs.MemorySize == nil {
		defaultMemSize := 1024 // Default 1GB
		inputs.MemorySize = &defaultMemSize
		logger.Debugf("Using default memory size: %d MB", defaultMemSize)
	}

	// Get VM generation - find the generation from VM settings
	// Use default Generation 2 if not able to determine
	if inputs.Generation == nil {
		defaultGeneration := 2
		inputs.Generation = &defaultGeneration
		logger.Debugf("Using default generation: %d", defaultGeneration)
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
			logger.Debugf("Found auto start action: %s", actionStr)
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
			logger.Debugf("Found auto stop action: %s", actionStr)
		}
	}

	// Get hard drives attached to the VM
	// This is an extension point for future implementation
	// Currently just preserving any hard drives specified in the input
	if len(inputs.HardDrives) == 0 {
		// In a real implementation, we would retrieve the actual hard drives from the VM
		logger.Debugf("No hard drives specified in input, and retrieval not implemented")
	}

	// Get network adapters attached to the VM
	// This is an extension point for future implementation
	// Currently just preserving any network adapters specified in the input
	if len(inputs.NetworkAdapters) == 0 {
		// In a real implementation, we would retrieve the actual network adapters from the VM
		logger.Debugf("No network adapters specified in input, and retrieval not implemented")
	}

	// Update outputs with populated inputs
	outputs.MachineInputs = inputs

	return outputs, nil
}

// This is the Create method. This will be run on every Machine resource creation.
// createVMWithPowerShell creates a virtual machine using PowerShell cmdlets when WMI services are unavailable.
// This is a fallback method used primarily on Windows 10/11 where VSMS might not be fully available.
func createVMWithPowerShell(ctx context.Context, id string, input MachineInputs) (string, MachineOutputs, error) {
	logger := logging.GetLogger(ctx)
	state := MachineOutputs{MachineInputs: input}

	// Ensure vmId is set
	EnsureVmId(&state, id)

	// Find PowerShell executable
	powershellExe, err := util.FindPowerShellExe()
	if err != nil {
		return id, state, err
	}

	// Build the PowerShell command for creating a VM
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "New-VM", "-Name", fmt.Sprintf("\"%s\"", id))

	// Add memory parameter
	if input.MemorySize != nil {
		cmdArgs = append(cmdArgs, "-MemoryStartupBytes", fmt.Sprintf("%d", *input.MemorySize*1024*1024)) // Convert MB to bytes
	} else {
		// Default to 1GB
		cmdArgs = append(cmdArgs, "-MemoryStartupBytes", "1073741824") // 1GB in bytes
	}

	// Add generation parameter
	if input.Generation != nil {
		cmdArgs = append(cmdArgs, "-Generation", fmt.Sprintf("%d", *input.Generation))
	} else {
		// Default to Gen 2
		cmdArgs = append(cmdArgs, "-Generation", "2")
	}

	// Use the first network adapter's switch name if available
	if len(input.NetworkAdapters) > 0 && input.NetworkAdapters[0].SwitchName != nil {
		cmdArgs = append(cmdArgs, "-SwitchName", fmt.Sprintf("\"%s\"", *input.NetworkAdapters[0].SwitchName))
	}

	// Create the VM without adding drives initially
	// We'll add drives and network adapters separately after creation
	cmdArgs = append(cmdArgs, "-NoVHD")

	// Execute the PowerShell command
	var cmd *exec.Cmd
	var output []byte
	var cmdErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				cmdErr = fmt.Errorf("recovered from panic in PowerShell execution: %v", r)
				logger.Warnf("Recovered from panic in PowerShell execution: %v", r)
			}
		}()

		cmd = exec.Command(powershellExe, cmdArgs...)
		output, cmdErr = cmd.CombinedOutput()
	}()

	if cmdErr != nil {
		return id, state, fmt.Errorf("failed to create VM using PowerShell: %v, output: %s", cmdErr, string(output))
	}

	logger.Debugf("Created VM %s with PowerShell", id)

	// Set processor count if specified
	if input.ProcessorCount != nil {
		procCmd := fmt.Sprintf("Set-VMProcessor -VMName \"%s\" -Count %d", id, *input.ProcessorCount)
		_, err := util.RunPowerShellCommand(procCmd)
		if err != nil {
			logger.Warnf("Failed to set processor count: %v", err)
			// Continue despite error
		}
	}

	// Configure dynamic memory if specified
	if input.DynamicMemory != nil && *input.DynamicMemory {
		var minMem, maxMem int64

		if input.MinimumMemory != nil {
			minMem = int64(*input.MinimumMemory) * 1024 * 1024 // Convert MB to bytes
		} else {
			minMem = 512 * 1024 * 1024 // 512MB default
		}

		if input.MaximumMemory != nil {
			maxMem = int64(*input.MaximumMemory) * 1024 * 1024 // Convert MB to bytes
		} else if input.MemorySize != nil {
			maxMem = int64(*input.MemorySize) * 2 * 1024 * 1024 // Double startup memory if maximum not specified
		} else {
			maxMem = 2 * 1024 * 1024 * 1024 // 2GB default
		}

		memCmd := fmt.Sprintf("Set-VMMemory -VMName \"%s\" -DynamicMemoryEnabled $true -MinimumBytes %d -MaximumBytes %d",
			id, minMem, maxMem)
		_, err := util.RunPowerShellCommand(memCmd)
		if err != nil {
			logger.Warnf("Failed to configure dynamic memory: %v", err)
			// Continue despite error
		}
	}

	// Set auto start/stop actions if specified
	if input.AutoStartAction != nil {
		var startAction string
		switch strings.ToLower(*input.AutoStartAction) {
		case "nothing":
			startAction = "Nothing"
		case "startifrunning":
			startAction = "StartIfRunning"
		case "start":
			startAction = "Start"
		default:
			startAction = "Nothing"
		}

		autoStartCmd := fmt.Sprintf("Set-VM -VMName \"%s\" -AutomaticStartAction %s", id, startAction)
		_, err := util.RunPowerShellCommand(autoStartCmd)
		if err != nil {
			logger.Warnf("Failed to set auto start action: %v", err)
			// Continue despite error
		}
	}

	if input.AutoStopAction != nil {
		var stopAction string
		switch strings.ToLower(*input.AutoStopAction) {
		case "turnoff":
			stopAction = "TurnOff"
		case "save":
			stopAction = "Save"
		case "shutdown":
			stopAction = "ShutDown"
		default:
			stopAction = "TurnOff"
		}

		autoStopCmd := fmt.Sprintf("Set-VM -VMName \"%s\" -AutomaticStopAction %s", id, stopAction)
		_, err := util.RunPowerShellCommand(autoStopCmd)
		if err != nil {
			logger.Warnf("Failed to set auto stop action: %v", err)
			// Continue despite error
		}
	}

	// Add hard drives if specified
	if len(input.HardDrives) > 0 {
		for i, hd := range input.HardDrives {
			if hd.Path == nil {
				logger.Debugf("Hard drive path not specified, skipping")
				continue
			}

			// Default values for controller
			controllerType := "SCSI"
			if hd.ControllerType != nil {
				controllerType = *hd.ControllerType
			}

			// Default to first port
			controllerNumber := 0
			if hd.ControllerNumber != nil {
				controllerNumber = *hd.ControllerNumber
			}

			// Default to sequential ports
			controllerLocation := i
			if hd.ControllerLocation != nil {
				controllerLocation = *hd.ControllerLocation
			}

			hdCmd := fmt.Sprintf("Add-VMHardDiskDrive -VMName \"%s\" -Path \"%s\" -ControllerType %s -ControllerNumber %d -ControllerLocation %d",
				id, *hd.Path, controllerType, controllerNumber, controllerLocation)
			_, err := util.RunPowerShellCommand(hdCmd)
			if err != nil {
				logger.Warnf("Failed to add hard drive: %v", err)
				// Continue with other hard drives despite error
			} else {
				logger.Debugf("Added hard drive %s to VM %s", *hd.Path, id)
			}
		}
	}

	// Add network adapters if specified
	if len(input.NetworkAdapters) > 0 {
		for i, na := range input.NetworkAdapters {
			if na.SwitchName == nil {
				logger.Debugf("Network adapter switch name not specified, skipping")
				continue
			}

			// Use the index as part of name if no name provided
			adapterName := fmt.Sprintf("Network Adapter %d", i+1)
			if na.Name != nil {
				adapterName = *na.Name
			}

			// Create the adapter and connect it to the switch
			naCmd := fmt.Sprintf("Add-VMNetworkAdapter -VMName \"%s\" -Name \"%s\" -SwitchName \"%s\"",
				id, adapterName, *na.SwitchName)

			// Add MAC address if specified
			if na.MacAddress != nil && *na.MacAddress != "" {
				naCmd += fmt.Sprintf(" -StaticMacAddress \"%s\"", *na.MacAddress)
			}

			_, err := util.RunPowerShellCommand(naCmd)
			if err != nil {
				logger.Warnf("Failed to add network adapter: %v", err)
				// Continue with other adapters despite error
			} else {
				logger.Debugf("Added network adapter %s to VM %s", adapterName, id)
			}
		}
	}

	// Start the VM automatically after creation
	startCmd := fmt.Sprintf("Start-VM -Name \"%s\"", id)
	startOutput, startErr := util.RunPowerShellCommand(startCmd)
	if startErr != nil {
		// Check for specific error conditions
		if strings.Contains(startOutput, "Not enough memory in the system to start the virtual machine") {
			// Handle out of memory error with specific guidance
			memoryRequiredStr := "default"
			if input.MemorySize != nil {
				memoryRequiredStr = fmt.Sprintf("%d", *input.MemorySize)
			}

			logger.Errorf("Memory error starting VM %s: %s", id, startOutput)
			return id, state, fmt.Errorf("failed to start VM due to insufficient memory: the system does not have enough memory to allocate %s MB for this VM. "+
				"Try reducing the memory allocation, closing other applications, or adding more RAM to the host system", memoryRequiredStr)
		} else if strings.Contains(startOutput, "0x8007000E") {
			// Generic out of resources error (could be memory or something else)
			return id, state, fmt.Errorf("failed to start VM due to insufficient system resources (error 0x8007000E). " +
				"Try reducing VM resource allocation, closing other applications, or adding more resources to the host system")
		} else if strings.Contains(startOutput, "could not initialize memory") {
			// Memory initialization error
			return id, state, fmt.Errorf("failed to start VM due to memory initialization error. " +
				"This could be due to insufficient memory, memory fragmentation, or a system configuration issue")
		} else {
			// Log the error but continue since the VM is created
			logger.Warnf("Failed to start VM: %v", startErr)
			logger.Debugf("Start-VM output: %s", startOutput)
		}
	} else {
		logger.Infof("Started VM %s", id)
	}

	return id, state, nil
}

// checkVMExistsPowerShell checks if a VM exists using PowerShell
func checkVMExistsPowerShell(vmName string) (bool, error) {

	// Command to check if VM exists
	// The -ErrorAction SilentlyContinue prevents errors if the VM doesn't exist
	cmd := fmt.Sprintf("(Get-VM -Name \"%s\" -ErrorAction SilentlyContinue | Measure-Object).Count", vmName)

	output, err := util.RunPowerShellCommand(cmd)
	if err != nil {
		return false, fmt.Errorf("failed to check VM existence: %v", err)
	}

	// Parse the output (should be "0" or "1")
	count, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		return false, fmt.Errorf("failed to parse PowerShell output: %v, output: %s", err, output)
	}

	return count > 0, nil
}

// readVMWithPowerShell reads VM properties using PowerShell
func readVMWithPowerShell(ctx context.Context, vmName string, inputs MachineInputs) (MachineOutputs, error) {
	logger := logging.GetLogger(ctx)
	outputs := MachineOutputs{MachineInputs: inputs}

	// Get VM details using PowerShell
	vmDetails := fmt.Sprintf(`
		$vm = Get-VM -Name "%s" -ErrorAction SilentlyContinue
		if ($vm) {
			$output = New-Object PSObject
			$output | Add-Member -Type NoteProperty -Name Name -Value $vm.Name
			$output | Add-Member -Type NoteProperty -Name ProcessorCount -Value $vm.ProcessorCount
			$output | Add-Member -Type NoteProperty -Name MemoryStartupBytes -Value $vm.MemoryStartup
			$output | Add-Member -Type NoteProperty -Name MemoryMinimumBytes -Value $vm.MemoryMinimum
			$output | Add-Member -Type NoteProperty -Name MemoryMaximumBytes -Value $vm.MemoryMaximum
			$output | Add-Member -Type NoteProperty -Name DynamicMemoryEnabled -Value $vm.DynamicMemoryEnabled
			$output | Add-Member -Type NoteProperty -Name Generation -Value $vm.Generation
			$output | Add-Member -Type NoteProperty -Name AutomaticStartAction -Value $vm.AutomaticStartAction
			$output | Add-Member -Type NoteProperty -Name AutomaticStopAction -Value $vm.AutomaticStopAction
			$output | ConvertTo-Json
		}
	`, vmName)

	output, err := util.RunPowerShellCommand(vmDetails)
	if err != nil {
		logger.Warnf("Failed to get VM details: %v", err)
		return outputs, nil
	}

	// Check if output is empty or invalid
	output = strings.TrimSpace(output)
	if output == "" {
		logger.Warnf("Empty response from PowerShell when getting VM details")
		return outputs, nil
	}

	// Parse key properties from output if present
	// We don't need to parse everything as we'll populate with input values

	// Try to extract processor count if not in inputs
	if inputs.ProcessorCount == nil {
		procCountMatch := strings.Index(output, `"ProcessorCount":`)
		if procCountMatch >= 0 {
			// This is a very simple parsing approach - in production you'd use proper JSON parsing
			// However, this illustrates the idea
			startIdx := procCountMatch + len(`"ProcessorCount":`)
			endIdx := startIdx
			for endIdx < len(output) && (output[endIdx] == ' ' || output[endIdx] == '\t') {
				endIdx++
			}
			for endIdx < len(output) && (output[endIdx] >= '0' && output[endIdx] <= '9') {
				endIdx++
			}

			if endIdx > startIdx {
				procCount, err := strconv.Atoi(strings.TrimSpace(output[startIdx:endIdx]))
				if err == nil {
					outputs.ProcessorCount = &procCount
					logger.Debugf("Found processor count: %d", procCount)
				}
			}
		}
	}

	// Try to extract memory if not in inputs
	if inputs.MemorySize == nil {
		memMatch := strings.Index(output, `"MemoryStartupBytes":`)
		if memMatch >= 0 {
			// This is a very simple parsing approach - in production you'd use proper JSON parsing
			startIdx := memMatch + len(`"MemoryStartupBytes":`)
			endIdx := startIdx
			for endIdx < len(output) && (output[endIdx] == ' ' || output[endIdx] == '\t') {
				endIdx++
			}
			for endIdx < len(output) && ((output[endIdx] >= '0' && output[endIdx] <= '9') || output[endIdx] == ',') {
				endIdx++
			}

			if endIdx > startIdx {
				memBytes, err := strconv.ParseInt(strings.TrimSpace(strings.Replace(output[startIdx:endIdx], ",", "", -1)), 10, 64)
				if err == nil {
					memMB := memBytes / (1024 * 1024) // Convert bytes to MB
					memMBInt := int(memMB)
					outputs.MemorySize = &memMBInt
					logger.Debugf("Found memory size: %d MB", memMB)
				}
			}
		}
	}

	// Get hard drives using PowerShell
	if len(inputs.HardDrives) == 0 {
		hddCmd := fmt.Sprintf(`
			$vm = Get-VM -Name "%s" -ErrorAction SilentlyContinue
			if ($vm) {
				$hds = Get-VMHardDiskDrive -VM $vm
				$hds | ForEach-Object { $_.Path }
			}
		`, vmName)

		hdOutput, err := util.RunPowerShellCommand(hddCmd)
		if err == nil && len(hdOutput) > 0 {
			// Process each line as a hard drive path
			hdPaths := strings.Split(strings.TrimSpace(hdOutput), "\n")
			hardDrives := make([]HardDriveInput, 0, len(hdPaths))

			for i, path := range hdPaths {
				trimmedPath := strings.TrimSpace(path)
				if trimmedPath != "" {
					hdPath := trimmedPath
					controllerType := "SCSI" // Default
					controllerNumber := 0    // Default
					controllerLocation := i

					hardDrive := HardDriveInput{
						Path:               &hdPath,
						ControllerType:     &controllerType,
						ControllerNumber:   &controllerNumber,
						ControllerLocation: &controllerLocation,
					}

					hardDrives = append(hardDrives, hardDrive)
				}
			}

			if len(hardDrives) > 0 {
				// Convert to pointer slice
				hardDrivePtrs := make([]*HardDriveInput, len(hardDrives))
				for i := range hardDrives {
					hardDrivePtrs[i] = &hardDrives[i]
				}
				outputs.HardDrives = hardDrivePtrs
				logger.Debugf("Found %d hard drives", len(hardDrives))
			}
		}
	}

	// Get network adapters using PowerShell
	if len(inputs.NetworkAdapters) == 0 {
		naCmd := fmt.Sprintf(`
			$vm = Get-VM -Name "%s" -ErrorAction SilentlyContinue
			if ($vm) {
				$adapters = Get-VMNetworkAdapter -VM $vm
				$adapters | ForEach-Object { 
					$adapterObj = New-Object PSObject
					$adapterObj | Add-Member -Type NoteProperty -Name Name -Value $_.Name
					$adapterObj | Add-Member -Type NoteProperty -Name SwitchName -Value $_.SwitchName
					$adapterObj | Add-Member -Type NoteProperty -Name MacAddress -Value $_.MacAddress
					$adapterObj | ConvertTo-Json
				}
			}
		`, vmName)

		naOutput, err := util.RunPowerShellCommand(naCmd)
		if err == nil && len(naOutput) > 0 {
			// Process the output - in a real implementation, you'd parse the JSON properly
			// For this example, we'll use a simple approach to illustrate

			// Check if we have multiple adapters (multiple JSON objects) or just one
			adapters := make([]NetworkAdapterInput, 0)

			if strings.Contains(naOutput, "SwitchName") {
				// For simplicity in this example, we'll just extract switch names
				naLines := strings.Split(naOutput, "\n")
				for _, line := range naLines {
					if strings.Contains(line, `"SwitchName":`) {
						startIdx := strings.Index(line, `"SwitchName":`) + len(`"SwitchName":`)
						endIdx := strings.Index(line[startIdx:], `"`)
						if endIdx > 0 {
							switchName := strings.TrimSpace(line[startIdx : startIdx+endIdx])
							switchName = strings.Trim(switchName, `"`)
							if switchName != "" {
								adapter := NetworkAdapterInput{
									SwitchName: &switchName,
								}
								adapters = append(adapters, adapter)
							}
						}
					}
				}

				if len(adapters) > 0 {
					// Convert NetworkAdapterInput to networkadapter.NetworkAdapterInputs
					naInputs := make([]*networkadapter.NetworkAdapterInputs, len(adapters))
					for i, adapter := range adapters {
						naInputs[i] = &networkadapter.NetworkAdapterInputs{
							SwitchName: adapter.SwitchName,
						}
					}
					outputs.NetworkAdapters = naInputs
					logger.Debugf("Found %d network adapters", len(adapters))
				}
			}
		}
	}

	return outputs, nil
}

// Delete method to delete a virtual machine
func (c *Machine) Delete(ctx context.Context, id string, props MachineOutputs) error {
	logger := logging.GetLogger(ctx)

	// Ensure vmId is set for the delete operation
	// If the VM ID is not available in props, use the resource ID
	vmName := id
	if props.VmId != nil {
		vmName = *props.VmId
	} else if props.MachineName != nil {
		vmName = *props.MachineName
	}

	logger.Infof("Deleting VM %s", vmName)

	// Try to stop the VM before deleting it
	// We'll use PowerShell as it's most reliable

	// Try to stop the VM using PowerShell - this is most reliable

	// First check if the VM exists
	existsCmd := fmt.Sprintf("Get-VM -Name \"%s\" -ErrorAction SilentlyContinue", vmName)
	existsOutput, existsErr := util.RunPowerShellCommand(existsCmd)
	if existsErr != nil || strings.TrimSpace(existsOutput) == "" {
		logger.Infof("VM %s does not exist or is not accessible, skipping deletion", vmName)
		return nil
	}

	// Check if the VM is running
	checkCmd := fmt.Sprintf("(Get-VM -Name \"%s\" -ErrorAction SilentlyContinue).State -eq 'Running'", vmName)
	output, err := util.RunPowerShellCommand(checkCmd)
	isRunning := false
	if err == nil && strings.TrimSpace(output) == "True" {
		isRunning = true
	}

	if isRunning {
		logger.Infof("Stopping VM %s before deletion", vmName)
		// Use -TurnOff to force an immediate shutdown rather than a graceful one
		stopCmd := fmt.Sprintf("Stop-VM -Name \"%s\" -Force -TurnOff", vmName)
		_, err = util.RunPowerShellCommand(stopCmd)
		if err != nil {
			logger.Warnf("Failed to stop VM with PowerShell: %v", err)
			// Try to delete anyway
		} else {
			logger.Infof("Successfully stopped VM %s", vmName)
		}
	}

	// Get OS version to check if we're running on Azure
	osVersion, osErr := util.GetOSVersion()
	if osErr == nil && strings.Contains(strings.ToLower(osVersion), "azure") {
		logger.Infof("Detected Azure datacenter edition, using alternative VM deletion approach")

		// On Azure, first check VM state again to ensure it's fully stopped
		checkStoppedCmd := fmt.Sprintf("(Get-VM -Name \"%s\" -ErrorAction SilentlyContinue).State -eq 'Off'", vmName)
		stoppedOutput, stoppedErr := util.RunPowerShellCommand(checkStoppedCmd)
		isStopped := false
		if stoppedErr == nil && strings.TrimSpace(stoppedOutput) == "True" {
			isStopped = true
		}

		if !isStopped {
			logger.Warnf("VM %s may not be fully stopped, force stopping again", vmName)
			stopAgainCmd := fmt.Sprintf("Stop-VM -Name \"%s\" -Force -TurnOff", vmName)
			_, _ = util.RunPowerShellCommand(stopAgainCmd)
			// Small delay to allow VM to fully stop
			logger.Infof("Waiting for VM to fully stop...")
		}

		// Try using the alternative deletion approach that works better on Azure
		// Uses Get-VM | Remove-VM pattern which can be more reliable than direct Remove-VM
		deleteAzureCmd := fmt.Sprintf("Get-VM -Name \"%s\" | Remove-VM -Force", vmName)
		_, azureErr := util.RunPowerShellCommand(deleteAzureCmd)
		if azureErr == nil {
			logger.Infof("Successfully deleted VM %s using Azure-specific approach", vmName)
			return nil
		}
		logger.Warnf("Alternative VM deletion failed: %v, trying standard approach", azureErr)
	}

	// Standard deletion approach
	deleteCmd := fmt.Sprintf("Remove-VM -Name \"%s\" -Force", vmName)
	_, err = util.RunPowerShellCommand(deleteCmd)
	if err != nil {
		// Check if VM still exists after failed deletion attempt
		existsCmd := fmt.Sprintf("Get-VM -Name \"%s\" -ErrorAction SilentlyContinue", vmName)
		existsOutput, _ := util.RunPowerShellCommand(existsCmd)
		if strings.TrimSpace(existsOutput) == "" {
			// VM doesn't exist anymore despite error, consider it successfully deleted
			logger.Infof("VM %s no longer exists despite deletion error, considering it successfully deleted", vmName)
			return nil
		}
		return fmt.Errorf("failed to delete VM %s with PowerShell: %v", vmName, err)
	}

	logger.Infof("Successfully deleted VM %s", vmName)
	return nil
}

func (c *Machine) Create(ctx context.Context, name string, input MachineInputs, preview bool) (string, MachineOutputs, error) {
	logger := logging.GetLogger(ctx)
	id := name
	if input.MachineName != nil {
		id = *input.MachineName
	}
	state := MachineOutputs{MachineInputs: input}

	// Always ensure vmId is set
	EnsureVmId(&state, id)

	// If in preview, don't run the command.
	if preview {
		return id, state, nil
	}
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		logger.Warnf("Error connecting to Hyper-V: %v", err)
		return id, state, fmt.Errorf("failed to connect to Hyper-V: %v", err)
	}

	// If both vmmsClient and vsms are nil, use PowerShell fallback entirely
	// If only vsms is nil but vmmsClient is available, use PowerShell for operations that need vsms
	if vmmsClient == nil && vsms == nil {
		// Using PowerShell fallback for the entire VM creation
		logger.Infof("Using PowerShell fallback to create VM %s", id)
		return createVMWithPowerShell(ctx, id, input)
	}

	// If we have vmmsClient but not vsms (common on Windows 10/11), also use PowerShell
	if vmmsClient != nil && vsms == nil {
		logger.Infof("VMMS available but VSMS is nil - using PowerShell fallback to create VM %s", id)
		return createVMWithPowerShell(ctx, id, input)
	}

	// If we have both vmmsClient and vsms, proceed with WMI implementation
	vConn := vmmsClient.GetVirtualizationConn()
	if vConn == nil {
		logger.Warnf("Virtualization connection is nil, falling back to PowerShell")
		return createVMWithPowerShell(ctx, id, input)
	}

	wmiHost := vConn.WMIHost
	if wmiHost == nil {
		logger.Warnf("WMI host is nil, falling back to PowerShell")
		return createVMWithPowerShell(ctx, id, input)
	}

	// Now create the VM settings with proper error handling
	var setting *virtualsystem.VirtualSystemSettingData
	var settingErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				settingErr = fmt.Errorf("recovered from panic in GetVirtualSystemSettingData: %v", r)
				logger.Warnf("Recovered from panic in GetVirtualSystemSettingData: %v", r)
			}
		}()

		setting, settingErr = virtualsystem.GetVirtualSystemSettingData(wmiHost, id)
	}()

	if settingErr != nil {
		logger.Warnf("Failed to get virtual system setting data: %v, falling back to PowerShell", settingErr)
		return createVMWithPowerShell(ctx, id, input)
	}

	if setting == nil {
		logger.Warnf("Virtual system setting data is nil, falling back to PowerShell")
		return createVMWithPowerShell(ctx, id, input)
	}

	err = setting.SetPropertyInstanceID(id)
	if err != nil {
		logger.Warnf("Failed to set property instance ID: %v, falling back to PowerShell", err)
		return createVMWithPowerShell(ctx, id, input)
	}

	defer setting.Close()
	logger.Debugf("Create VMSettings")

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
	logger.Debugf("Created VM")

	// Add hard drives if specified
	if len(input.HardDrives) > 0 {
		for _, hd := range input.HardDrives {
			if hd.Path == nil {
				logger.Debugf("Hard drive path not specified, skipping")
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

			logger.Debugf("Adding hard drive %s to VM %s", *hd.Path, id)
			logger.Debugf("Controller details for VM %s: type=%s, number=%d, location=%d",
				id, controllerType, controllerNumber, controllerLocation)

			// Wrap in recovery block to prevent panics in case of type errors
			var addErr error
			func() {
				defer func() {
					if r := recover(); r != nil {
						addErr = fmt.Errorf("recovered from panic in AttachVirtualHardDisk: %v", r)
						logger.Errorf("Recovered from panic in AttachVirtualHardDisk: %v", r)
					}
				}()

				// Use the new robust VMMS method with fallback logic
				addErr = vmmsClient.AttachVirtualHardDisk(vm, *hd.Path, controllerType, controllerNumber, controllerLocation, logger)
			}()

			if addErr != nil {
				logger.Errorf("Failed to add hard drive: %v", addErr)
				logger.Warnf("Saving machine to state despite hard drive attachment failure")
				return id, state, fmt.Errorf("failed to add hard drive: %v", addErr)
			} else {
				logger.Infof("Successfully added hard drive to VM %s", id)
			}
		}
	}

	// Add network adapters if specified
	if len(input.NetworkAdapters) > 0 {
		for i, na := range input.NetworkAdapters {
			if na.SwitchName == nil {
				logger.Debugf("Network adapter switch name not specified, skipping")
				continue
			}

			// Use the index as part of name if no name provided
			adapterName := fmt.Sprintf("Network Adapter %d", i+1)
			if na.Name != nil {
				adapterName = *na.Name
			}

			logger.Debugf("Adding network adapter %s to VM %s, connected to switch %s",
				adapterName, id, *na.SwitchName)

			// Wrap in recovery block to prevent panics in case of type errors
			var addErr error
			func() {
				defer func() {
					if r := recover(); r != nil {
						addErr = fmt.Errorf("recovered from panic in AddVirtualNetworkAdapterAndConnect: %v", r)
						logger.Errorf("Recovered from panic in AddVirtualNetworkAdapterAndConnect: %v", r)
					}
				}()

				// Use the new robust VMMS method with fallback logic
				addErr = vmmsClient.AddVirtualNetworkAdapterAndConnect(vm, adapterName, *na.SwitchName, logger)
			}()

			if addErr != nil {
				logger.Errorf("Failed to add network adapter: %v", addErr)
				// Return the state with created VM even though adding the network adapter failed
				// This allows the VM to exist in the Pulumi state so it can be cleaned up properly
				logger.Warnf("Saving machine to state despite network adapter attachment failure")
				return id, state, fmt.Errorf("failed to add network adapter: %v", addErr)
			} else {
				logger.Infof("Successfully added network adapter %s to VM %s", adapterName, id)
			}
		}
	}

	// Start the VM after all configuration is done
	logger.Infof("Starting VM %s", id)

	startCmd := fmt.Sprintf("Start-VM -Name \"%s\"", id)
	output, startErr := util.RunPowerShellCommand(startCmd)
	if startErr != nil {
		// Check for specific error conditions
		if strings.Contains(output, "Not enough memory in the system to start the virtual machine") {
			// Handle out of memory error with specific guidance
			memoryRequiredStr := fmt.Sprintf("%d", memorySizeMB)
			if input.MemorySize != nil {
				memoryRequiredStr = fmt.Sprintf("%d", *input.MemorySize)
			}

			logger.Errorf("Memory error starting VM %s: %s", id, output)
			return id, state, fmt.Errorf("failed to start VM due to insufficient memory: the system does not have enough memory to allocate %s MB for this VM. "+
				"Try reducing the memory allocation, closing other applications, or adding more RAM to the host system", memoryRequiredStr)
		} else if strings.Contains(output, "0x8007000E") {
			// Generic out of resources error (could be memory or something else)
			return id, state, fmt.Errorf("failed to start VM due to insufficient system resources (error 0x8007000E). " +
				"Try reducing VM resource allocation, closing other applications, or adding more resources to the host system")
		} else if strings.Contains(output, "could not initialize memory") {
			// Memory initialization error
			return id, state, fmt.Errorf("failed to start VM due to memory initialization error. " +
				"This could be due to insufficient memory, memory fragmentation, or a system configuration issue")
		} else {
			// Log the error but continue since the VM is created
			logger.Warnf("Failed to start VM with PowerShell: %v", startErr)
			logger.Debugf("Start-VM output: %s", output)
		}
	} else {
		logger.Infof("Started VM %s using PowerShell", id)
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
	logger := logging.GetLogger(ctx)
	logger.Infof("Updating VM %s", id)

	// Initialize the output state with the new inputs
	state := MachineOutputs{MachineInputs: news}

	// Always ensure vmId is set - carry over from old state if available, otherwise use id
	if olds.VmId != nil {
		state.VmId = olds.VmId
	} else {
		EnsureVmId(&state, id)
	}

	// If in preview, don't run the command.
	if preview {
		return state, nil
	}

	// Get the VM name from id, olds, or news
	vmName := id
	if olds.MachineName != nil {
		vmName = *olds.MachineName
	} else if news.MachineName != nil {
		vmName = *news.MachineName
	}
	logger.Infof("Using VM name: %s", vmName)

	// Check if the VM exists
	exists, err := checkVMExistsPowerShell(vmName)
	if err != nil {
		logger.Errorf("Error checking if VM exists: %v", err)
		return state, fmt.Errorf("failed to check if VM exists: %v", err)
	}
	if !exists {
		logger.Errorf("VM %s does not exist", vmName)
		return state, fmt.Errorf("VM %s does not exist", vmName)
	}

	// Connect to Hyper-V
	vmmsClient, vsms, err := c.Connect(ctx)
	if err != nil {
		logger.Warnf("Error connecting to Hyper-V: %v", err)
		// Fall back to PowerShell completely
		logger.Infof("Using PowerShell fallback for VM update")
		return updateVMWithPowerShell(ctx, vmName, olds, news)
	}

	// Check if we need to stop the VM to make changes
	// Always check if VM is running with PowerShell as it's most reliable
	needsRestart := false
	wasRunning, err := isVMRunningPowerShell(vmName)
	if err != nil {
		logger.Warnf("Error checking if VM is running: %v", err)
		// Continue anyway, we'll handle errors later
	}

	// Determine if we need to stop the VM based on properties being changed
	needsVMStopped := false

	// Check changes that require the VM to be stopped
	if (olds.ProcessorCount != nil && news.ProcessorCount != nil && *olds.ProcessorCount != *news.ProcessorCount) ||
		(olds.MemorySize != nil && news.MemorySize != nil && *olds.MemorySize != *news.MemorySize) ||
		((olds.DynamicMemory == nil || !*olds.DynamicMemory) && news.DynamicMemory != nil && *news.DynamicMemory) ||
		(olds.DynamicMemory != nil && *olds.DynamicMemory && news.DynamicMemory != nil && !*news.DynamicMemory) ||
		// Changes in minimum or maximum memory for dynamic memory
		((olds.MinimumMemory == nil && news.MinimumMemory != nil) ||
			(olds.MinimumMemory != nil && news.MinimumMemory != nil && *olds.MinimumMemory != *news.MinimumMemory)) ||
		((olds.MaximumMemory == nil && news.MaximumMemory != nil) ||
			(olds.MaximumMemory != nil && news.MaximumMemory != nil && *olds.MaximumMemory != *news.MaximumMemory)) {
		needsVMStopped = true
		logger.Infof("VM update requires stopping the VM because processor, memory, or dynamic memory settings are changing")
	}

	// Network adapters or hard drives changes also require VM to be stopped
	if len(olds.NetworkAdapters) != len(news.NetworkAdapters) || len(olds.HardDrives) != len(news.HardDrives) {
		needsVMStopped = true
		logger.Infof("VM update requires stopping the VM because network adapters or hard drives are changing")
	}

	// If VM needs to be stopped and is running, stop it
	if needsVMStopped && wasRunning {
		logger.Infof("Stopping VM %s before updating", vmName)
		stopErr := stopVMPowerShell(vmName)
		if stopErr != nil {
			logger.Errorf("Failed to stop VM %s: %v", vmName, stopErr)
			return state, fmt.Errorf("failed to stop VM %s before update: %v", vmName, stopErr)
		}
		needsRestart = true
		logger.Infof("VM %s stopped successfully", vmName)
	}

	// If we don't have VMMS client or VSMS, use PowerShell for everything
	if vmmsClient == nil || vsms == nil {
		logger.Infof("Using PowerShell fallback for VM update because VMMS or VSMS is nil")
		result, err := updateVMWithPowerShell(ctx, vmName, olds, news)

		// Start the VM if we stopped it and update was successful
		if err == nil && needsRestart {
			logger.Infof("Restarting VM %s after update", vmName)
			startErr := startVMPowerShell(vmName)
			if startErr != nil {
				logger.Warnf("Failed to restart VM %s after update: %v", vmName, startErr)
				// We'll warn but not fail the update since the changes were applied
			}
		}

		return result, err
	}

	// Use WMI when available
	vm, err := vsms.GetVirtualMachineByName(vmName)
	if err != nil {
		logger.Warnf("Failed to get VM %s using WMI: %v", vmName, err)
		logger.Infof("Falling back to PowerShell for VM update")
		result, err := updateVMWithPowerShell(ctx, vmName, olds, news)

		// Start the VM if we stopped it and update was successful
		if err == nil && needsRestart {
			logger.Infof("Restarting VM %s after update", vmName)
			startErr := startVMPowerShell(vmName)
			if startErr != nil {
				logger.Warnf("Failed to restart VM %s after update: %v", vmName, startErr)
				// We'll warn but not fail the update since the changes were applied
			}
		}

		return result, err
	}
	defer vm.Close()

	// Get VM settings data
	vmSettings, err := virtualsystem.GetVirtualSystemSettingData(vmmsClient.GetVirtualizationConn().WMIHost, vmName)
	if err != nil {
		logger.Warnf("Failed to get VM settings: %v", err)
		logger.Infof("Falling back to PowerShell for VM update")
		result, err := updateVMWithPowerShell(ctx, vmName, olds, news)

		// Start the VM if we stopped it and update was successful
		if err == nil && needsRestart {
			logger.Infof("Restarting VM %s after update", vmName)
			startErr := startVMPowerShell(vmName)
			if startErr != nil {
				logger.Warnf("Failed to restart VM %s after update: %v", vmName, startErr)
				// We'll warn but not fail the update since the changes were applied
			}
		}

		return result, err
	}
	defer vmSettings.Close()

	// Update processor count if changed
	if news.ProcessorCount != nil && (olds.ProcessorCount == nil || *olds.ProcessorCount != *news.ProcessorCount) {
		logger.Infof("Updating processor count from %v to %d", olds.ProcessorCount, *news.ProcessorCount)

		// Try WMI first
		processorSettings, err := processor.GetProcessorSettingDataFromVirtualSystemSettingData(vmmsClient.GetVirtualizationConn().WMIHost, vmSettings)
		if err != nil || processorSettings == nil {
			logger.Warnf("Failed to get processor settings: %v", err)
			// Fallback to PowerShell for this setting
			procCmd := fmt.Sprintf("Set-VMProcessor -VMName \"%s\" -Count %d", vmName, *news.ProcessorCount)
			_, psErr := util.RunPowerShellCommand(procCmd)
			if psErr != nil {
				logger.Warnf("Failed to update processor count: %v", psErr)
				// Continue with other updates despite error
			} else {
				logger.Infof("Updated processor count to %d using PowerShell", *news.ProcessorCount)
			}
		} else {
			defer processorSettings.Close()
			err = processorSettings.SetCPUCount(uint64(*news.ProcessorCount))
			if err != nil {
				logger.Warnf("Failed to set CPU count: %v", err)
				// Fallback to PowerShell
				procCmd := fmt.Sprintf("Set-VMProcessor -VMName \"%s\" -Count %d", vmName, *news.ProcessorCount)
				_, psErr := util.RunPowerShellCommand(procCmd)
				if psErr != nil {
					logger.Warnf("Failed to update processor count with PowerShell fallback: %v", psErr)
					// Continue with other updates despite error
				} else {
					logger.Infof("Updated processor count to %d using PowerShell fallback", *news.ProcessorCount)
				}
			} else {
				logger.Infof("Updated processor count to %d using WMI", *news.ProcessorCount)
			}
		}
	}

	// Update memory settings if changed
	memoryChanged := false
	if news.MemorySize != nil && (olds.MemorySize == nil || *olds.MemorySize != *news.MemorySize) {
		memoryChanged = true
	}

	dynamicMemoryChanged := false
	if (news.DynamicMemory != nil && olds.DynamicMemory == nil) ||
		(news.DynamicMemory != nil && olds.DynamicMemory != nil && *news.DynamicMemory != *olds.DynamicMemory) {
		dynamicMemoryChanged = true
	}

	minMemoryChanged := false
	if news.MinimumMemory != nil && (olds.MinimumMemory == nil || *olds.MinimumMemory != *news.MinimumMemory) {
		minMemoryChanged = true
	}

	maxMemoryChanged := false
	if news.MaximumMemory != nil && (olds.MaximumMemory == nil || *olds.MaximumMemory != *news.MaximumMemory) {
		maxMemoryChanged = true
	}

	if memoryChanged || dynamicMemoryChanged || minMemoryChanged || maxMemoryChanged {
		logger.Infof("Updating memory settings for VM %s", vmName)

		// Try WMI first
		memorySettings, err := memory.GetMemorySettingDataFromVirtualSystemSettingData(vmmsClient.GetVirtualizationConn().WMIHost, vmSettings)
		if err != nil || memorySettings == nil {
			logger.Warnf("Failed to get memory settings: %v", err)
			// Fallback to PowerShell for all memory settings
			updateMemoryWithPowerShell(ctx, vmName, news)
		} else {
			defer memorySettings.Close()

			// Update memory size if changed
			if memoryChanged {
				err = memorySettings.SetSizeMB(uint64(*news.MemorySize))
				if err != nil {
					logger.Warnf("Failed to set memory size: %v", err)
					// Will fall back to PowerShell below
				} else {
					logger.Infof("Updated memory size to %d MB using WMI", *news.MemorySize)
				}
			}

			// Update dynamic memory settings if changed
			if dynamicMemoryChanged {
				err = memorySettings.SetPropertyDynamicMemoryEnabled(*news.DynamicMemory)
				if err != nil {
					logger.Warnf("Failed to set dynamic memory enabled: %v", err)
					// Will fall back to PowerShell below
				} else {
					logger.Infof("Updated dynamic memory to %v using WMI", *news.DynamicMemory)
				}
			}

			// Update minimum memory if changed
			if minMemoryChanged && news.MinimumMemory != nil {
				err = memorySettings.SetProperty("MinimumBytes", uint64(*news.MinimumMemory)*1024*1024) // Convert MB to bytes
				if err != nil {
					logger.Warnf("Failed to set minimum memory: %v", err)
					// Will fall back to PowerShell below
				} else {
					logger.Infof("Updated minimum memory to %d MB using WMI", *news.MinimumMemory)
				}
			}

			// Update maximum memory if changed
			if maxMemoryChanged && news.MaximumMemory != nil {
				err = memorySettings.SetProperty("MaximumBytes", uint64(*news.MaximumMemory)*1024*1024) // Convert MB to bytes
				if err != nil {
					logger.Warnf("Failed to set maximum memory: %v", err)
					// Will fall back to PowerShell below
				} else {
					logger.Infof("Updated maximum memory to %d MB using WMI", *news.MaximumMemory)
				}
			}

			// If any memory update failed, use PowerShell as fallback
			if err != nil {
				logger.Infof("Using PowerShell fallback for memory settings")
				updateMemoryWithPowerShell(ctx, vmName, news)
			}
		}
	}

	// Update auto start action if changed
	if news.AutoStartAction != nil && (olds.AutoStartAction == nil || *olds.AutoStartAction != *news.AutoStartAction) {
		logger.Infof("Updating auto start action from %v to %s", olds.AutoStartAction, *news.AutoStartAction)

		// Try WMI first
		var autoStartValue uint16
		switch *news.AutoStartAction {
		case "Nothing":
			autoStartValue = 0
		case "StartIfRunning":
			autoStartValue = 1
		case "Start":
			autoStartValue = 2
		default:
			logger.Errorf("Invalid auto start action: %s, setting Nothing", *news.AutoStartAction)
			autoStartValue = 0
		}

		err = vmSettings.SetProperty("AutoStartAction", autoStartValue)
		if err != nil {
			logger.Warnf("Failed to set auto start action: %v", err)
			// Fallback to PowerShell
			autoStartCmd := fmt.Sprintf("Set-VM -VMName \"%s\" -AutomaticStartAction %s", vmName, *news.AutoStartAction)
			_, psErr := util.RunPowerShellCommand(autoStartCmd)
			if psErr != nil {
				logger.Warnf("Failed to update auto start action: %v", psErr)
				// Continue with other updates despite error
			} else {
				logger.Infof("Updated auto start action to %s using PowerShell", *news.AutoStartAction)
			}
		} else {
			logger.Infof("Updated auto start action to %s using WMI", *news.AutoStartAction)
		}
	}

	// Update auto stop action if changed
	if news.AutoStopAction != nil && (olds.AutoStopAction == nil || *olds.AutoStopAction != *news.AutoStopAction) {
		logger.Infof("Updating auto stop action from %v to %s", olds.AutoStopAction, *news.AutoStopAction)

		// Try WMI first
		var autoStopValue uint16
		switch *news.AutoStopAction {
		case "TurnOff":
			autoStopValue = 0
		case "Save":
			autoStopValue = 1
		case "ShutDown":
			autoStopValue = 2
		default:
			logger.Errorf("Invalid auto stop action: %s, setting TurnOff", *news.AutoStopAction)
			autoStopValue = 0
		}

		err = vmSettings.SetProperty("AutoStopAction", autoStopValue)
		if err != nil {
			logger.Warnf("Failed to set auto stop action: %v", err)
			// Fallback to PowerShell
			autoStopCmd := fmt.Sprintf("Set-VM -VMName \"%s\" -AutomaticStopAction %s", vmName, *news.AutoStopAction)
			_, psErr := util.RunPowerShellCommand(autoStopCmd)
			if psErr != nil {
				logger.Warnf("Failed to update auto stop action: %v", psErr)
				// Continue with other updates despite error
			} else {
				logger.Infof("Updated auto stop action to %s using PowerShell", *news.AutoStopAction)
			}
		} else {
			logger.Infof("Updated auto stop action to %s using WMI", *news.AutoStopAction)
		}
	}

	// Update hard drives if changed
	if !compareHardDrives(olds.HardDrives, news.HardDrives) {
		logger.Infof("Updating hard drives for VM %s", vmName)

		// First remove all existing hard drives using PowerShell (more reliable)
		removeHDCmd := fmt.Sprintf("Get-VMHardDiskDrive -VMName \"%s\" | Remove-VMHardDiskDrive", vmName)
		_, removeErr := util.RunPowerShellCommand(removeHDCmd)
		if removeErr != nil {
			logger.Warnf("Failed to remove existing hard drives: %v", removeErr)
			// Try to continue anyway
		}

		// Add all new hard drives
		for i, hd := range news.HardDrives {
			if hd.Path == nil {
				logger.Debugf("Hard drive path not specified, skipping")
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

			// Default to sequential ports
			controllerLocation := i
			if hd.ControllerLocation != nil {
				controllerLocation = *hd.ControllerLocation
			}

			// Use PowerShell to add hard drive (more reliable)
			hdCmd := fmt.Sprintf("Add-VMHardDiskDrive -VMName \"%s\" -Path \"%s\" -ControllerType %s -ControllerNumber %d -ControllerLocation %d",
				vmName, *hd.Path, controllerType, controllerNumber, controllerLocation)
			_, err := util.RunPowerShellCommand(hdCmd)
			if err != nil {
				logger.Warnf("Failed to add hard drive: %v", err)
				// Continue with other hard drives despite error
			} else {
				logger.Debugf("Added hard drive %s to VM %s", *hd.Path, vmName)
			}
		}
	}

	// Update network adapters if changed
	if !compareNetworkAdapters(olds.NetworkAdapters, news.NetworkAdapters) {
		logger.Infof("Updating network adapters for VM %s", vmName)

		// First remove all existing network adapters using PowerShell (more reliable)
		removeNACmd := fmt.Sprintf("Get-VMNetworkAdapter -VMName \"%s\" | Remove-VMNetworkAdapter", vmName)
		_, removeErr := util.RunPowerShellCommand(removeNACmd)
		if removeErr != nil {
			logger.Warnf("Failed to remove existing network adapters: %v", removeErr)
			// Try to continue anyway
		}

		// Add all new network adapters
		for i, na := range news.NetworkAdapters {
			if na.SwitchName == nil {
				logger.Debugf("Network adapter switch name not specified, skipping")
				continue
			}

			// Use the index as part of name if no name provided
			adapterName := fmt.Sprintf("Network Adapter %d", i+1)
			if na.Name != nil {
				adapterName = *na.Name
			}

			// Create the adapter and connect it to the switch
			naCmd := fmt.Sprintf("Add-VMNetworkAdapter -VMName \"%s\" -Name \"%s\" -SwitchName \"%s\"",
				vmName, adapterName, *na.SwitchName)

			// Add MAC address if specified
			if na.MacAddress != nil && *na.MacAddress != "" {
				naCmd += fmt.Sprintf(" -StaticMacAddress \"%s\"", *na.MacAddress)
			}

			_, err := util.RunPowerShellCommand(naCmd)
			if err != nil {
				logger.Warnf("Failed to add network adapter: %v", err)
				// Continue with other adapters despite error
			} else {
				logger.Debugf("Added network adapter %s to VM %s", adapterName, vmName)
			}
		}
	}

	// Start the VM if we stopped it
	if needsRestart {
		logger.Infof("Restarting VM %s after update", vmName)
		startErr := startVMPowerShell(vmName)
		if startErr != nil {
			logger.Warnf("Failed to restart VM %s after update: %v", vmName, startErr)
			// We'll warn but not fail the update since the changes were applied
		}
	}

	return state, nil
}

// Helper functions for the Update method

// updateVMWithPowerShell updates a virtual machine using PowerShell cmdlets
func updateVMWithPowerShell(ctx context.Context, vmName string, olds MachineOutputs, news MachineInputs) (MachineOutputs, error) {
	logger := logging.GetLogger(ctx)
	state := MachineOutputs{MachineInputs: news}

	// Always ensure vmId is set - carry over from old state if available
	if olds.VmId != nil {
		state.VmId = olds.VmId
	} else {
		EnsureVmId(&state, vmName)
	}

	// Update processor count if changed
	if news.ProcessorCount != nil && (olds.ProcessorCount == nil || *olds.ProcessorCount != *news.ProcessorCount) {
		logger.Infof("Updating processor count from %v to %d", olds.ProcessorCount, *news.ProcessorCount)
		procCmd := fmt.Sprintf("Set-VMProcessor -VMName \"%s\" -Count %d", vmName, *news.ProcessorCount)
		_, err := util.RunPowerShellCommand(procCmd)
		if err != nil {
			logger.Warnf("Failed to update processor count: %v", err)
			// Continue with other updates despite error
		}
	}

	// Update memory settings
	updateMemoryWithPowerShell(ctx, vmName, news)

	// Update auto start action if changed
	if news.AutoStartAction != nil && (olds.AutoStartAction == nil || *olds.AutoStartAction != *news.AutoStartAction) {
		logger.Infof("Updating auto start action from %v to %s", olds.AutoStartAction, *news.AutoStartAction)
		autoStartCmd := fmt.Sprintf("Set-VM -VMName \"%s\" -AutomaticStartAction %s", vmName, *news.AutoStartAction)
		_, err := util.RunPowerShellCommand(autoStartCmd)
		if err != nil {
			logger.Warnf("Failed to update auto start action: %v", err)
			// Continue with other updates despite error
		}
	}

	// Update auto stop action if changed
	if news.AutoStopAction != nil && (olds.AutoStopAction == nil || *olds.AutoStopAction != *news.AutoStopAction) {
		logger.Infof("Updating auto stop action from %v to %s", olds.AutoStopAction, *news.AutoStopAction)
		autoStopCmd := fmt.Sprintf("Set-VM -VMName \"%s\" -AutomaticStopAction %s", vmName, *news.AutoStopAction)
		_, err := util.RunPowerShellCommand(autoStopCmd)
		if err != nil {
			logger.Warnf("Failed to update auto stop action: %v", err)
			// Continue with other updates despite error
		}
	}

	// Update hard drives if changed
	if !compareHardDrives(olds.HardDrives, news.HardDrives) {
		logger.Infof("Updating hard drives for VM %s", vmName)

		// First remove all existing hard drives
		removeHDCmd := fmt.Sprintf("Get-VMHardDiskDrive -VMName \"%s\" | Remove-VMHardDiskDrive", vmName)
		_, removeErr := util.RunPowerShellCommand(removeHDCmd)
		if removeErr != nil {
			logger.Warnf("Failed to remove existing hard drives: %v", removeErr)
			// Try to continue anyway
		}

		// Add all new hard drives
		for i, hd := range news.HardDrives {
			if hd.Path == nil {
				logger.Debugf("Hard drive path not specified, skipping")
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

			// Default to sequential ports
			controllerLocation := i
			if hd.ControllerLocation != nil {
				controllerLocation = *hd.ControllerLocation
			}

			hdCmd := fmt.Sprintf("Add-VMHardDiskDrive -VMName \"%s\" -Path \"%s\" -ControllerType %s -ControllerNumber %d -ControllerLocation %d",
				vmName, *hd.Path, controllerType, controllerNumber, controllerLocation)
			_, err := util.RunPowerShellCommand(hdCmd)
			if err != nil {
				logger.Warnf("Failed to add hard drive: %v", err)
				// Continue with other hard drives despite error
			} else {
				logger.Debugf("Added hard drive %s to VM %s", *hd.Path, vmName)
			}
		}
	}

	// Update network adapters if changed
	if !compareNetworkAdapters(olds.NetworkAdapters, news.NetworkAdapters) {
		logger.Infof("Updating network adapters for VM %s", vmName)

		// First remove all existing network adapters
		removeNACmd := fmt.Sprintf("Get-VMNetworkAdapter -VMName \"%s\" | Remove-VMNetworkAdapter", vmName)
		_, removeErr := util.RunPowerShellCommand(removeNACmd)
		if removeErr != nil {
			logger.Warnf("Failed to remove existing network adapters: %v", removeErr)
			// Try to continue anyway
		}

		// Add all new network adapters
		for i, na := range news.NetworkAdapters {
			if na.SwitchName == nil {
				logger.Debugf("Network adapter switch name not specified, skipping")
				continue
			}

			// Use the index as part of name if no name provided
			adapterName := fmt.Sprintf("Network Adapter %d", i+1)
			if na.Name != nil {
				adapterName = *na.Name
			}

			// Create the adapter and connect it to the switch
			naCmd := fmt.Sprintf("Add-VMNetworkAdapter -VMName \"%s\" -Name \"%s\" -SwitchName \"%s\"",
				vmName, adapterName, *na.SwitchName)

			// Add MAC address if specified
			if na.MacAddress != nil && *na.MacAddress != "" {
				naCmd += fmt.Sprintf(" -StaticMacAddress \"%s\"", *na.MacAddress)
			}

			_, err := util.RunPowerShellCommand(naCmd)
			if err != nil {
				logger.Warnf("Failed to add network adapter: %v", err)
				// Continue with other adapters despite error
			} else {
				logger.Debugf("Added network adapter %s to VM %s", adapterName, vmName)
			}
		}
	}

	return state, nil
}

// updateMemoryWithPowerShell updates memory settings using PowerShell
func updateMemoryWithPowerShell(ctx context.Context, vmName string, news MachineInputs) {
	logger := logging.GetLogger(ctx)

	// Set dynamic memory if specified
	if news.DynamicMemory != nil {
		var memCmd string

		if *news.DynamicMemory {
			// Enable dynamic memory
			var minMem, maxMem int64

			if news.MinimumMemory != nil {
				minMem = int64(*news.MinimumMemory) * 1024 * 1024 // Convert MB to bytes
			} else {
				minMem = 512 * 1024 * 1024 // 512MB default
			}

			if news.MaximumMemory != nil {
				maxMem = int64(*news.MaximumMemory) * 1024 * 1024 // Convert MB to bytes
			} else if news.MemorySize != nil {
				maxMem = int64(*news.MemorySize) * 2 * 1024 * 1024 // Double startup memory if maximum not specified
			} else {
				maxMem = 2 * 1024 * 1024 * 1024 // 2GB default
			}

			var startupMem int64
			if news.MemorySize != nil {
				startupMem = int64(*news.MemorySize) * 1024 * 1024 // Convert MB to bytes
			} else {
				startupMem = 1024 * 1024 * 1024 // 1GB default
			}

			memCmd = fmt.Sprintf("Set-VMMemory -VMName \"%s\" -DynamicMemoryEnabled $true -MinimumBytes %d -MaximumBytes %d -StartupBytes %d",
				vmName, minMem, maxMem, startupMem)
		} else {
			// Disable dynamic memory and set static memory
			var memorySize int64
			if news.MemorySize != nil {
				memorySize = int64(*news.MemorySize) * 1024 * 1024 // Convert MB to bytes
			} else {
				memorySize = 1024 * 1024 * 1024 // 1GB default
			}

			memCmd = fmt.Sprintf("Set-VMMemory -VMName \"%s\" -DynamicMemoryEnabled $false -StartupBytes %d",
				vmName, memorySize)
		}

		_, err := util.RunPowerShellCommand(memCmd)
		if err != nil {
			logger.Warnf("Failed to update memory settings: %v", err)
			// Continue despite error
		} else {
			logger.Infof("Updated memory settings for VM %s", vmName)
		}
	} else if news.MemorySize != nil {
		// Just update memory size without changing dynamic memory setting
		memorySize := int64(*news.MemorySize) * 1024 * 1024 // Convert MB to bytes
		memCmd := fmt.Sprintf("Set-VMMemory -VMName \"%s\" -StartupBytes %d", vmName, memorySize)

		_, err := util.RunPowerShellCommand(memCmd)
		if err != nil {
			logger.Warnf("Failed to update memory size: %v", err)
			// Continue despite error
		} else {
			logger.Infof("Updated memory size to %d MB for VM %s", *news.MemorySize, vmName)
		}
	}
}

// compareHardDrives compares two slices of hard drives to see if they're different
func compareHardDrives(olds []*HardDriveInput, news []*HardDriveInput) bool {
	if len(olds) != len(news) {
		return false
	}

	// Create maps for comparison
	oldMap := make(map[string]bool)
	newMap := make(map[string]bool)

	for _, hd := range olds {
		if hd.Path != nil {
			oldMap[*hd.Path] = true
		}
	}

	for _, hd := range news {
		if hd.Path != nil {
			newMap[*hd.Path] = true
			// Check if this path exists in old hard drives
			if !oldMap[*hd.Path] {
				return false
			}
		}
	}

	// Check if all old paths are in new hard drives
	for _, hd := range olds {
		if hd.Path != nil {
			if !newMap[*hd.Path] {
				return false
			}
		}
	}

	return true
}

// compareNetworkAdapters compares two slices of network adapters to see if they're different
func compareNetworkAdapters(olds []*networkadapter.NetworkAdapterInputs, news []*networkadapter.NetworkAdapterInputs) bool {
	if len(olds) != len(news) {
		return false
	}

	// Create maps for comparison
	oldMap := make(map[string]bool)
	newMap := make(map[string]bool)

	for _, na := range olds {
		if na.SwitchName != nil {
			key := *na.SwitchName
			if na.Name != nil {
				key = *na.Name + ":" + key
			}
			oldMap[key] = true
		}
	}

	for _, na := range news {
		if na.SwitchName != nil {
			key := *na.SwitchName
			if na.Name != nil {
				key = *na.Name + ":" + key
			}
			newMap[key] = true
			// Check if this combination exists in old adapters
			if !oldMap[key] {
				return false
			}
		}
	}

	// Check if all old adapters are in new adapters
	for _, na := range olds {
		if na.SwitchName != nil {
			key := *na.SwitchName
			if na.Name != nil {
				key = *na.Name + ":" + key
			}
			if !newMap[key] {
				return false
			}
		}
	}

	return true
}

// isVMRunningPowerShell checks if a VM is running using PowerShell
func isVMRunningPowerShell(vmName string) (bool, error) {
	checkCmd := fmt.Sprintf("(Get-VM -Name \"%s\" -ErrorAction SilentlyContinue).State -eq 'Running'", vmName)
	output, err := util.RunPowerShellCommand(checkCmd)
	if err != nil {
		return false, fmt.Errorf("failed to check VM state: %v", err)
	}

	return strings.TrimSpace(output) == "True", nil
}

// stopVMPowerShell stops a VM using PowerShell
func stopVMPowerShell(vmName string) error {
	// Check if the VM is running first
	isRunning, err := isVMRunningPowerShell(vmName)
	if err != nil {
		return fmt.Errorf("failed to check if VM is running: %v", err)
	}

	if !isRunning {
		return nil // VM is already stopped
	}

	stopCmd := fmt.Sprintf("Stop-VM -Name \"%s\" -Force", vmName)
	_, err = util.RunPowerShellCommand(stopCmd)
	if err != nil {
		return fmt.Errorf("failed to stop VM: %v", err)
	}

	// Verify the VM is stopped
	isStillRunning, verifyErr := isVMRunningPowerShell(vmName)
	if verifyErr != nil {
		return fmt.Errorf("failed to verify VM stopped state: %v", verifyErr)
	}

	if isStillRunning {
		return fmt.Errorf("VM did not stop after Stop-VM command")
	}

	return nil
}

// startVMPowerShell starts a VM using PowerShell
func startVMPowerShell(vmName string) error {
	// Check if the VM is already running
	isRunning, err := isVMRunningPowerShell(vmName)
	if err != nil {
		return fmt.Errorf("failed to check if VM is running: %v", err)
	}

	if isRunning {
		return nil // VM is already running
	}

	startCmd := fmt.Sprintf("Start-VM -Name \"%s\"", vmName)
	output, err := util.RunPowerShellCommand(startCmd)
	if err != nil {
		// Check for specific error conditions
		if strings.Contains(output, "Not enough memory in the system to start the virtual machine") {
			return fmt.Errorf("failed to start VM due to insufficient memory: %v", err)
		} else if strings.Contains(output, "0x8007000E") {
			return fmt.Errorf("failed to start VM due to insufficient system resources (error 0x8007000E)")
		} else if strings.Contains(output, "could not initialize memory") {
			return fmt.Errorf("failed to start VM due to memory initialization error")
		}

		return fmt.Errorf("failed to start VM: %v", err)
	}

	return nil
}
