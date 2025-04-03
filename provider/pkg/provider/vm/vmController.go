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

package vm

import (
	"context"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// The following statements are not required. They are type assertions to indicate to Go that Vm implements the following interfaces.
// If the function signature doesn't match or isn't implemented, we get nice compile time errors at this location.

// They would normally be included in the vmController.go file, but they're located here for instructive purposes.
var _ = (infer.CustomResource[VmInputs, VmOutputs])((*Vm)(nil))
var _ = (infer.CustomUpdate[VmInputs, VmOutputs])((*Vm)(nil))
var _ = (infer.CustomDelete[VmOutputs])((*Vm)(nil))
var _ = (infer.CustomRead[VmInputs, VmOutputs])((*Vm)(nil))

// This is the Get Metadata method.
func (c *Vm) Read(ctx context.Context, id string, inputs VmInputs, preview bool) (VmOutputs, error) {
	// This is a no-op. We don't need to do anything here.
	return VmOutputs{}, nil
}

// This is the Create method. This will be run on every Vm resource creation.
func (c *Vm) Create(ctx context.Context, name string, input VmInputs, preview bool) (string, VmOutputs, error) {
	state := VmOutputs{VmInputs: input}
	id, err := resource.NewUniqueHex(name, 8, 0)
	if err != nil {
		return id, state, err
	}

	// If in preview, don't run the command.
	if preview {
		return id, state, nil
	}
	if input.Create == nil {
		return id, state, nil
	}
	//cmd := *input.Create
	return id, state, err
}

// WireDependencies controls how secrets and unknowns flow through a resource.
//
//	var _ = (infer.ExplicitDependencies[VmInputs, VmOutputs])((*Vm)(nil))
//	func (r *Vm) WireDependencies(f infer.FieldSelector, args *VmInputs, state *VmOutputs) { .. }
//
// Because we want every output to depend on every input, we can leave the default behavior.

// The Update method will be run on every update.
func (c *Vm) Update(ctx context.Context, id string, olds VmOutputs, news VmInputs, preview bool) (VmOutputs, error) {
	state := VmOutputs{VmInputs: news, BaseOutputs: olds.BaseOutputs}
	// If in preview, don't run the command.
	if preview {
		return state, nil
	}
	// Use Create command if Update is unspecified.
	cmd := news.Create
	if news.Update != nil {
		cmd = news.Update
	}
	// If neither are specified, do nothing.
	if cmd == nil {
		return state, nil
	}
	return state, nil
}

// The Delete method will run when the resource is deleted.
func (c *Vm) Delete(ctx context.Context, id string, props VmOutputs) error {
	if props.Delete == nil {
		return nil
	}
}
