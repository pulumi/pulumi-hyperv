package virtualswitch

import (
	"fmt"

	wmi "github.com/microsoft/wmi/pkg/wmiinstance" // Updated import path
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vmms"
)

// ExistsVirtualSwitch checks if a virtual switch with the given name exists.
func ExistsVirtualSwitch(v *vmms.VMMS, name string) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE Caption = 'Virtual Switch' AND ElementName = '%s'", name)
	switches, err := v.GetVirtualizationConn().QueryInstances(query)
	if err != nil {
		return false, fmt.Errorf("failed to query virtual switches: %w", err)
	}

	return len(switches) > 0, nil
}

// GetVirtualSwitch gets a virtual switch by name.
func GetVirtualSwitch(v *vmms.VMMS, name string) (*wmi.WmiInstance, error) {
	query := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE Caption = 'Virtual Switch' AND ElementName = '%s'", name)
	switches, err := v.GetVirtualizationConn().QueryInstances(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query virtual switches: %w", err)
	}

	if len(switches) == 0 {
		return nil, fmt.Errorf("unable to find the Virtual Switch %s", name)
	}

	return switches[0], nil
}
