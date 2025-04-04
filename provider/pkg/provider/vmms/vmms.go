// Package vmms provides functionality to interact with the Hyper-V Virtual Machine Management Service (VMMS).
package vmms

import (
	"fmt"
	"log"

	"github.com/microsoft/wmi/pkg/base/host"
	securitysvc "github.com/microsoft/wmi/pkg/virtualization/core/security/service"
	vmmsvc "github.com/microsoft/wmi/pkg/virtualization/core/service"
	imsvc "github.com/microsoft/wmi/pkg/virtualization/core/storage/service"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance" // Updated import path
)

// VMMS represents the Hyper-V Virtual Machine Management Service.
type VMMS struct {
	host                *host.WmiHost
	virtualizationConn  *wmi.WmiSession
	hgsConn             *wmi.WmiSession
	securityService     *securitysvc.SecurityService
	imageManagementSvc  *imsvc.ImageManagementService
	vmManagementService *vmmsvc.VirtualSystemManagementService
}

// NewVMMS creates a new VMMS instance.
func NewVMMS(host *host.WmiHost) (*VMMS, error) {
	vmms := &VMMS{
		host: host,
	}

	sm := wmi.NewWmiSessionManager()
	defer sm.Close()
	defer sm.Dispose()

	// Set up virtualization connection
	virtConn, err := sm.GetLocalSession("root\\virtualization\\v2")
	if err != nil {
		return nil, fmt.Errorf("failed to create virtualization connection: %w", err)
	}
	_, err = virtConn.Connect()
	if err != nil {
		log.Printf("Could not connect session %v", err)
		return nil, fmt.Errorf("failed to connect to virtconn virtualization namespace: %w", err)
	}
	defer virtConn.Close()
	defer virtConn.Dispose()
	vmms.virtualizationConn = virtConn

	// Set up HGS connection
	hgsConn, err := sm.GetLocalSession("root\\Microsoft\\Windows\\Hgs")
	if err != nil {
		return nil, fmt.Errorf("failed to create HGS connection: %w", err)
	}
	_, err = hgsConn.Connect()
	if err != nil {
		log.Printf("Could not connect session %v", err)
		return nil, fmt.Errorf("failed to connect to hgs virtualization namespace: %w", err)
	}
	defer hgsConn.Close()
	defer hgsConn.Dispose()
	vmms.hgsConn = hgsConn

	// Get services
	ss, err := securitysvc.GetSecurityService(vmms.virtualizationConn.WMIHost)
	if err != nil {
		return nil, err
	}
	vmms.securityService = ss

	ims, err := imsvc.GetImageManagementService(vmms.virtualizationConn.WMIHost)
	if err != nil {
		return nil, err
	}
	vmms.imageManagementSvc = ims

	vmmSvc, err := vmmsvc.GetVirtualSystemManagementService(vmms.virtualizationConn.WMIHost)
	if err != nil {
		return nil, err
	}
	vmms.vmManagementService = vmmSvc

	return vmms, nil
}

// GetVirtualizationConn returns the virtualization connection.
func (v *VMMS) GetVirtualizationConn() *wmi.WmiSession {
	return v.virtualizationConn
}

// GetHgsConn returns the HGS connection.
func (v *VMMS) GetHgsConn() *wmi.WmiSession {
	return v.hgsConn
}

// GetSecurityService returns the security service.
func (v *VMMS) GetSecurityService() *securitysvc.SecurityService {
	return v.securityService
}

// GetImageManagementService returns the image management service.
func (v *VMMS) GetImageManagementService() *imsvc.ImageManagementService {
	return v.imageManagementSvc
}

// GetVirtualSystemManagementService returns the virtual machine management service.
func (v *VMMS) GetVirtualSystemManagementService() *vmmsvc.VirtualSystemManagementService {
	return v.vmManagementService
}

// GetUntrustedGuardian gets the untrusted guardian.
// func (v *VMMS) GetUntrustedGuardian() (*wmi.Result, error) {
// 	query := "SELECT * FROM MSFT_HgsGuardian WHERE Name = 'UntrustedGuardian'"
// 	guardians, err := v.hgsConn.Query(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query guardians: %w", err)
// 	}

// 	if len(guardians) == 0 {
// 		return nil, nil
// 	}

// 	return guardians[0], nil
// }

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

// RequestedState represents the state to request for a virtual machine.
type RequestedState uint16

