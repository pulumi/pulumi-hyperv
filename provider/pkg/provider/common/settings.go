package common

import (
	"fmt"

	wmi "github.com/microsoft/wmi/pkg/wmiinstance" // Updated import path
	"github.com/pulumi/pulumi-hyperv/provider/pkg/provider/vmms"
)

// Setting types
type Setting int

const (
	SettingSystem Setting = iota
	SettingSecurity
	SettingResource
	SettingMemory
	SettingProcessor
	SettingStorage
	SettingVirtualHardDisk
	SettingNetworkAdapter
	SettingSwitchPort
	SettingSwitchPortOffload
	SettingShutdown
	SettingTimeSynchronization
	SettingDataExchange
	SettingHeartbeat
	SettingVolumeShadowCopy
	SettingGuestServices
)

// SettingsClass returns the WMI class name for a settings type.
func SettingsClass(setting Setting) string {
	switch setting {
	case SettingSystem:
		return "Msvm_VirtualSystemSettingData"
	case SettingSecurity:
		return "Msvm_SecuritySettingData"
	case SettingResource:
		return "Msvm_ResourceAllocationSettingData"
	case SettingMemory:
		return "Msvm_MemorySettingData"
	case SettingProcessor:
		return "Msvm_ProcessorSettingData"
	case SettingStorage:
		return "Msvm_StorageAllocationSettingData"
	case SettingVirtualHardDisk:
		return "Msvm_VirtualHardDiskSettingData"
	case SettingNetworkAdapter:
		return "Msvm_SyntheticEthernetPortSettingData"
	case SettingSwitchPort:
		return "Msvm_EthernetPortAllocationSettingData"
	case SettingSwitchPortOffload:
		return "Msvm_EthernetSwitchPortOffloadSettingData"
	case SettingShutdown:
		return "Msvm_ShutdownComponentSettingData"
	case SettingTimeSynchronization:
		return "Msvm_TimeSyncComponentSettingData"
	case SettingDataExchange:
		return "Msvm_KvpExchangeComponentSettingData"
	case SettingHeartbeat:
		return "Msvm_HeartbeatComponentSettingData"
	case SettingVolumeShadowCopy:
		return "Msvm_VssComponentSettingData"
	case SettingGuestServices:
		return "Msvm_GuestServiceInterfaceComponentSettingData"
	}

	return ""
}

// CreateSettings creates settings of the specified type.
func CreateSettings(v *vmms.VMMS, setting Setting) (*wmi.WmiClass, error) {
	className := SettingsClass(setting)
	if className == "" {
		return nil, fmt.Errorf("invalid setting type: %d", setting)
	}

	return v.GetVirtualizationConn().GetClass(className)
}

// GetRelatedSettings gets settings of the specified type related to an instance.
func GetRelatedSettings(v *vmms.VMMS, instance *wmi.WmiInstance, setting Setting) (*wmi.WmiInstance, error) {
	className := SettingsClass(setting)
	if className == "" {
		return nil, fmt.Errorf("invalid setting type: %d", setting)
	}

	assocQuery := fmt.Sprintf("ASSOCIATORS OF {%s} WHERE ResultClass=%s", instance.InstancePath(), className)
	settings, err := v.GetVirtualizationConn().QueryInstances(assocQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query related settings: %w", err)
	}

	if len(settings) == 0 {
		return nil, fmt.Errorf("no related settings found of type %s", className)
	}

	return settings[0], nil
}
