package vhd

import (
	"fmt"

	wmi "github.com/microsoft/wmi/pkg/wmiinstance" // Updated import path
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
)

// CreateVirtualHardDisk creates a virtual hard disk.
func CreateVirtualHardDisk(v *vmms.VMMS, vhdSettings *wmi.WmiInstance) error {
	return fmt.Errorf("CreateVirtualHardDisk not implemented")
	// Get the WMI text representation of the VHD settings
	// 	vhdText, err := vhdSettings.GetText()
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get VHD settings text: %w", err)
	// 	}

	// 	params := map[string]interface{}{
	// 		"VirtualDiskSettingData": vhdText,
	// 	}

	// 	result, err := v.GetImageManagementService().InvokeMethod("CreateVirtualHardDisk", params)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create virtual hard disk: %w", err)
	// 	}

	// return v.ValidateOutput(result)
}