const (
	// RequestedStateOther represents another state.
	RequestedStateOther RequestedState = 1
	// RequestedStateEnabled represents an enabled state.
	RequestedStateEnabled RequestedState = 2
	// RequestedStateDisabled represents a disabled state.
	RequestedStateDisabled RequestedState = 3
	// RequestedStateShutDown represents a shutdown state.
	RequestedStateShutDown RequestedState = 4
	// RequestedStateOffline represents an offline state.
	RequestedStateOffline RequestedState = 6
	// RequestedStateTest represents a test state.
	RequestedStateTest RequestedState = 7
	// RequestedStateDefer represents a deferred state.
	RequestedStateDefer RequestedState = 8
	// RequestedStateQuiesce represents a quiesced state.
	RequestedStateQuiesce RequestedState = 9
	// RequestedStateReboot represents a reboot state.
	RequestedStateReboot RequestedState = 10
	// RequestedStateReset represents a reset state.
	RequestedStateReset RequestedState = 11
	// RequestedStateSaving represents a saving state.
	RequestedStateSaving RequestedState = 32773
	// RequestedStatePausing represents a pausing state.
	RequestedStatePausing RequestedState = 32776
	// RequestedStateResuming represents a resuming state.
	RequestedStateResuming RequestedState = 32777
	// RequestedStateFastSaved represents a fast saved state.
	RequestedStateFastSaved RequestedState = 32779
	// RequestedStateFastSaving represents a fast saving state.
	RequestedStateFastSaving RequestedState = 32780
	// RequestedStateRunningCritical represents a running critical state.
	RequestedStateRunningCritical RequestedState = 32781
	// RequestedStateOffCritical represents an off critical state.
	RequestedStateOffCritical RequestedState = 32782
	// RequestedStateStoppingCritical represents a stopping critical state.
	RequestedStateStoppingCritical RequestedState = 32783
	// RequestedStateSavedCritical represents a saved critical state.
	RequestedStateSavedCritical RequestedState = 32784
	// RequestedStatePausedCritical represents a paused critical state.
	RequestedStatePausedCritical RequestedState = 32785
	// RequestedStateStartingCritical represents a starting critical state.
	RequestedStateStartingCritical RequestedState = 32786
	// RequestedStateResetCritical represents a reset critical state.
	RequestedStateResetCritical RequestedState = 32787
	// RequestedStateSavingCritical represents a saving critical state.
	RequestedStateSavingCritical RequestedState = 32788
	// RequestedStatePausingCritical represents a pausing critical state.
	RequestedStatePausingCritical RequestedState = 32789
	// RequestedStateResumingCritical represents a resuming critical state.
	RequestedStateResumingCritical RequestedState = 32790
	// RequestedStateFastSavedCritical represents a fast saved critical state.
	RequestedStateFastSavedCritical RequestedState = 32791
	// RequestedStateFastSavingCritical represents a fast saving critical state.
	RequestedStateFastSavingCritical RequestedState = 32792
)

// RequestStateChange requests a state change for a virtual machine.
// func RequestStateChange(v VMMS, virtualMachine *vmmsvc.VirtualSystemManagementService, requestedState RequestedState) error {
// 	params := map[string]interface{}{
// 		"RequestedState": uint16(requestedState),
// 	}

// 	result, err := virtualMachine.InvokeMethod("RequestStateChange", params)
// 	if err != nil {
// 		return fmt.Errorf("failed to request state change: %w", err)
// 	}

// 	return v.ValidateOutput(result)
// }

// validateOutput validates the output of a WMI method call.
// func (v *VMMS) ValidateOutput(output *wmi.Result) error {
// 	returnValue, err := output.GetUint32("ReturnValue")
// 	if err != nil {
// 		return fmt.Errorf("failed to get return value: %w", err)
// 	}

// 	if returnValue == 4096 {
// 		// Job started - wait for completion
// 		jobPath, err := output.GetString("Job")
// 		if err != nil {
// 			return fmt.Errorf("failed to get job path: %w", err)
// 		}

// 		job, err := v.virtualizationConn.Get(jobPath)
// 		if err != nil {
// 			return fmt.Errorf("failed to get job object: %w", err)
// 		}

// 		for {
// 			jobState, err := job.GetUint16("JobState")
// 			if err != nil {
// 				return fmt.Errorf("failed to get job state: %w", err)
// 			}

// 			if IsJobComplete(jobState) {
// 				if !IsJobSuccessful(jobState) {
// 					errorDesc, err := job.GetString("ErrorDescription")
// 					if err != nil || errorDesc == "" {
// 						return fmt.Errorf("job failed: %s", ErrorCodeMeaning(uint32(jobState)))
// 					}
// 					return fmt.Errorf(errorDesc)
// 				}
// 				break
// 			}

// 			time.Sleep(500 * time.Millisecond)
// 			job, err = v.virtualizationConn.Get(jobPath)
// 			if err != nil {
// 				return fmt.Errorf("failed to refresh job object: %w", err)
// 			}
// 		}
// 	} else if returnValue != 0 {
// 		return fmt.Errorf(ErrorCodeMeaning(returnValue))
// 	}

// 	return nil
// }
