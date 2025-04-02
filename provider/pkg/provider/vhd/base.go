package vhd

import (
	"fmt"

	"github.com/microsoft/wmi"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider"
)

// CreateVirtualHardDisk creates a virtual hard disk.
func CreateVirtualHardDisk(v *provider.VMMS, vhdSettings *wmi.Result) error {
	// Get the WMI text representation of the VHD settings
	vhdText, err := vhdSettings.GetText()
	if err != nil {
		return fmt.Errorf("failed to get VHD settings text: %w", err)
	}

	params := map[string]interface{}{
		"VirtualDiskSettingData": vhdText,
	}

	result, err := v.GetImageManagementService().InvokeMethod("CreateVirtualHardDisk", params)
	if err != nil {
		return fmt.Errorf("failed to create virtual hard disk: %w", err)
	}

	return v.ValidateOutput(result)
}
