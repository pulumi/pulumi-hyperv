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
	"os/exec"
	"strings"

	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/storage/disk"
	"github.com/pulumi/pulumi-go-provider/infer"

	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/logging"
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
	// Initialize all the parameters.
	config := infer.GetConfig[common.Config](ctx)
	var whost *host.WmiHost
	if config.Host != "" {
		whost = host.NewWmiHost(config.Host)
	} else {
		whost = host.NewWmiLocalHost()
	}

	vmmsClient, err := vmms.NewVMMS(whost)
	return vmmsClient, nil, err
}

// Delete deletes a VHD file.
func (c *VhdFile) Delete(ctx context.Context, id string, state VhdFileOutputs) error {
	// Delete the VHD file.
	logger := logging.GetLogger(ctx)
	logger.Infof("Deleting vhd [%v]", state.Path)

	// In case any mounts or attachments to this file are missing.
	// Warn about them instead of failing.
	//
	// logger.Warnf("Failed to delete vhd [%v] because [%v]", state.Path, err)
	// return fmt.Errorf("Failed to delete vhd [%v] because [%v]", state.Path, err)
	//
	// return

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

	// First try using the VirtualSystemManagementService
	vsms := vmmsClient.GetVirtualSystemManagementService()
	if vsms == nil {
		return fmt.Errorf("cannot delete VHD - VirtualSystemManagementService is unavailable")
	}

	params := map[string]interface{}{
		"Path": vhdPath,
	}
	_, err = vsms.InvokeMethod("DeleteVirtualHardDisk", params)
	if err != nil {
		logger.Warnf("Failed to delete vhd [%v] because [%v]", vhdPath, err)
		return fmt.Errorf("Failed to delete vhd [%v] because [%v]", vhdPath, err)
	}

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
	// Create the VHD file.
	vhdFileName := *input.Path

	// Check if this is a differencing disk or a regular disk
	if input.DiskType != nil && *input.DiskType == "Differencing" {
		// Handle differencing disk creation
		if input.ParentPath == nil {
			return id, state, fmt.Errorf("ParentPath is required for Differencing disk type")
		}

		parentVhdPath := *input.ParentPath
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
		// Set the block size to 512 bytes if not specified.
		if input.BlockSize == nil {
			blockSize := int64(512)
			input.BlockSize = &blockSize
		}
		vhdFileBlockSize := *input.BlockSize
		// Set the disk type to "fixed" if not specified.
		dynamicDiskType := true
		if input.DiskType != nil && *input.DiskType == "fixed" {
			dynamicDiskType = false
		}
		ims := vmmsClient.GetImageManagementService()

		if ims == nil {
			logger.Warnf("ImageManagementService is unavailable, trying alternative method via VSMS")

			// If the ImageManagementService is unavailable, we can try using the VirtualSystemManagementService
			vsms := vmmsClient.GetVirtualSystemManagementService()
			// If both services are unavailable, we can fall back to PowerShell
			// This is a last resort and should be avoided if possible.
			if vsms == nil {
				logger.Warnf("Both ImageManagementService and VirtualSystemManagementService are unavailable, falling back to PowerShell")

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

				err := CreateVirtualHardDiskFallback(vhdFileName, vhdFileSize, blockSizeVal, *input.DiskType, input.ParentPath)
				if err != nil {
					return id, state, fmt.Errorf("failed to create VHD using PowerShell fallback: %v", err)
				}
				logger.Infof("Created VHD [%s] using PowerShell fallback", vhdFileName)
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

	ims := vmmsClient.GetImageManagementService()
	if ims == nil {
		logger.Infof("ImageManagementService not available, using basic file existence check")
		// Just verify the file exists
		params := map[string]interface{}{
			"Path": vhdFileName,
		}
		vsms := vmmsClient.GetVirtualSystemManagementService()
		_, err = vsms.InvokeMethod("TestVirtualHardDiskExists", params)
		if err != nil {
			return id, inputs, currentState, fmt.Errorf("Failed to check if VHD exists: %v", err)
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
		return id, inputs, currentState, fmt.Errorf("Failed to validate VHD: %v", err)
	}

	return id, inputs, outputs, nil
}

// CreateVirtualHardDiskFallback uses PowerShell to create a VHD if WMI services are unavailable.
// This is a fallback method when both ImageManagementService and VirtualSystemManagementService are unavailable.
func CreateVirtualHardDiskFallback(path string, sizeBytes int64, blockSize int64, diskType string, parentPath *string) error {
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

	// Construct the PowerShell command with proper escaping
	var cmdArgs []string

	// Base command
	cmdArgs = append(cmdArgs, "New-VHD", "-Path", fmt.Sprintf("\"%s\"", path))

	// Add size parameter for non-differencing disks
	if diskTypeNormalized != "differencing" {
		cmdArgs = append(cmdArgs, "-SizeBytes", fmt.Sprintf("%d", sizeBytes))
	}

	// Add block size if specified and valid
	if blockSize > 0 {
		cmdArgs = append(cmdArgs, "-BlockSizeBytes", fmt.Sprintf("%d", blockSize))
	}

	// Set the disk type
	switch diskTypeNormalized {
	case "fixed":
		cmdArgs = append(cmdArgs, "-Fixed")
	case "dynamic":
		cmdArgs = append(cmdArgs, "-Dynamic")
	case "differencing":
		cmdArgs = append(cmdArgs, "-Differencing", "-ParentPath", fmt.Sprintf("\"%s\"", *parentPath))
	}

	// Execute the PowerShell command
	cmd := exec.Command("powershell.exe", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create VHD using PowerShell: %v, output: %s", err, string(output))
	}

	return nil
}

// WireDependencies controls how secrets and unknowns flow through a resource.
//
//	var _ = (infer.ExplicitDependencies[VhdFileInputs, VhdFileOutputs])((*VhdFile)(nil))
//	func (r *VhdFile) WireDependencies(f infer.FieldSelector, args *VhdFileInputs, state *VhdFileOutputs) { .. }
//
// Because we want every output to depend on every input, we can leave the default behavior.
