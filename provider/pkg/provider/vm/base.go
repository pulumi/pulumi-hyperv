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

	// Updated import path
	"github.com/pulumi/pulumi-go-provider/infer"
)

// BaseInputs is the common set of inputs for all local commands.
type BaseInputs struct {
	VmName *string `pulumi:"vmname"`
}

// Implementing Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *BaseInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.VmName, "Name of the Virtual Machine")
}

type BaseOutputs struct{}

// Implementing Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *BaseOutputs) Annotate(a infer.Annotator) {}
