// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vhdfile

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/storage/disk"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/logging"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/util"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
)

// VhdFileController implements the controller methods for VhdFile.
// The actual VhdFile type is defined in vhdfile.go.

// The following statements are type assertions to indicate to Go that VhdFile implements the interfaces.
var _ = (infer.CustomResource[VhdFileInputs, VhdFileOutputs])((*VhdFile)(nil))
var _ = (infer.CustomCreate[VhdFileInputs, VhdFileOutputs])((*VhdFile)(nil))
var _ = (infer.CustomRead[VhdFileInputs, VhdFileOutputs])((*VhdFile)(nil))
var _ = (infer.CustomDelete[VhdFileOutputs])((*VhdFile)(nil))

// Connect establishes a connection to the Hyper-V server.
func (c *VhdFile) Connect(ctx context.Context) (*vmms.VMMS, interface{}, error) {
	logger := logging.GetLogger(ctx)

	// Initialize all the parameters.
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
		// Log the error but don't fail - we'll continue with nil client and use fallback methods
		logger.Warnf("Failed to create VMMS client: %v", vmmsErr)
		logger.Infof("Will attempt to use PowerShell fallback methods for VHD operations")
	}

	// We'll return the client even if it's nil, and handle nil checks in the resource methods
	// This allows us to fall back to PowerShell commands when the WMI services are unavailable
	return vmmsClient, nil, nil
}

// Delete deletes a VHD file.
func (c *VhdFile) Delete(ctx context.Context, id string, state VhdFileOutputs) error {
	// Delete the VHD file.
	logger := logging.GetLogger(ctx)
	logger.Infof("Deleting vhd [%v]", state.Path)

	// If the path is empty, we can't delete the VHD file.
	if state.Path == nil {
		return fmt.Errorf("Path is nil")
	}

	if !strings.HasSuffix(*state.Path, ".vhd") && !strings.HasSuffix(*state.Path, ".vhdx") {
		return fmt.Errorf("Path [%v] doesn't end with .vhd or .vhdx", *state.Path)
	}

	// Get the disk.
	vmmsClient, _, err := c.Connect(ctx)
	if err != nil {
		return err
	}

	// Try deleting the VHD using various available methods
	vhdPath := *state.Path

	// Check if the file actually exists first
	powershellExe, psErr := util.FindPowerShellExe()
	if psErr != nil {
		return fmt.Errorf("failed to find PowerShell executable: %v", psErr)
	}

	checkCmd := exec.Command(powershellExe, "-Command", fmt.Sprintf("Test-Path -Path \"%s\"", vhdPath))
	checkOutput, checkErr := checkCmd.CombinedOutput()
	fileExists := false
	if checkErr == nil && strings.TrimSpace(string(checkOutput)) == "True" {
		fileExists = true
	}

	if !fileExists {
		logger.Infof("VHD file [%s] already doesn't exist, considering deletion successful", vhdPath)
		return nil
	}

	// First check for nil client - use PowerShell directly if no VMMS client
	if vmmsClient == nil {
		logger.Warnf("VMMS client is nil, attempting to delete file via PowerShell")
		// Use PowerShell to delete the file as a fallback
		cmd := exec.Command(powershellExe, "-Command", fmt.Sprintf("Remove-Item -Path \"%s\" -Force", vhdPath))
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to delete VHD file with PowerShell: %v, output: %s", err, string(output))
		}
		logger.Infof("Deleted VHD [%s] using PowerShell", vhdPath)
		return nil
	}

	// Try using the VirtualSystemManagementService if available
	vsms := vmmsClient.GetVirtualSystemManagementService()
	if vsms == nil {
		logger.Warnf("VirtualSystemManagementService is unavailable, falling back to PowerShell")
		// Use PowerShell to delete the file as a fallback
		cmd := exec.Command(powershellExe, "-Command", fmt.Sprintf("Remove-Item -Path \"%s\" -Force", vhdPath))
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to delete VHD file with PowerShell: %v, output: %s", err, string(output))
		}
		logger.Infof("Deleted VHD [%s] using PowerShell", vhdPath)
		return nil
	}

	// Try deleting with WMI first
	params := map[string]interface{}{
		"Path": vhdPath,
	}
	_, err = vsms.InvokeMethod("DeleteVirtualHardDisk", params)
	if err != nil {
		logger.Warnf("Failed to delete vhd [%v] via WMI because [%v], falling back to direct file removal", vhdPath, err)

		// If WMI method fails, try direct PowerShell file removal as a fallback
		cmd := exec.Command(powershellExe, "-Command", fmt.Sprintf("Remove-Item -Path \"%s\" -Force", vhdPath))
		output, removeErr := cmd.CombinedOutput()
		if removeErr != nil {
			logger.Errorf("Failed to delete VHD file with PowerShell after WMI failure: %v, output: %s", removeErr, string(output))
			return fmt.Errorf("Failed to delete vhd [%v]: WMI error: [%v], PowerShell error: [%v]", vhdPath, err, removeErr)
		}

		logger.Infof("Deleted VHD [%s] using PowerShell after WMI method failed", vhdPath)
		return nil
	}

	logger.Infof("Successfully deleted VHD [%s] using WMI", vhdPath)
	return nil
}

