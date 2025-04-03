// Package vmms provides functionality to interact with the Hyper-V Virtual Machine Management Service (VMMS).
package vmms

import (
	"fmt"
	"strings"
	"time"

	"github.com/microsoft/wmi"
)

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

// JobState represents the state of a WMI job.
type JobState uint16

const (
	// JobStateNew indicates the job is new.
	JobStateNew JobState = 2
	// JobStateStarting indicates the job is starting.
	JobStateStarting JobState = 3
	// JobStateRunning indicates the job is running.
	JobStateRunning JobState = 4
	// JobStateSuspended indicates the job is suspended.
	JobStateSuspended JobState = 5
	// JobStateShuttingDown indicates the job is shutting down.
	JobStateShuttingDown JobState = 6
	// JobStateCompleted indicates the job is completed.
	JobStateCompleted JobState = 7
	// JobStateTerminated indicates the job is terminated.
	JobStateTerminated JobState = 8
	// JobStateKilled indicates the job is killed.
	JobStateKilled JobState = 9
	// JobStateException indicates the job encountered an exception.
	JobStateException JobState = 10
	// JobStateCompletedWithWarnings indicates the job completed with warnings.
	JobStateCompletedWithWarnings JobState = 32768
)

// IsJobComplete checks if a job has completed.
func IsJobComplete(jobState uint16) bool {
	state := JobState(jobState)
	return state == JobStateCompleted ||
		state == JobStateCompletedWithWarnings ||
		state == JobStateTerminated ||
		state == JobStateException ||
		state == JobStateKilled
}

// IsJobSuccessful checks if a job has completed successfully.
func IsJobSuccessful(jobState uint16) bool {
	state := JobState(jobState)
	return state == JobStateCompleted ||
		state == JobStateCompletedWithWarnings
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
func RequestStateChange(v VMMS, virtualMachine *wmi.Result, requestedState RequestedState) error {
	params := map[string]interface{}{
		"RequestedState": uint16(requestedState),
	}

	result, err := virtualMachine.InvokeMethod("RequestStateChange", params)
	if err != nil {
		return fmt.Errorf("failed to request state change: %w", err)
	}

	return v.ValidateOutput(result)
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

			if IsJobComplete(jobState) {
				if !IsJobSuccessful(jobState) {
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
