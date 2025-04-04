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

package machine

import (
	_ "embed"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
)

//go:embed machine.md
var resourceDoc string

// This is the type that implements the Vm resource methods.
// The methods are declared in the vmController.go file.
type Machine struct{}

// The following statement is not required. It is a type assertion to indicate to Go that Vm
// implements the following interfaces. If the function signature doesn't match or isn't implemented,
// we get nice compile time errors at this location.

var _ = (infer.Annotated)((*Machine)(nil))

// Implementing Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (c *Machine) Annotate(a infer.Annotator) {
	a.Describe(&c, resourceDoc)
}

// This is the type that implements the network adapter inputs.
type NetworkAdapterInput struct {
	Name       *string `pulumi:"name"`
	SwitchName *string `pulumi:"switchName"`
}

type HardDriveInput struct {
	Path               *string `pulumi:"path"`
	ControllerType     *string `pulumi:"controllerType"`
	ControllerNumber   *int    `pulumi:"controllerNumber"`
	ControllerLocation *int    `pulumi:"controllerLocation"`
}

// These are the inputs (or arguments) to a Vm resource.
type MachineInputs struct {
	common.ResourceInputs
	MachineName    *string `pulumi:"machineName"`
	Generation     *int    `pulumi:"generation,optional"`
	ProcessorCount *int    `pulumi:"processorCount,optional"`
	MemorySize     *int    `pulumi:"memorySize,optional"`
	// NetworkAdapters []*NetworkAdapterInput `pulumi:"networkAdapters,optional"`
	// HardDrives      []*HardDriveInput      `pulumi:"hardDrives,optional"`
}

func (c *MachineInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.MachineName, "Name of the Virtual Machine")
	a.Describe(&c.ProcessorCount, "Number of processors to allocate to the Virtual Machine. Defaults to 1.")
	a.Describe(&c.MemorySize, "Amount of memory to allocate to the Virtual Machine in MB. Defaults to 1024.")
	a.Describe(&c.Generation, "Generation of the Virtual Machine. Defaults to 2.")
	// a.Describe(&c.NetworkAdapters, "Network adapters to attach to the Virtual Machine.")
	// a.Describe(&c.HardDrives, "Hard drives to attach to the Virtual Machine.")
	// a.Describe(&c.HardDrives[0].Path, "Path to the hard drive file.")
	// a.Describe(&c.HardDrives[0].ControllerType, "Type of controller to use for the hard drive. Defaults to SCSI.")
	// a.Describe(&c.HardDrives[0].ControllerNumber, "Number of the controller to use for the hard drive. Defaults to 0.")
	// a.Describe(&c.HardDrives[0].ControllerLocation, "Location of the controller to use for the hard drive. Defaults to 0.")
	// a.Describe(&c.NetworkAdapters[0].Name, "Name of the network adapter.")
	// a.Describe(&c.NetworkAdapters[0].SwitchName, "Name of the virtual switch to connect the network adapter to.")
}

// These are the outputs (or properties) of a Vm resource.
type MachineOutputs struct {
	MachineInputs
}

func (c *MachineOutputs) Annotate(a infer.Annotator) {}