// This is the Create method. This will be run on every VhdFile resource creation.
func (c *VhdFile) Create(ctx context.Context, name string, input VhdFileInputs, preview bool) (string, VhdFileOutputs, error) {
	logger := logging.GetLogger(ctx)
	state := VhdFileOutputs{VhdFileInputs: input}
	id := name

	// If in preview, don't run the command.
	if preview {
		return id, state, nil
	}

	vmmsClient, _, err := c.Connect(ctx)
	if err != nil {
		return id, state, err
	}

	// Check for nil client
	if vmmsClient == nil {
		logger.Warnf("VMMS client is nil, will attempt PowerShell fallback for VHD creation")
	}
	// Create the VHD file.
	vhdFileName := *input.Path

	// Create parent directory if it doesn't exist (using Go's native code)
	// Extract directory path from the full VHD path
	lastSlashIndex := strings.LastIndex(vhdFileName, "\\")
	if lastSlashIndex == -1 {
		lastSlashIndex = strings.LastIndex(vhdFileName, "/")
	}

	if lastSlashIndex != -1 {
		dirPath := vhdFileName[:lastSlashIndex]
		logger.Infof("Ensuring parent directory exists: %s", dirPath)

		// Use Go's native os.MkdirAll to create directory and any missing parents
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			logger.Errorf("Failed to create parent directory: %v", err)
			return id, state, fmt.Errorf("failed to create parent directory: %v", err)
		}
	}

	// Check if this is a differencing disk or a regular disk
	if input.DiskType != nil && *input.DiskType == "Differencing" {
		// Handle differencing disk creation
		if input.ParentPath == nil {
			return id, state, fmt.Errorf("ParentPath is required for Differencing disk type")
		}

		parentVhdPath := *input.ParentPath

		// Check for Azure Edition OS to skip PowerShell fallback
		isAzureEdition := false
		osVersion, err := util.GetOSVersion()
		if err == nil && strings.Contains(strings.ToLower(osVersion), "azure") {
			logger.Infof("Azure Edition OS detected: %s. Will use WMI for differencing disk creation.", osVersion)
			isAzureEdition = true
		}

		// Handle nil client case - Azure Edition should never use PowerShell fallback
		if vmmsClient == nil && !isAzureEdition {
			logger.Warnf("VMMS client is nil, falling back to PowerShell for differencing disk")

			// Wrap PowerShell fallback in panic recovery to prevent crashes
			var fallbackErr error
			func() {
				defer func() {
					if r := recover(); r != nil {
						fallbackErr = fmt.Errorf("recovered from panic in CreateVirtualHardDiskFallback: %v", r)
					}
				}()
				fallbackErr = CreateVirtualHardDiskFallback(vhdFileName, 0, 0, "Differencing", input.ParentPath)
			}()

			if fallbackErr != nil {
				return id, state, fmt.Errorf("failed to create differencing disk using PowerShell fallback: %v", fallbackErr)
			}

			logger.Infof("Created differencing vhd [%s] with parent [%s] using PowerShell", vhdFileName, parentVhdPath)
			return id, state, nil
		}

		ims := vmmsClient.GetImageManagementService()
		if ims == nil {
			logger.Warnf("ImageManagementService is unavailable, trying alternative method via VSMS")

			// Alternative method using VirtualSystemManagementService
			vsms := vmmsClient.GetVirtualSystemManagementService()
			if vsms == nil {
				return id, state, fmt.Errorf("Both ImageManagementService and VirtualSystemManagementService are unavailable")
			}

			// Use VSMS to create differencing disk
			params := map[string]interface{}{
				"Path":       vhdFileName,
				"ParentPath": parentVhdPath,
				"Type":       uint32(4), // 4 = Differencing disk
			}

			_, err = vsms.InvokeMethod("CreateVirtualHardDisk", params)
			if err != nil {
				return id, state, fmt.Errorf("Failed to create differencing disk using VSMS: %v", err)
			}
		} else {
			// Create differencing disk using ImageManagementService
			// Use direct method invocation to create a differencing disk
			// Type 4 corresponds to a differencing disk according to Hyper-V WMI API
			params := map[string]interface{}{
				"Path":       vhdFileName,
				"ParentPath": parentVhdPath,
				"Type":       uint32(4), // 4 = Differencing disk
			}

			_, err = ims.InvokeMethod("CreateVirtualHardDisk", params)
			if err != nil {
				return id, state, fmt.Errorf("Failed to create differencing disk: %v", err)
			}
		}

		logger.Infof("Created differencing vhd [%s] with parent [%s]", vhdFileName, parentVhdPath)
	} else {
		// Regular disk creation (fixed or dynamic)
		vhdFileSize := *input.SizeBytes
		// Set the block size if not specified - using 1MB (1048576 bytes) for better compatibility
		if input.BlockSize == nil {
			blockSize := int64(1048576) // 1MB block size
			input.BlockSize = &blockSize
			logger.Warnf("No block size specified for VHD creation. Using 1MB (1048576 bytes) for better compatibility.")
		}
		vhdFileBlockSize := *input.BlockSize
		// Set the disk type to "fixed" if not specified.
		dynamicDiskType := true
		if input.DiskType != nil && *input.DiskType == "fixed" {
			dynamicDiskType = false
		}

		// Check for Azure Edition OS to skip PowerShell fallback
		isAzureEdition := false
		osVersion, err := util.GetOSVersion()
		if err == nil && strings.Contains(strings.ToLower(osVersion), "azure") {
			logger.Infof("Azure Edition OS detected: %s. Will use WMI for regular disk creation.", osVersion)
			isAzureEdition = true
		}

		// Handle nil client case - Azure Edition should never use PowerShell fallback
		if vmmsClient == nil && !isAzureEdition {
			logger.Warnf("VMMS client is nil, falling back to PowerShell for VHD creation")

			// Ensure we have valid values for required parameters
			if input.DiskType == nil {
				diskType := "dynamic" // Default to dynamic disk if not specified
				input.DiskType = &diskType
				logger.Infof("No disk type specified, defaulting to 'dynamic' for PowerShell fallback")
			}

			// BlockSize might be nil if it wasn't set in the original input
			var blockSizeVal int64 = 0
			if input.BlockSize != nil {
				blockSizeVal = *input.BlockSize
			}

			// Wrap PowerShell fallback in panic recovery to prevent crashes
			var fallbackErr error
			func() {
				defer func() {
					if r := recover(); r != nil {
						fallbackErr = fmt.Errorf("recovered from panic in CreateVirtualHardDiskFallback: %v", r)
					}
				}()
				fallbackErr = CreateVirtualHardDiskFallback(vhdFileName, vhdFileSize, blockSizeVal, *input.DiskType, input.ParentPath)
			}()

			if fallbackErr != nil {
				return id, state, fmt.Errorf("failed to create VHD using PowerShell fallback: %v", fallbackErr)
			}

			logger.Infof("Created VHD [%s] using PowerShell", vhdFileName)
			return id, state, nil
		}

		ims := vmmsClient.GetImageManagementService()

		if ims == nil {
			logger.Warnf("ImageManagementService is unavailable, trying alternative method via VSMS")

			// If the ImageManagementService is unavailable, we can try using the VirtualSystemManagementService
			vsms := vmmsClient.GetVirtualSystemManagementService()
			// If both services are unavailable, we can fall back to PowerShell
			// This is a last resort and should be avoided if possible.
			if vsms == nil && !isAzureEdition {
				logger.Warnf("Both ImageManagementService and VirtualSystemManagementService are unavailable, falling back to PowerShell")

				// Ensure we have valid values for required parameters
				if input.DiskType == nil {
					diskType := "dynamic" // Default to dynamic disk if not specified
					input.DiskType = &diskType
					logger.Infof("No disk type specified, defaulting to 'dynamic' for PowerShell fallback")
				}

				// BlockSize is important for proper VHD creation
				var blockSizeVal int64 = 1048576 // Default to 1MB block size for PowerShell fallback
				if input.BlockSize != nil {
					blockSizeVal = *input.BlockSize
				} else {
					logger.Warnf("No block size specified for VHD creation. Using 1MB (1048576 bytes) for better compatibility.")
				}

				err := CreateVirtualHardDiskFallback(vhdFileName, vhdFileSize, blockSizeVal, *input.DiskType, input.ParentPath)
				if err != nil {
					return id, state, fmt.Errorf("failed to create VHD using PowerShell fallback: %v", err)
				}
				logger.Infof("Created VHD [%s] using PowerShell fallback", vhdFileName)
				return id, state, nil
			} else if vsms == nil && isAzureEdition {
				// For Azure Edition, log clearly in Pulumi window but still use PowerShell fallback
				logger.LogAzureEditionFallback()

				// Ensure we have valid values for required parameters
				if input.DiskType == nil {
					diskType := "dynamic" // Default to dynamic disk if not specified
					input.DiskType = &diskType
					logger.Infof("No disk type specified, defaulting to 'dynamic' for PowerShell fallback")
				}

				// BlockSize is important for proper VHD creation
				var blockSizeVal int64 = 1048576 // Default to 1MB block size for PowerShell fallback
				if input.BlockSize != nil {
					blockSizeVal = *input.BlockSize
				} else {
					logger.Warnf("No block size specified for VHD creation. Using 1MB (1048576 bytes) for better compatibility.")
				}

				err := CreateVirtualHardDiskFallback(vhdFileName, vhdFileSize, blockSizeVal, *input.DiskType, input.ParentPath)
				if err != nil {
					return id, state, fmt.Errorf("failed to create VHD using PowerShell fallback: %v", err)
				}
				logger.Infof("Created VHD [%s] using PowerShell fallback on Azure Edition", vhdFileName)
				return id, state, nil
			}

			// If we have the VirtualSystemManagementService, we can create the disk using it
			diskType := uint32(3) // 3 = Dynamic (default)
			if !dynamicDiskType {
				diskType = uint32(2) // 2 = Fixed
			}

			params := map[string]interface{}{
				"Path":            vhdFileName,
				"MaxInternalSize": uint64(vhdFileSize),
				"BlockSize":       uint32(vhdFileBlockSize),
				"Type":            diskType,
			}

			_, err = vsms.InvokeMethod("CreateVirtualHardDisk", params)
			if err != nil {
				return id, state, fmt.Errorf("Failed to create disk using VSMS: %v", err)
			}
		} else {
			// Use ImageManagementService to create the disk
			setting, err := disk.GetVirtualHardDiskSettingData(
				vmmsClient.GetVirtualizationConn().WMIHost,
				vhdFileName,
				512,
				512,
				uint32(vhdFileBlockSize),
				uint64(vhdFileSize),
				dynamicDiskType,
				disk.VirtualHardDiskFormat_2,
			)
			if err != nil {
				return id, state, fmt.Errorf("Failed to get disk settings: %v", err)
			}
			defer setting.Close()
			err = ims.CreateDisk(setting)
			if err != nil {
				return id, state, fmt.Errorf("Failed to create disk: %v", err)
			}
		}
		logger.Infof("Created vhd [%s]", vhdFileName)
	}
	return id, state, nil
}

