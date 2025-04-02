package common

import (
	"fmt"

	"github.com/microsoft/wmi"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider"
)

// Resource types
type Resource int

const (
	ResourceProcessor Resource = iota
	ResourceMemory
	ResourceSCSIController
	ResourceVirtualHardDrive
	ResourceVirtualHardDisk
	ResourceVirtualDvdDrive
	ResourceVirtualDvdDisk
	ResourceNetworkAdapter
	ResourceSwitchPort
)

func ResourceSubType(r Resource) string {
	var resourceSubType string
	switch r {
	case ResourceProcessor:
		resourceSubType = "Microsoft:Hyper-V:Processor"
	case ResourceMemory:
		resourceSubType = "Microsoft:Hyper-V:Memory"
	case ResourceSCSIController:
		resourceSubType = "Microsoft:Hyper-V:Synthetic SCSI Controller"
	case ResourceVirtualHardDrive:
		resourceSubType = "Microsoft:Hyper-V:Synthetic Disk Drive"
	case ResourceVirtualHardDisk:
		resourceSubType = "Microsoft:Hyper-V:Virtual Hard Disk"
	case ResourceVirtualDvdDrive:
		resourceSubType = "Microsoft:Hyper-V:Synthetic DVD Drive"
	case ResourceVirtualDvdDisk:
		resourceSubType = "Microsoft:Hyper-V:Virtual CD/DVD Disk"
	case ResourceNetworkAdapter:
		resourceSubType = "Microsoft:Hyper-V:Synthetic Ethernet Port"
	case ResourceSwitchPort:
		resourceSubType = "Microsoft:Hyper-V:Ethernet Connection"
	}
	return resourceSubType
}

// AddResourceSettings adds resource settings to a system.
func AddResourceSettings(v *provider.VMMS, systemSettings *wmi.Result, resourceSettings []*wmi.Result) ([]*wmi.Result, error) {
	var resultingResourceSettings []*wmi.Result

	// Convert resource settings to an array of strings
	rsStrings := make([]string, len(resourceSettings))
	for i, rs := range resourceSettings {
		rsStrings[i] = rs.Path()
	}

	params := map[string]interface{}{
		"AffectedConfiguration": systemSettings.Path(),
		"ResourceSettings":      rsStrings,
	}

	result, err := v.VirtualMachineManagementService().InvokeMethod("AddResourceSettings", params)
	if err != nil {
		return nil, fmt.Errorf("failed to add resource settings: %w", err)
	}

	if err := v.ValidateOutput(result); err != nil {
		return nil, err
	}

	resultStrings, err := result.GetStringArray("ResultingResourceSettings")
	if err != nil {
		return nil, fmt.Errorf("failed to get resulting resource settings: %w", err)
	}

	for _, path := range resultStrings {
		obj, err := v.VirtualizationConn().Get(path)
		if err != nil {
			return nil, fmt.Errorf("failed to get resource setting object: %w", err)
		}
		resultingResourceSettings = append(resultingResourceSettings, obj)
	}

	return resultingResourceSettings, nil
}
