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

package vhdfile

import (
	_ "embed"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
)

//go:embed vhdfile.md
var resourceDoc string

// This is the type that implements the Vm resource methods.
// The methods are declared in the vmController.go file.
type VhdFile struct{}

// The following statement is not required. It is a type assertion to indicate to Go that Vm
// implements the following interfaces. If the function signature doesn't match or isn't implemented,
// we get nice compile time errors at this location.

var _ = (infer.Annotated)((*VhdFile)(nil))

// Implementing Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (c *VhdFile) Annotate(a infer.Annotator) {
	a.Describe(&c, resourceDoc)
}

// These are the inputs (or arguments) to a Vm resource.
type VhdFileInputs struct {
	common.ResourceInputs
	Path       *string `pulumi:"path"`
	SizeBytes  *int64  `pulumi:"sizeBytes,optional"`
	BlockSize  *int64  `pulumi:"blockSize,optional"`
	ParentPath *string `pulumi:"parentPath,optional"`
	// DiskType is a string that can be either "fixed", "dynamic", or "differencing".
	// It is used to specify the type of VHD file to create.
	// "fixed" means that the VHD file will be created with a fixed size.
	// "dynamic" means that the VHD file will be created with a dynamic size.
	// "differencing" means that the VHD file will be created as a differencing disk.
	// The default value is "fixed".
	DiskType *string `pulumi:"diskType,optional"`
}

func (c *VhdFileInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Path, "Path to the VHD file")
	a.Describe(&c.SizeBytes, "Size of the VHD file in bytes")
	a.Describe(&c.BlockSize, "Block size of the VHD file in bytes")
	a.Describe(&c.ParentPath, "Path to the parent VHD file when creating a differencing disk")
	a.Describe(&c.DiskType, "Type of the VHD file (Fixed, Dynamic, or Differencing)")
}

// These are the outputs (or properties) of a Vm resource.
type VhdFileOutputs struct {
	VhdFileInputs
}

func (c *VhdFileOutputs) Annotate(a infer.Annotator) {}