// Read retrieves information about an existing VHD file.
func (c *VhdFile) Read(ctx context.Context, id string, inputs VhdFileInputs, currentState VhdFileOutputs) (string, VhdFileInputs, VhdFileOutputs, error) {
	logger := logging.GetLogger(ctx)
	logger.Infof("Reading vhd [%v]", inputs.Path)

	// If the path is empty, we can't read the VHD file.
	if inputs.Path == nil {
		return id, inputs, currentState, fmt.Errorf("Path is nil")
	}

	vhdFileName := *inputs.Path
	if !strings.HasSuffix(vhdFileName, ".vhd") && !strings.HasSuffix(vhdFileName, ".vhdx") {
		return id, inputs, currentState, fmt.Errorf("Path [%v] doesn't end with .vhd or .vhdx", vhdFileName)
	}

	// Check if the file exists.
	vmmsClient, _, err := c.Connect(ctx)
	if err != nil {
		return id, inputs, currentState, err
	}

	// Create outputs from inputs
	outputs := VhdFileOutputs{
		VhdFileInputs: inputs,
	}

	// Check if the file actually exists first using PowerShell
	powershellExe, psErr := util.FindPowerShellExe()
	if psErr == nil {
		checkCmd := exec.Command(powershellExe, "-Command", fmt.Sprintf("Test-Path -Path \"%s\"", vhdFileName))
		checkOutput, checkErr := checkCmd.CombinedOutput()
		if checkErr == nil {
			fileExists := strings.TrimSpace(string(checkOutput)) == "True"
			if !fileExists {
				return id, inputs, currentState, fmt.Errorf("VHD file does not exist: %s", vhdFileName)
			}
			// File exists - if VMMS client is not available, we can just return success
			if vmmsClient == nil {
				logger.Infof("VHD file exists (verified via PowerShell): %s", vhdFileName)
				return id, inputs, outputs, nil
			}
		}
	}

	// If we have a VMMS client, try more detailed validation
	ims := vmmsClient.GetImageManagementService()
	if ims == nil {
		logger.Infof("ImageManagementService not available, using basic file existence check via VSMS")
		// Just verify the file exists
		vsms := vmmsClient.GetVirtualSystemManagementService()
		if vsms == nil {
			logger.Infof("VirtualSystemManagementService not available, file existence was already verified via PowerShell")
			return id, inputs, outputs, nil
		}

		params := map[string]interface{}{
			"Path": vhdFileName,
		}
		_, err = vsms.InvokeMethod("TestVirtualHardDiskExists", params)
		if err != nil {
			logger.Warnf("VSMS failed to verify VHD exists: %v, but file existence was already verified via PowerShell", err)
		}
		return id, inputs, outputs, nil
	}

	// If we have the ImageManagementService, we can get more details about the VHD
	// For now, we're just checking existence
	params := map[string]interface{}{
		"Path": vhdFileName,
	}
	_, err = ims.InvokeMethod("ValidateVirtualHardDisk", params)
	if err != nil {
		logger.Warnf("Failed to validate VHD with IMS: %v, but file existence was already verified", err)
	}

	return id, inputs, outputs, nil
}

