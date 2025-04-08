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

package virtualswitch

import (
	_ "embed"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
)

//go:embed virtualswitch.md
var resourceDoc string

// This is the type that implements the VirtualSwitch resource methods.
// The methods are declared in the virtualSwitchController.go file.
type VirtualSwitch struct{}

// The following statement is not required. It is a type assertion to indicate to Go that VirtualSwitch
// implements the following interfaces. If the function signature doesn't match or isn't implemented,
// we get nice compile time errors at this location.

var _ = (infer.Annotated)((*VirtualSwitch)(nil))

// Implementing Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (c *VirtualSwitch) Annotate(a infer.Annotator) {
	a.Describe(&c, resourceDoc)
}

// These are the inputs (or arguments) to a VirtualSwitch resource.
type VirtualSwitchInputs struct {
	common.ResourceInputs
	Name              *string `pulumi:"name"`
	SwitchType        *string `pulumi:"switchType"`
	AllowManagementOs *bool   `pulumi:"allowManagementOs,optional"`
	NetAdapterName    *string `pulumi:"netAdapterName,optional"`
}

func (c *VirtualSwitchInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Name, "Name of the virtual switch")
	a.Describe(&c.SwitchType, "Type of switch: 'External', 'Internal', or 'Private'")
	a.Describe(&c.AllowManagementOs, "Allow the management OS to access the switch (External switches)")
	a.Describe(&c.NetAdapterName, "Name of the physical network adapter to bind to (External switches)")
}

// These are the outputs (or properties) of a VirtualSwitch resource.
type VirtualSwitchOutputs struct {
	VirtualSwitchInputs
}

func (c *VirtualSwitchOutputs) Annotate(a infer.Annotator) {}
