package common

import (
	"fmt"

	wmi "github.com/microsoft/wmi/pkg/wmiinstance" // Updated import path
	"github.com/pulumi/pulumi-hyperv/provider/pkg/provider/vmms"
)

// Feature types
type Feature int

const (
	FeatureBandwidth Feature = iota
	FeatureOffload
	FeatureSecurity
	FeatureVlan
)

// func featureGUID(feature Feature) string {
// 	var featureGUID string
// 	switch feature {
// 	case FeatureBandwidth:
// 		featureGUID = "24AD3CE1-69BD-4978-B2AC-DAAD389D699C"
// 	case FeatureOffload:
// 		featureGUID = "C885BFD1-ABB7-418F-8163-9F379C9F7166"
// 	case FeatureSecurity:
// 		featureGUID = "776E0BA7-94A1-41C8-8F28-951F524251B5"
// 	case FeatureVlan:
// 		featureGUID = "952C5004-4465-451C-8CB8-FA9AB382B773"
// 	}
// 	return featureGUID
// }

// CreateFeatureSettings creates feature settings for a specified feature.
func CreateFeatureSettings(v *vmms.VMMS, feature Feature) (*wmi.WmiInstance, error) {

	return nil, fmt.Errorf("feature %d not implemented", feature)
	//featureGUID := featureGUID(feature)
	// // Query for feature capabilities matching the feature GUID
	// query := fmt.Sprintf("SELECT * FROM Msvm_EthernetSwitchFeatureCapabilities WHERE FeatureId = '%s'", featureGUID)
	// featureCaps, err := v.VirtualizationConn().QueryInstances(query)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to query feature capabilities: %w", err)
	// }

	// if len(featureCaps) == 0 {
	// 	return nil, fmt.Errorf("no feature capabilities found for feature ID %s", featureGUID)
	// }

	// // Get associated default feature settings
	// assocQuery := fmt.Sprintf("ASSOCIATORS OF {%s} WHERE ResultClass=Msvm_FeatureSettingsDefineCapabilities", featureCaps[0].Path())
	// featureSettingAssocs, err := v.VirtualizationConn().QueryInstances(assocQuery)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to query feature setting associations: %w", err)
	// }

	// var defaultFeatureSettingPath string
	// for _, assoc := range featureSettingAssocs {
	// 	valueRole, err := assoc.GetUint16("ValueRole")
	// 	if err != nil {
	// 		continue
	// 	}
	// 	if valueRole == 0 {
	// 		defaultFeatureSettingPath, err = assoc.GetString("PartComponent")
	// 		if err != nil {
	// 			continue
	// 		}
	// 		break
	// 	}
	// }

	// if defaultFeatureSettingPath == "" {
	// 	return nil, fmt.Errorf("unable to find the Default Feature Settings")
	// }

	// return v.VirtualizationConn().Get(defaultFeatureSettingPath)
}

// ModifyFeatureSettings modifies feature settings.
// func ModifyFeatureSettings(v *vmms.VMMS, featureSettings []*wmi.WmiInstance) ([]*wmi.Result, error) {
// 	var resultingFeatureSettings []*wmi.Result

// 	// Convert feature settings to an array of strings
// 	fsTexts := make([]string, len(featureSettings))
// 	for i, fs := range featureSettings {
// 		text, err := fs.GetText()
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get feature setting text: %w", err)
// 		}
// 		fsTexts[i] = text
// 	}

// 	params := map[string]interface{}{
// 		"FeatureSettings": fsTexts,
// 	}

// 	result, err := v.VirtualMachineManagementService().InvokeMethod("ModifyFeatureSettings", params)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to modify feature settings: %w", err)
// 	}

// 	if err := v.ValidateOutput(result); err != nil {
// 		return nil, err
// 	}

// 	resultStrings, err := result.GetStringArray("ResultingFeatureSettings")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get resulting feature settings: %w", err)
// 	}

// 	for _, path := range resultStrings {
// 		obj, err := v.VirtualizationConn().Get(path)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get feature setting object: %w", err)
// 		}
// 		resultingFeatureSettings = append(resultingFeatureSettings, obj)
// 	}

// 	return resultingFeatureSettings, nil
// }

// // AddFeatureSettings adds feature settings to an ethernet port allocation.
// func AddFeatureSettings(v *vmms.VMMS, ethernetPortAllocationSettings *wmi.Result, featureSettings []*wmi.Result) ([]*wmi.Result, error) {
// 	var resultingFeatureSettings []*wmi.Result

// 	// Convert feature settings to an array of strings
// 	fsStrings := make([]string, len(featureSettings))
// 	for i, fs := range featureSettings {
// 		fsStrings[i] = fs.Path()
// 	}

// 	params := map[string]interface{}{
// 		"AffectedConfiguration": ethernetPortAllocationSettings.Path(),
// 		"FeatureSettings":       fsStrings,
// 	}

// 	result, err := v.VirtualMachineManagementService().InvokeMethod("AddFeatureSettings", params)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to add feature settings: %w", err)
// 	}

// 	if err := v.ValidateOutput(result); err != nil {
// 		return nil, err
// 	}

// 	resultStrings, err := result.GetStringArray("ResultingFeatureSettings")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get resulting feature settings: %w", err)
// 	}

// 	for _, path := range resultStrings {
// 		obj, err := v.VirtualizationConn().Get(path)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get feature setting object: %w", err)
// 		}
// 		resultingFeatureSettings = append(resultingFeatureSettings, obj)
// 	}

// 	return resultingFeatureSettings, nil
// }
