// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"
	"strings"
	"time"

	"github.com/microsoft/wmi"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider/vm"
)

const (
	Name = "hyperv"
)

// This provider uses the `pulumi-go-provider` library to produce a code-first provider definition.
func NewProvider() p.Provider {
	return infer.Provider(infer.Options{
		// This is the metadata for the provider
		Metadata: schema.Metadata{
			DisplayName: "Hyperv",
			Description: "The Pulumi hyperv Provider enables you to use Hyper-V resources in your Pulumi programs.",
			Keywords: []string{
				"pulumi",
				"hyperv",
				"category/utility",
				"kind/native",
			},
			Homepage:   "https://pulumi.com",
			License:    "Apache-2.0",
			Repository: "https://github.com/pulumi/pulumi-hyperv-provider",
			Publisher:  "Pulumi",
			LogoURL:    "https://raw.githubusercontent.com/pulumi/pulumi-hyperv-provider/master/assets/logo.svg",
			// This contains language specific details for generating the provider's SDKs
			LanguageMap: map[string]any{
				"csharp": map[string]any{
					"respectSchemaVersion": true,
					"packageReferences": map[string]string{
						"Pulumi": "3.*",
					},
				},
				"go": map[string]any{
					"respectSchemaVersion":           true,
					"generateResourceContainerTypes": true,
					"importBasePath":                 "github.com/pulumi/pulumi-hyperv-provider/provider/go/hyperv",
				},
				"nodejs": map[string]any{
					"respectSchemaVersion": true,
				},
				"python": map[string]any{
					"respectSchemaVersion": true,
					"pyproject": map[string]bool{
						"enabled": true,
					},
				},
				"java": map[string]any{
					"buildFiles":                      "gradle",
					"gradleNexusPublishPluginVersion": "2.0.0",
					"dependencies": map[string]any{
						"com.pulumi:pulumi":               "1.0.0",
						"com.google.code.gson:gson":       "2.8.9",
						"com.google.code.findbugs:jsr305": "3.0.2",
					},
				},
			},
		},
		// A list of `infer.Resource` that are provided by the provider.
		Resources: []infer.InferredResource{
			// The hyperv resource implementation is commented extensively for new pulumi-go-provider developers.
			infer.Resource[
				// 1. This type is an interface that implements the logic for the Resource
				//    these methods include `Create`, `Update`, `Delete`, and `WireDependencies`.
				//    `WireDependencies` should be implemented to preserve the secretness of an input
				*vm.Vm,
				// 2. The type of the Inputs/Arguments to supply to the Resource.
				vm.VmInputs,
				// 3. The type of the Output/Properties/Fields of a created Resource.
				vm.VmOutputs,
			](),
		},
		// Functions or invokes that are provided by the provider.
		Functions: []infer.InferredFunction{
			// The Run function is commented extensively for new pulumi-go-provider developers.
			infer.Function[*vm.Run, vm.RunInputs, vm.RunOutputs](),
		},
	})
}

// VMMS represents the Hyper-V Virtual Machine Management Service.
type VMMS struct {
	host                string
	virtualizationConn  *wmi.Connection
	hgsConn             *wmi.Connection
	securityService     *wmi.Result
	imageManagementSvc  *wmi.Result
	vmManagementService *wmi.Result
}

// NewVMMS creates a new VMMS instance.
func NewVMMS(host string) (*VMMS, error) {
	vmms := &VMMS{
		host: host,
	}

	// Set up virtualization connection
	virtConn, err := wmi.NewConnection("root\\virtualization\\v2")
	if err != nil {
		return nil, fmt.Errorf("failed to create virtualization connection: %w", err)
	}
	vmms.virtualizationConn = virtConn

	// Set up HGS connection
	hgsConn, err := wmi.NewConnection("root\\Microsoft\\Windows\\Hgs")
	if err != nil {
		return nil, fmt.Errorf("failed to create HGS connection: %w", err)
	}
	vmms.hgsConn = hgsConn

	// Get services
	ss, err := vmms.GetSecurityService()
	if err != nil {
		return nil, err
	}
	vmms.securityService = ss

	ims, err := vmms.GetImageManagementService()
	if err != nil {
		return nil, err
	}
	vmms.imageManagementSvc = ims

	vmmSvc, err := vmms.GetVirtualMachineManagementService()
	if err != nil {
		return nil, err
	}
	vmms.vmManagementService = vmmSvc

	return vmms, nil
}

