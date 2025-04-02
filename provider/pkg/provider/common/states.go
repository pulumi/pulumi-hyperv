package common

import (
	"fmt"

	"github.com/microsoft/wmi"
	"github.com/pulumi/pulumi-hyperv-provider/provider/pkg/provider"
)

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
func RequestStateChange(v provider.VMMS, virtualMachine *wmi.Result, requestedState RequestedState) error {
	params := map[string]interface{}{
		"RequestedState": uint16(requestedState),
	}

	result, err := virtualMachine.InvokeMethod("RequestStateChange", params)
	if err != nil {
		return fmt.Errorf("failed to request state change: %w", err)
	}

	return v.ValidateOutput(result)
}