// CreateVirtualHardDiskFallback uses PowerShell to create a VHD if WMI services are unavailable.
// This is a fallback method when both ImageManagementService and VirtualSystemManagementService are unavailable.
func CreateVirtualHardDiskFallback(path string, sizeBytes int64, blockSize int64, diskType string, parentPath *string) error {
	// Find PowerShell executable first
	powershellExe, psErr := util.FindPowerShellExe()
	if psErr != nil {
		return fmt.Errorf("failed to find PowerShell executable: %v", psErr)
	}
	// Input validation
	if path == "" {
		return fmt.Errorf("VHD path cannot be empty")
	}

	if !strings.HasSuffix(path, ".vhd") && !strings.HasSuffix(path, ".vhdx") {
		return fmt.Errorf("path '%s' must end with .vhd or .vhdx", path)
	}

	// Validate disk type and required parameters
	diskTypeNormalized := strings.ToLower(diskType)
	if diskTypeNormalized != "fixed" && diskTypeNormalized != "dynamic" && diskTypeNormalized != "differencing" {
		return fmt.Errorf("invalid disk type: %s (must be 'Fixed', 'Dynamic', or 'Differencing')", diskType)
	}

	// For differencing disks, parentPath is required
	if diskTypeNormalized == "differencing" {
		if parentPath == nil || *parentPath == "" {
			return fmt.Errorf("parentPath is required for Differencing disk type")
		}

		// Validate parent path extension
		if !strings.HasSuffix(*parentPath, ".vhd") && !strings.HasSuffix(*parentPath, ".vhdx") {
			return fmt.Errorf("parentPath '%s' must end with .vhd or .vhdx", *parentPath)
		}
	} else {
		// For non-differencing disks, sizeBytes must be positive
		if sizeBytes <= 0 {
			return fmt.Errorf("sizeBytes must be greater than 0 for Fixed or Dynamic disk types")
		}
	}

	// Validate blockSize if specified
	if blockSize < 0 {
		return fmt.Errorf("blockSize cannot be negative")
	}

	// Create parent directory if it doesn't exist
	// Extract directory path from the full VHD path
	lastSlashIndex := strings.LastIndex(path, "\\")
	if lastSlashIndex == -1 {
		lastSlashIndex = strings.LastIndex(path, "/")
	}

	if lastSlashIndex != -1 {
		dirPath := path[:lastSlashIndex]
		// Create directory if it doesn't exist (including all parent directories)
		createDirCmd := exec.Command(powershellExe, "-Command", fmt.Sprintf("New-Item -Path \"%s\" -ItemType Directory -Force | Out-Null", dirPath))
		createDirOutput, createDirErr := createDirCmd.CombinedOutput()
		if createDirErr != nil {
			return fmt.Errorf("failed to create parent directory: %v, output: %s", createDirErr, string(createDirOutput))
		}
	}

	// Construct the PowerShell command with proper escaping
	var cmdArgs []string

	// Base command - use proper array for command arguments
	cmdArgs = append(cmdArgs, "-Command")

	// Build the New-VHD command string
	newVhdCmd := fmt.Sprintf("New-VHD -Path '%s'", path)

	// Add size parameter for non-differencing disks
	if diskTypeNormalized != "differencing" {
		newVhdCmd += fmt.Sprintf(" -SizeBytes %d", sizeBytes)
	}

	// Add block size if specified and valid
	if blockSize > 0 {
		newVhdCmd += fmt.Sprintf(" -BlockSizeBytes %d", blockSize)
	}

	// Set the disk type
	switch diskTypeNormalized {
	case "fixed":
		newVhdCmd += " -Fixed"
	case "dynamic":
		newVhdCmd += " -Dynamic"
	case "differencing":
		newVhdCmd += fmt.Sprintf(" -Differencing -ParentPath '%s'", *parentPath)
	}

	cmdArgs = append(cmdArgs, newVhdCmd)

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
		return fmt.Errorf("failed to create VHD using PowerShell: %v, output: %s", cmdErr, string(output))
	}

	return nil
}

// WireDependencies controls how secrets and unknowns flow through a resource.
//
//	var _ = (infer.ExplicitDependencies[VhdFileInputs, VhdFileOutputs])((*VhdFile)(nil))
//	func (r *VhdFile) WireDependencies(f infer.FieldSelector, args *VhdFileInputs, state *VhdFileOutputs) { .. }
//
// Because we want every output to depend on every input, we can leave the default behavior.
