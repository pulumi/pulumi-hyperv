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

package networkadapter

import (
	_ "embed"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
)

//go:embed networkadapter.md
var resourceDoc string

// This is the type that implements the NetworkAdapter resource methods.
// The methods are declared in the networkadapterController.go file.
type NetworkAdapter struct{}

// The following statement is not required. It is a type assertion to indicate to Go that NetworkAdapter
// implements the following interfaces. If the function signature doesn't match or isn't implemented,
// we get nice compile time errors at this location.

var _ = (infer.Annotated)((*NetworkAdapter)(nil))

// Implementing Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (c *NetworkAdapter) Annotate(a infer.Annotator) {
	a.Describe(&c, resourceDoc)
}

// These are the inputs (or arguments) to a NetworkAdapter resource.
type NetworkAdapterInputs struct {
	common.ResourceInputs
	Name            *string `pulumi:"name"`
	VMName          *string `pulumi:"vmName"`
	SwitchName      *string `pulumi:"switchName"`
	MacAddress      *string `pulumi:"macAddress,optional"`
	VlanId          *int    `pulumi:"vlanId,optional"`
	DHCPGuard       *bool   `pulumi:"dhcpGuard,optional"`
	RouterGuard     *bool   `pulumi:"routerGuard,optional"`
	PortMirroring   *string `pulumi:"portMirroring,optional"`
	IeeePriorityTag *bool   `pulumi:"ieeePriorityTag,optional"`
	VMQWeight       *int    `pulumi:"vmqWeight,optional"`
	IPAddresses     *string `pulumi:"ipAddresses,optional"`
}

func (c *NetworkAdapterInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Name, "Name of the network adapter")
	a.Describe(&c.VMName, "Name of the virtual machine to attach the network adapter to")
	a.Describe(&c.SwitchName, "Name of the virtual switch to connect the network adapter to")
	a.Describe(&c.MacAddress, "MAC address for the network adapter. If not specified, a dynamic MAC address will be generated.")
	a.Describe(&c.VlanId, "VLAN ID for the network adapter. If not specified, no VLAN tagging is used.")
	a.Describe(&c.DHCPGuard, "Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages.")
	a.Describe(&c.RouterGuard, "Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages.")
	a.Describe(&c.PortMirroring, "Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None.")
	a.Describe(&c.IeeePriorityTag, "Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value.")
	a.Describe(&c.VMQWeight, "VMQ weight for the network adapter. A value of 0 disables VMQ.")
	a.Describe(&c.IPAddresses, "Comma-separated list of IP addresses to assign to the network adapter.")
}

// These are the outputs (or properties) of a NetworkAdapter resource.
type NetworkAdapterOutputs struct {
	NetworkAdapterInputs
	AdapterId *string `pulumi:"adapterId"`
}

func (c *NetworkAdapterOutputs) Annotate(a infer.Annotator) {
	a.Describe(&c.AdapterId, "The ID of the network adapter")
}