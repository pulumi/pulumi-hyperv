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

package vhdfile

import (
	"context"
	"fmt"
	"os"

	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	"github.com/microsoft/wmi/pkg/virtualization/core/storage/disk"

	// Updated import path
	provider "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
)

// The following statements are not required. They are type assertions to indicate to Go that VhdFile implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors at this location.

// They would normally be included in the vhdfileController.go file, but they're located here for instructive purposes.
var _ = (infer.CustomResource[VhdFileInputs, VhdFileOutputs])((*VhdFile)(nil))
var _ = (infer.CustomUpdate[VhdFileInputs, VhdFileOutputs])((*VhdFile)(nil))
var _ = (infer.CustomDelete[VhdFileOutputs])((*VhdFile)(nil))

func (c *VhdFile) Connect(ctx context.Context) (*vmms.VMMS, *service.VirtualSystemManagementService, error) {
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
func (c *VhdFile) Read(ctx context.Context, id string, inputs VhdFileInputs, preview bool) (VhdFileOutputs, error) {
	logger := provider.GetLogger(ctx)

	// If we don't have a path, we can't read the file
	if inputs.Path == nil {
		return VhdFileOutputs{}, fmt.Errorf("path is required to read VHD file")
	}

	// Get the path to the VHD file
	vhdFilePath := *inputs.Path

	// Check if the file exists
	_, err := os.Stat(vhdFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist
			logger.Debug(fmt.Sprintf("VHD file %s does not exist", vhdFilePath))
			return VhdFileOutputs{}, nil
		}
		// Some other error occurred
		return VhdFileOutputs{}, fmt.Errorf("failed to stat VHD file %s: %v", vhdFilePath, err)
	}

	// If the file exists, we try to read the VHD info
	vmmsClient, _, err := c.Connect(ctx)
	if err != nil {
		return VhdFileOutputs{}, err
	}

	// Get VHD settings to verify it exists in Hyper-V
	// For existing VHD files, we pass path but use zeros for the other parameters
	// since we're not creating a new VHD file
	setting, err := disk.GetVirtualHardDiskSettingData(
		vmmsClient.GetVirtualizationConn().WMIHost,
		vhdFilePath,
		0,     // Default sector size
		0,     // Default logical sector size
		0,     // Block size (not needed for reads)
		0,     // Size bytes (not needed for reads)
		false, // Fixed/Dynamic doesn't matter for reads
		disk.VirtualHardDiskFormat_2,
	)

	if err != nil {
		logger.Debug(fmt.Sprintf("Failed to get VHD settings for %s: %v", vhdFilePath, err))
		// Even if we can't get the settings, we still want to return the outputs
		// since the file exists on disk
	} else {
		defer setting.Close()
		logger.Debug(fmt.Sprintf("Successfully retrieved VHD settings for %s", vhdFilePath))

		// Extract properties from the VHD
		// Get block size if available
		blockSize, err := setting.GetPropertyBlockSize()
		if err == nil && blockSize > 0 && inputs.BlockSize == nil {
			blockSizeInt64 := int64(blockSize)
			inputs.BlockSize = &blockSizeInt64
			logger.Debug(fmt.Sprintf("Retrieved block size: %d", blockSize))
		}

		// Get disk type if available
		isDynamic, err := setting.GetPropertyType()
		if err == nil && inputs.DiskType == nil {
			var diskType string
			if isDynamic == disk.VirtualHardDiskType_SPARSE {
				diskType = "dynamic"
			} else {
				diskType = "fixed"
			}
			inputs.DiskType = &diskType
			logger.Debug(fmt.Sprintf("Retrieved disk type: %s", diskType))
		}
	}

	// Create outputs based on inputs
	outputs := VhdFileOutputs{
		VhdFileInputs: inputs,
	}

	logger.Debug(fmt.Sprintf("Successfully read VHD file %s", vhdFilePath))
	return outputs, nil
}

// This is the Create method. This will be run on every VhdFile resource creation.
func (c *VhdFile) Create(ctx context.Context, name string, input VhdFileInputs, preview bool) (string, VhdFileOutputs, error) {
	logger := provider.GetLogger(ctx)
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
		return id, state, fmt.Errorf("Failed [%+v]", err)
	}
	defer setting.Close()
	err = ims.CreateDisk(setting)
	if err != nil {
		return id, state, fmt.Errorf("Failed [%+v]", err)
	}
	logger.Infof("Created vhd [%s]", vhdFileName)
	return id, state, nil
}

// WireDependencies controls how secrets and unknowns flow through a resource.
//
//	var _ = (infer.ExplicitDependencies[VhdFileInputs, VhdFileOutputs])((*VhdFile)(nil))
//	func (r *VhdFile) WireDependencies(f infer.FieldSelector, args *VhdFileInputs, state *VhdFileOutputs) { .. }
//
// Because we want every output to depend on every input, we can leave the default behavior.

// The Update method will be run on every update.
func (c *VhdFile) Update(ctx context.Context, id string, olds VhdFileOutputs, news VhdFileInputs, preview bool) (VhdFileOutputs, error) {
	// This is a no-op. We don't need to do anything here.
	state := VhdFileOutputs{VhdFileInputs: news}
	// If in preview, don't run the command.
	if preview {
		return state, nil
	}
	return state, nil
}

// The Delete method will run when the resource is deleted.
func (c *VhdFile) Delete(ctx context.Context, id string, props VhdFileOutputs) error {
	logger := provider.GetLogger(ctx)
	vhdFileName := *props.Path
	err := os.RemoveAll(vhdFileName)
	if err != nil {
		return fmt.Errorf("Failed to delete VHD file [%s]: %v", vhdFileName, err)
	}
	// If the file was deleted successfully, return nil.
	logger.Infof("Deleted VHD file [%s]\n", vhdFileName)
	return nil
}