// GetVirtualizationConn returns the virtualization connection.
func (v *VMMS) VirtualizationConn() *wmi.Connection {
	return v.virtualizationConn
}

// GetHgsConn returns the HGS connection.
func (v *VMMS) HgsConn() *wmi.Connection {
	return v.hgsConn
}

// GetSecurityService returns the security service.
func (v *VMMS) SecurityService() *wmi.Result {
	return v.securityService
}

// GetImageManagementService returns the image management service.
func (v *VMMS) ImageManagementService() *wmi.Result {
	return v.imageManagementSvc
}

// GetVirtualMachineManagementService returns the virtual machine management service.
func (v *VMMS) VirtualMachineManagementService() *wmi.Result {
	return v.vmManagementService
}

// Close closes the VMMS connections.
func (v *VMMS) Close() error {
	var errs []string

	if v.securityService != nil {
		v.securityService = nil
	}

	if v.imageManagementSvc != nil {
		v.imageManagementSvc = nil
	}

	if v.vmManagementService != nil {
		v.vmManagementService = nil
	}

	if v.virtualizationConn != nil {
		if err := v.virtualizationConn.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("failed to close virtualization connection: %v", err))
		}
	}

	if v.hgsConn != nil {
		if err := v.hgsConn.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("failed to close HGS connection: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}
	return nil
}

// GetSecurityService returns the Hyper-V security service.
func (v *VMMS) GetSecurityService() (*wmi.Result, error) {
	objs, err := v.virtualizationConn.GetAll("Msvm_SecurityService")
	if err != nil {
		return nil, fmt.Errorf("failed to query security service: %w", err)
	}

	if len(objs) == 0 {
		return nil, fmt.Errorf("no security service found")
	}

	return objs[0], nil
}

// GetImageManagementService returns the Hyper-V image management service.
func (v *VMMS) GetImageManagementService() (*wmi.Result, error) {
	objs, err := v.virtualizationConn.GetAll("Msvm_ImageManagementService")
	if err != nil {
		return nil, fmt.Errorf("failed to query image management service: %w", err)
	}

	if len(objs) == 0 {
		return nil, fmt.Errorf("no image management service found")
	}

	return objs[0], nil
}

// GetVirtualMachineManagementService returns the Hyper-V virtual machine management service.
func (v *VMMS) GetVirtualMachineManagementService() (*wmi.Result, error) {
	objs, err := v.virtualizationConn.GetAll("Msvm_VirtualSystemManagementService")
	if err != nil {
		return nil, fmt.Errorf("failed to query virtual machine management service: %w", err)
	}

	if len(objs) == 0 {
		return nil, fmt.Errorf("no virtual machine management service found")
	}

	return objs[0], nil
}

// GetUntrustedGuardian gets the untrusted guardian.
func (v *VMMS) GetUntrustedGuardian() (*wmi.Result, error) {
	query := "SELECT * FROM MSFT_HgsGuardian WHERE Name = 'UntrustedGuardian'"
	guardians, err := v.hgsConn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query guardians: %w", err)
	}

	if len(guardians) == 0 {
		return nil, nil
	}

	return guardians[0], nil
}

