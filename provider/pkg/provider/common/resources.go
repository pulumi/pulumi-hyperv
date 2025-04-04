package common

import (
	"fmt"

	wmi "github.com/microsoft/wmi/pkg/wmiinstance" // Updated import path
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
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
func AddResourceSettings(v *vmms.VMMS, systemSettings *wmi.WmiInstance, resourceSettings []*wmi.WmiInstance) ([]*wmi.WmiInstance, error) {
	return nil, fmt.Errorf("AddResourceSettings not implemented")
	// 	var resultingResourceSettings []*wmi.WmiInstance

	// 	// Convert resource settings to an array of strings
	// 	rsStrings := make([]string, len(resourceSettings))
	// 	for i, rs := range resourceSettings {
	// 		rsStrings[i] = rs.InstancePath()
	// 	}

	// 	params := map[string]interface{}{
	// 		"AffectedConfiguration": systemSettings.InstancePath(),
	// 		"ResourceSettings":      rsStrings,
	// 	}

	// 	result, err := v.VirtualMachineManagementService().InvokeMethod("AddResourceSettings", params)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to add resource settings: %w", err)
	// 	}

	// 	if err := v.ValidateOutput(result); err != nil {
	// 		return nil, err
	// 	}

	// 	resultStrings, err := result.GetStringArray("ResultingResourceSettings")
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to get resulting resource settings: %w", err)
	// 	}

	// 	for _, path := range resultStrings {
	// 		obj, err := v.VirtualizationConn().Get(path)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("failed to get resource setting object: %w", err)
	// 		}
	// 		resultingResourceSettings = append(resultingResourceSettings, obj)
	// 	}

	// 	return resultingResourceSettings, nil
	// }

	// // CreateResource creates a resource of the specified type.
	// func CreateResource(v *vmms.VMMS, resource Resource) (*wmi.WmiInstance, error) {
	// 	resourcePoolClass := "Msvm_ResourcePool"
	// 	if resource == ResourceProcessor {
	// 		resourcePoolClass = "Msvm_ProcessorPool"
	// 	}

	// 	resourceSubType := ResourceSubType(resource)

	// 	// Query for the resource pool
	// 	query := fmt.Sprintf("SELECT * FROM %s WHERE ResourceSubType = '%s' AND Primordial = TRUE", resourcePoolClass, resourceSubType)
	// 	resourcePools, err := v.VirtualizationConn().Query(query)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to query resource pool: %w", err)
	// 	}

	// 	if len(resourcePools) == 0 {
	// 		return nil, fmt.Errorf("no resource pool found for resource subtype %s", resourceSubType)
	// 	}

	// 	// Get allocation capabilities
	// 	capQuery := fmt.Sprintf("ASSOCIATORS OF {%s} WHERE ResultClass=Msvm_AllocationCapabilities", resourcePools[0].Path())
	// 	caps, err := v.VirtualizationConn().Query(capQuery)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to query allocation capabilities: %w", err)
	// 	}

	// 	if len(caps) == 0 {
	// 		return nil, fmt.Errorf("no allocation capabilities found for resource subtype %s", resourceSubType)
	// 	}

	// 	// Get default settings
	// 	settingQuery := fmt.Sprintf("ASSOCIATORS OF {%s} WHERE ResultClass=Msvm_SettingsDefineCapabilities", caps[0].Path())
	// 	settings, err := v.VirtualizationConn().Query(settingQuery)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to query settings: %w", err)
	// 	}

	// 	var defaultSettingPath string
	// 	for _, setting := range settings {
	// 		valueRole, err := setting.GetUint16("ValueRole")
	// 		if err != nil {
	// 			continue
	// 		}
	// 		if valueRole == 0 {
	// 			defaultSettingPath, err = setting.GetString("PartComponent")
	// 			if err != nil {
	// 				continue
	// 			}
	// 			break
	// 		}
	// 	}

	// 	if defaultSettingPath == "" {
	// 		return nil, fmt.Errorf("unable to find the Default Resource Settings")
	// 	}

	// return v.VirtualizationConn().Get(defaultSettingPath)
}