// ErrorCodeMeaning returns a string description for a WMI error code.
func ErrorCodeMeaning(returnValue uint32) string {
	switch returnValue {
	case 0:
		return "Completed with No Error."
	case 1:
		return "Not Supported."
	case 2:
		return "Failed."
	case 3:
		return "Timeout."
	case 4:
		return "Invalid Parameter."
	case 5:
		return "Invalid State."
	case 6:
		return "Invalid Type."
	case 4096:
		return "Method Parameters Checked - Job Started."
	case 32768:
		return "Failed."
	case 32769:
		return "Access Denied."
	case 32770:
		return "Not Supported."
	case 32771:
		return "Status is Unknown."
	case 32772:
		return "Timeout."
	case 32773:
		return "Invalid Parameter."
	case 32774:
		return "System is In Use."
	case 32775:
		return "Invalid State for this Operation."
	case 32776:
		return "Incorrect Data Type."
	case 32777:
		return "System is Not Available."
	case 32778:
		return "Out of Memory."
	default:
		return "The Method Failed. The Reason is Unknown."
	}
}

// validateOutput validates the output of a WMI method call.
func (v *VMMS) ValidateOutput(output *wmi.Result) error {
	returnValue, err := output.GetUint32("ReturnValue")
	if err != nil {
		return fmt.Errorf("failed to get return value: %w", err)
	}

	if returnValue == 4096 {
		// Job started - wait for completion
		jobPath, err := output.GetString("Job")
		if err != nil {
			return fmt.Errorf("failed to get job path: %w", err)
		}

		job, err := v.virtualizationConn.Get(jobPath)
		if err != nil {
			return fmt.Errorf("failed to get job object: %w", err)
		}

		for {
			jobState, err := job.GetUint16("JobState")
			if err != nil {
				return fmt.Errorf("failed to get job state: %w", err)
			}

			if common.IsJobComplete(jobState) {
				if !common.IsJobSuccessful(jobState) {
					errorDesc, err := job.GetString("ErrorDescription")
					if err != nil || errorDesc == "" {
						return fmt.Errorf("job failed: %s", ErrorCodeMeaning(uint32(jobState)))
					}
					return fmt.Errorf(errorDesc)
				}
				break
			}

			time.Sleep(500 * time.Millisecond)
			job, err = v.virtualizationConn.Get(jobPath)
			if err != nil {
				return fmt.Errorf("failed to refresh job object: %w", err)
			}
		}
	} else if returnValue != 0 {
		return fmt.Errorf(ErrorCodeMeaning(returnValue))
	}

	return nil
}

// CreateResource creates a resource of the specified type.
func (v *VMMS) CreateResource(resource common.Resource) (*wmi.Result, error) {
	resourcePoolClass := "Msvm_ResourcePool"
	if resource == common.ResourceProcessor {
		resourcePoolClass = "Msvm_ProcessorPool"
	}

	resourceSubType := common.ResourceSubType(resource)

	// Query for the resource pool
	query := fmt.Sprintf("SELECT * FROM %s WHERE ResourceSubType = '%s' AND Primordial = TRUE", resourcePoolClass, resourceSubType)
	resourcePools, err := v.virtualizationConn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query resource pool: %w", err)
	}

	if len(resourcePools) == 0 {
		return nil, fmt.Errorf("no resource pool found for resource subtype %s", resourceSubType)
	}

	// Get allocation capabilities
	capQuery := fmt.Sprintf("ASSOCIATORS OF {%s} WHERE ResultClass=Msvm_AllocationCapabilities", resourcePools[0].Path())
	caps, err := v.virtualizationConn.Query(capQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query allocation capabilities: %w", err)
	}

	if len(caps) == 0 {
		return nil, fmt.Errorf("no allocation capabilities found for resource subtype %s", resourceSubType)
	}

	// Get default settings
	settingQuery := fmt.Sprintf("ASSOCIATORS OF {%s} WHERE ResultClass=Msvm_SettingsDefineCapabilities", caps[0].Path())
	settings, err := v.virtualizationConn.Query(settingQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query settings: %w", err)
	}

	var defaultSettingPath string
	for _, setting := range settings {
		valueRole, err := setting.GetUint16("ValueRole")
		if err != nil {
			continue
		}
		if valueRole == 0 {
			defaultSettingPath, err = setting.GetString("PartComponent")
			if err != nil {
				continue
			}
			break
		}
	}

	if defaultSettingPath == "" {
		return nil, fmt.Errorf("unable to find the Default Resource Settings")
	}

	return v.virtualizationConn.Get(defaultSettingPath)
}
