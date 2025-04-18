// Package vmms provides functionality to interact with the Hyper-V Virtual Machine Management Service (VMMS).
package vmms

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/microsoft/wmi/pkg/base/host"
	securitysvc "github.com/microsoft/wmi/pkg/virtualization/core/security/service"
	vmmsvc "github.com/microsoft/wmi/pkg/virtualization/core/service"
	imsvc "github.com/microsoft/wmi/pkg/virtualization/core/storage/service"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	"github.com/microsoft/wmi/pkg/virtualization/network/virtualswitch"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance" // Updated import path
	"github.com/pulumi/pulumi-hyperv/provider/pkg/provider/logging"
	"github.com/pulumi/pulumi-hyperv/provider/pkg/provider/util"
)

// VMMS represents the Hyper-V Virtual Machine Management Service.
type VMMS struct {
	host                *host.WmiHost
	virtualizationConn  *wmi.WmiSession
	hgsConn             *wmi.WmiSession
	securityService     *securitysvc.SecurityService
	imageManagementSvc  *imsvc.ImageManagementService
	vmManagementService *vmmsvc.VirtualSystemManagementService
	logger              logging.Logger
}

// NewVMMS creates a new VMMS instance.
func NewVMMS(ctx context.Context, host *host.WmiHost) (*VMMS, error) {
	logger := logging.GetLogger(ctx)

	// Check if host is nil
	if host == nil {
		logger.Errorf("WMI host is nil")
		return nil, fmt.Errorf("WMI host is nil")
	}

	vmms := &VMMS{
		host:   host,
		logger: logger,
	}

	// Each session manager operation is wrapped in panic recovery
	var sm *wmi.WmiSessionManager
	var smErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("Recovered from panic in NewWmiSessionManager: %v", r)
				smErr = fmt.Errorf("panic in NewWmiSessionManager: %v", r)
			}
		}()

		sm = wmi.NewWmiSessionManager()
	}()

	if smErr != nil {
		return nil, smErr
	}

	if sm == nil {
		return nil, fmt.Errorf("failed to create WMI session manager")
	}

	defer func() {
		if sm != nil {
			sm.Close()
			sm.Dispose()
		}
	}()

	// Set up virtualization connection with robust error handling
	var virtConn *wmi.WmiSession
	var virtConnErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("Recovered from panic in GetLocalSession: %v", r)
				virtConnErr = fmt.Errorf("panic in GetLocalSession: %v", r)
			}
		}()

		virtConn, virtConnErr = sm.GetLocalSession("root\\virtualization\\v2")
		if virtConnErr == nil && virtConn != nil {
			// Try to connect, but handle errors gracefully
			_, connectErr := virtConn.Connect()
			if connectErr != nil {
				logger.Errorf("Could not connect session: %v", connectErr)
				virtConnErr = fmt.Errorf("failed to connect to virtualization namespace: %w", connectErr)

				// Close and dispose if connection failed but session was created
				virtConn.Close()
				virtConn.Dispose()
				virtConn = nil
			}
		}
	}()

	if virtConnErr != nil {
		logger.Errorf("Failed to create virtualization connection: %v", virtConnErr)
		return nil, fmt.Errorf("failed to create virtualization connection: %w", virtConnErr)
	}

	if virtConn == nil {
		logger.Errorf("Virtualization connection is nil")
		return nil, fmt.Errorf("virtualization connection is nil after creation")
	}

	defer func() {
		if virtConn != nil {
			virtConn.Close()
			virtConn.Dispose()
		}
	}()

	vmms.virtualizationConn = virtConn

	// Set up HGS connection (optional - not needed for basic Hyper-V functionality)
	func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Warnf("Recovered from panic in HGS connection setup: %v", r)
				logger.Infof("Continuing without HGS support (required only for advanced security features)")
			}
		}()

		hgsConn, err := sm.GetLocalSession("root\\Microsoft\\Windows\\Hgs")
		if err != nil {
			logger.Warnf("HGS connection not available: %v", err)
			logger.Infof("Continuing without HGS support (required only for advanced security features)")
			return
		}

		if hgsConn == nil {
			logger.Warnf("HGS session is nil")
			return
		}

		_, err = hgsConn.Connect()
		if err != nil {
			logger.Warnf("Could not connect to HGS session: %v", err)
			logger.Infof("Continuing without HGS support (required only for advanced security features)")
			hgsConn.Close()
			hgsConn.Dispose()
			return
		}

		// Only store the HGS connection if everything succeeded
		vmms.hgsConn = hgsConn

		// Don't defer close here since we want to keep this connection
		// It will be closed when the VMMS instance is disposed
	}()

	// Get security service (optional - not required for basic functionality)
	func() {
		// Skip if virtualization connection is not available
		if vmms.virtualizationConn == nil || vmms.virtualizationConn.WMIHost == nil {
			logger.Warnf("Virtualization connection or WMI host is nil, skipping SecurityService")
			return
		}

		defer func() {
			if r := recover(); r != nil {
				logger.Warnf("Recovered from panic in GetSecurityService: %v", r)
				logger.Infof("Continuing without security service (needed only for advanced security features)")
			}
		}()

		ss, err := securitysvc.GetSecurityService(vmms.virtualizationConn.WMIHost)
		if err != nil {
			logger.Warnf("Could not get security service: %v", err)
			logger.Infof("Continuing without security service (needed only for advanced security features)")
			return
		}

		if ss == nil {
			logger.Warnf("Security service is nil")
			return
		}

		vmms.securityService = ss
	}()

	// Get image management service (optional - not required for all functionality)
	func() {
		// Skip if virtualization connection is not available
		if vmms.virtualizationConn == nil || vmms.virtualizationConn.WMIHost == nil {
			logger.Warnf("Virtualization connection or WMI host is nil, skipping ImageManagementService")
			return
		}

		defer func() {
			if r := recover(); r != nil {
				logger.Warnf("Recovered from panic in GetImageManagementService: %v", r)
				logger.Infof("Continuing without image management service")
			}
		}()

		ims, err := imsvc.GetImageManagementService(vmms.virtualizationConn.WMIHost)
		if err != nil {
			logger.Warnf("Could not get image management service: %v", err)
			logger.Infof("Continuing without image management service (needed for some disk operations)")
			return
		}

		if ims == nil {
			logger.Warnf("Image management service is nil")
			return
		}

		vmms.imageManagementSvc = ims
	}()

	// Get Virtual System Management Service (critical service)
	// Even though this is considered critical, we'll make our code robust enough to
	// continue without it and handle the nil value in each resource's controller
	var vsmsErr error

	func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("Recovered from panic in GetVirtualSystemManagementService: %v", r)
				if err, ok := r.(error); ok {
					vsmsErr = err
				} else {
					vsmsErr = fmt.Errorf("panic in GetVirtualSystemManagementService: %v", r)
				}
			}
		}()

		// Skip if virtualization connection is not available
		if vmms.virtualizationConn == nil || vmms.virtualizationConn.WMIHost == nil {
			logger.Errorf("Virtualization connection or WMI host is nil, cannot get VirtualSystemManagementService")
			vsmsErr = fmt.Errorf("virtualization connection or WMI host is nil, cannot get VirtualSystemManagementService")
			return
		}

		vmmSvc, err := vmmsvc.GetVirtualSystemManagementService(vmms.virtualizationConn.WMIHost)
		if err != nil {
			logger.Errorf("Failed to get virtual system management service: %v", err)
			logger.Warnf("Continuing without VSMS - some operations will not be available")
			vsmsErr = fmt.Errorf("failed to get virtual system management service: %w", err)
			return
		}

		if vmmSvc == nil {
			logger.Errorf("Virtual system management service is nil")
			logger.Warnf("Continuing without VSMS - some operations will not be available")
			vsmsErr = fmt.Errorf("virtual system management service is nil")
			return
		}

		vmms.vmManagementService = vmmSvc
	}()

	// Log the warning but return the VMMS object anyway so controllers can handle this case
	if vsmsErr != nil {
		logger.Warnf("Failed to initialize VirtualSystemManagementService: %v", vsmsErr)
		logger.Warnf("Continuing with limited functionality - some operations may fail")
	}

	// Return VMMS instance even with limited functionality
	// Each controller will need to handle nil services appropriately
	return vmms, nil
}

// GetVirtualizationConn returns the virtualization connection.
func (v *VMMS) GetVirtualizationConn() *wmi.WmiSession {
	// Add nil check to prevent panic when VMMS is nil
	if v == nil {
		log.Printf("[ERROR] VMMS object is nil when trying to get virtualization connection")
		return nil
	}
	return v.virtualizationConn
}

// GetHgsConn returns the HGS connection or nil if not available.
// Callers must check for nil before using the returned connection.
func (v *VMMS) GetHgsConn() *wmi.WmiSession {
	return v.hgsConn
}

// GetSecurityService returns the security service or nil if not available.
// Callers must check for nil before using the returned service.
func (v *VMMS) GetSecurityService() *securitysvc.SecurityService {
	return v.securityService
}

// GetImageManagementService returns the image management service or nil if not available.
// Callers must check for nil before using the returned service.
func (v *VMMS) GetImageManagementService() *imsvc.ImageManagementService {
	if v == nil {
		v.logger.Errorf("VMMS object is nil when trying to get ImageManagementService")
		return nil
	}
	if v.imageManagementSvc == nil {
		v.logger.Warnf("ImageManagementService is not initialized")
		return nil
	}
	return v.imageManagementSvc
}

// GetVirtualSystemManagementService returns the virtual machine management service.
func (v *VMMS) GetVirtualSystemManagementService() *vmmsvc.VirtualSystemManagementService {
	// Add nil check to prevent panics when the service couldn't be initialized
	if v == nil {
		v.logger.Errorf("VMMS object is nil when trying to get VirtualSystemManagementService")
		return nil
	}
	return v.vmManagementService
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

func (v *VMMS) AttachVirtualHardDisk(vm *virtualsystem.VirtualMachine, hdPath string, controllerType string, controllerNumber int, controllerLocation int, logger logging.Logger) error {
	if v == nil {
		return fmt.Errorf("VMMS object is nil")
	}

	if vm == nil {
		return fmt.Errorf("virtual machine is nil")
	}

	vsms := v.GetVirtualSystemManagementService()
	if vsms == nil {
		logger.Warnf("VirtualSystemManagementService is unavailable, falling back to PowerShell")
		return attachVirtualHardDiskPowerShell(vm, hdPath, controllerType, controllerNumber, controllerLocation, logger)
	}

	// Determine disk type based on controller type
	var diskType virtualsystem.VirtualHardDiskType
	switch strings.ToUpper(controllerType) {
	case "IDE":
		diskType = virtualsystem.VirtualHardDiskType_DATADISK_VIRTUALHARDDISK
	case "SCSI":
		diskType = virtualsystem.VirtualHardDiskType_DATADISK_VIRTUALHARDDISK
	default:
		logger.Warnf("Unknown controller type '%s', defaulting to SCSI", controllerType)
		diskType = virtualsystem.VirtualHardDiskType_DATADISK_VIRTUALHARDDISK
	}

	// Attempt to attach using WMI API
	err := vsms.AddSCSIController(vm)
	if err != nil {
		logger.Warnf("Failed to add SCSI controller: %v", err)
		return attachVirtualHardDiskPowerShell(vm, hdPath, controllerType, controllerNumber, controllerLocation, logger)
	}
	_, _, err = vsms.AttachVirtualHardDisk(vm, hdPath, diskType)
	if err == nil {
		logger.Infof("[INFO] Successfully attached VHD [%s] using WMI", hdPath)
		return nil
	}

	logger.Warnf("Failed to attach VHD [%s] using WMI: %v, trying direct API", hdPath, err)

	// Attempt direct API call
	err = v.AttachVirtualHardDiskDirectApi(vm, hdPath, controllerNumber, controllerLocation, logger)
	if err == nil {
		logger.Infof("[INFO] Successfully attached VHD [%s] using direct API", hdPath)
		return nil
	}

	logger.Warnf("Failed to attach VHD [%s] using direct API: %v, falling back to PowerShell", hdPath, err)

	// Fallback to PowerShell
	return attachVirtualHardDiskPowerShell(vm, hdPath, controllerType, controllerNumber, controllerLocation, logger)
}

// attachVirtualHardDiskPowerShell attaches a VHD using PowerShell as a fallback.
func attachVirtualHardDiskPowerShell(vm *virtualsystem.VirtualMachine, hdPath string, controllerType string, controllerNumber int, controllerLocation int, logger logging.Logger) error {
	vmName, err := vm.GetPropertyElementName()
	if err != nil {
		return fmt.Errorf("failed to get VM name: %w", err)
	}

	cmd := fmt.Sprintf("Add-VMHardDiskDrive -VMName \"%s\" -Path \"%s\" -ControllerType %s -ControllerNumber %d -ControllerLocation %d",
		vmName, hdPath, controllerType, controllerNumber, controllerLocation)

	output, err := util.RunPowerShellCommand(cmd)
	if err != nil {
		outputStr := string(output)

		// Check for common PowerShell error patterns and provide more user-friendly messages
		if strings.Contains(outputStr, "ObjectNotFound") && strings.Contains(outputStr, "Add-VMHardDiskDrive") {
			return fmt.Errorf("failed to attach disk: virtual machine '%s' not found. Please verify the VM exists and you have permission to modify it", vmName)
		}

		if strings.Contains(outputStr, "no available locations were found on the disk controller") {
			return fmt.Errorf("failed to attach disk: no available locations found on %s controller %d. Try using a different controller number or location, or add a new controller",
				controllerType, controllerNumber)
		}

		if strings.Contains(outputStr, "The system cannot find the file specified") || strings.Contains(outputStr, "Cannot find path") {
			return fmt.Errorf("failed to attach disk: VHD file '%s' not found. Please verify the path is correct and accessible", hdPath)
		}

		if strings.Contains(outputStr, "The parameter is incorrect") {
			return fmt.Errorf("failed to attach disk: incorrect parameter. This often happens with incompatible VHDX formats or block sizes. Try a different VHD file or format")
		}

		if strings.Contains(outputStr, "Access is denied") || strings.Contains(outputStr, "AccessDenied") {
			return fmt.Errorf("failed to attach disk: access denied. Please verify you have administrator privileges")
		}

		if strings.Contains(outputStr, "The requested operation cannot be performed on a file with a user-mapped section open") {
			return fmt.Errorf("failed to attach disk: the VHD file '%s' is in use by another process. Make sure the disk is not mounted elsewhere", hdPath)
		}

		// Default error with full details if we couldn't match a specific pattern
		return fmt.Errorf("failed to attach VHD using PowerShell: %v\nDetailed error: %s", err, outputStr)
	}

	logger.Infof("[INFO] Successfully attached VHD [%s] to VM [%s] using PowerShell", hdPath, vmName)
	return nil
}

func (v *VMMS) AttachVirtualHardDiskDirectApi(vm *virtualsystem.VirtualMachine, path string, controllerNumber int, controllerLocation int, logger logging.Logger) error {
	if v == nil {
		return fmt.Errorf("VMMS object is nil")
	}

	if vm == nil {
		return fmt.Errorf("virtual machine is nil")
	}

	if v.vmManagementService == nil {
		return fmt.Errorf("virtual system management service is nil")
	}

	// Get VM element name for identifying the VM
	vmName, err := vm.GetPropertyElementName()
	if err != nil {
		return fmt.Errorf("failed to get VM element name: %w", err)
	}

	// Create resource settings for the hard drive
	resourceSettings := []interface{}{
		map[string]interface{}{
			"ResourceType":       uint16(31), // 31 = Disk drive
			"ResourceSubType":    "Microsoft:Hyper-V:Virtual Hard Disk",
			"Path":               path,
			"ControllerType":     "Microsoft:Hyper-V:Synthetic SCSI Controller",
			"ControllerNumber":   uint32(controllerNumber),
			"ControllerLocation": uint32(controllerLocation),
		},
	}

	// Get the VM path or system name to use in AddResourceSettings
	systemName := fmt.Sprintf("\\\\%s\\root\\virtualization\\v2:Msvm_ComputerSystem.CreationClassName=\"Msvm_ComputerSystem\",Name=\"%s\"",
		v.host.HostName, vm.InstanceID)

	// Add the resource settings to the VM
	result, err := v.vmManagementService.InvokeMethod("AddResourceSettings", []interface{}{systemName, resourceSettings})
	if err != nil {
		return fmt.Errorf("failed to add hard drive: %w", err)
	}

	// Check return value
	if len(result) == 0 {
		return fmt.Errorf("empty result from AddResourceSettings")
	}

	// The first element should be a map with ReturnValue
	resultMap, ok := result[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected result type from AddResourceSettings")
	}

	// Check return value
	returnValue, ok := resultMap["ReturnValue"].(uint32)
	if !ok {
		return fmt.Errorf("unexpected ReturnValue type from AddResourceSettings")
	}

	if returnValue != 0 {
		return fmt.Errorf("AddResourceSettings failed with error: %s", ErrorCodeMeaning(returnValue))
	}

	logger.Infof("[INFO] Successfully attached virtual hard disk %s to VM %s", path, vmName)
	return nil
}

func (v *VMMS) AddVirtualNetworkAdapterAndConnect(vm *virtualsystem.VirtualMachine, adapterName string, switchName string, logger logging.Logger) error {
	if v == nil {
		return fmt.Errorf("VMMS object is nil")
	}

	if vm == nil {
		return fmt.Errorf("virtual machine is nil")
	}

	vsms := v.GetVirtualSystemManagementService()
	if vsms == nil {
		logger.Warnf("VirtualSystemManagementService is unavailable, falling back to PowerShell")
		return addVirtualNetworkAdapterPowerShell(vm, adapterName, switchName, logger)
	}

	// Attempt to add adapter using WMI API
	var addErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				addErr = fmt.Errorf("recovered from panic in AddVirtualNetworkAdapter: %v", r)
				logger.Errorf("Recovered from panic in AddVirtualNetworkAdapter: %v", r)
			}
		}()

		adapter, err := vsms.AddVirtualNetworkAdapter(vm, adapterName)
		if err != nil {
			addErr = fmt.Errorf("failed to add network adapter via WMI: %v", err)
			return
		}
		defer adapter.Close()

		// Connect adapter to switch
		vs, err := virtualswitch.GetVirtualSwitch(v.virtualizationConn.WMIHost, switchName)
		if err != nil {
			logger.Errorf("Failed [%+v]", err)
		}
		defer vs.Close()
		logger.Infof("Got VirtualSwitch[%s]", "test")
		err = vsms.ConnectAdapterToVirtualSwitch(vm, adapterName, vs)
		if err != nil {
			addErr = fmt.Errorf("failed to connect network adapter to switch via WMI: %v", err)
			return
		}
	}()

	if addErr == nil {
		logger.Infof("Successfully added and connected network adapter [%s] to switch [%s] using WMI", adapterName, switchName)
		return nil
	}

	logger.Warnf("Failed to add/connect network adapter [%s] using WMI: %v, falling back to PowerShell", adapterName, addErr)

	// Fallback to PowerShell
	return addVirtualNetworkAdapterPowerShell(vm, adapterName, switchName, logger)
}

// addVirtualNetworkAdapterPowerShell adds and connects a network adapter using PowerShell as a fallback.
func addVirtualNetworkAdapterPowerShell(vm *virtualsystem.VirtualMachine, adapterName string, switchName string, logger logging.Logger) error {
	vmName, err := vm.GetPropertyElementName()
	if err != nil {
		return fmt.Errorf("failed to get VM name: %w", err)
	}
	cmd := fmt.Sprintf("Add-VMNetworkAdapter -VMName \"%s\" -Name \"%s\" -SwitchName \"%s\"",
		vmName, adapterName, switchName)
	output, err := util.RunPowerShellCommand(cmd)
	if err != nil {
		outputStr := string(output)

		// Check for common PowerShell error patterns and provide more user-friendly messages
		if strings.Contains(outputStr, "ObjectNotFound") && strings.Contains(outputStr, "Add-VMNetworkAdapter") {
			return fmt.Errorf("failed to add network adapter: virtual machine '%s' not found. Please verify the VM exists and you have permission to modify it", vmName)
		}

		if strings.Contains(outputStr, "Cannot find virtual switch") {
			return fmt.Errorf("failed to add network adapter: virtual switch '%s' not found. Please verify the switch exists", switchName)
		}

		if strings.Contains(outputStr, "already exists") {
			return fmt.Errorf("failed to add network adapter: an adapter named '%s' already exists on VM '%s'", adapterName, vmName)
		}

		if strings.Contains(outputStr, "Access is denied") || strings.Contains(outputStr, "AccessDenied") {
			return fmt.Errorf("failed to add network adapter: access denied. Please verify you have administrator privileges")
		}

		// Default error with full details if we couldn't match a specific pattern
		return fmt.Errorf("failed to add/connect network adapter using PowerShell: %v\nDetailed error: %s", err, outputStr)
	}

	logger.Infof("[INFO] Successfully added and connected network adapter [%s] to switch [%s] on VM [%s] using PowerShell", adapterName, switchName, vmName)
	return nil
}

// AddVirtualNetworkAdapterAndConnect adds a virtual network adapter to a virtual machine
// and connects it to a virtual switch.
func (v *VMMS) AddVirtualNetworkAdapterAndConnectApi(vm *virtualsystem.VirtualMachine, adapterName string, switchName string, logger logging.Logger) error {
	if v == nil {
		return fmt.Errorf("VMMS object is nil")
	}

	if vm == nil {
		return fmt.Errorf("virtual machine is nil")
	}

	if v.vmManagementService == nil {
		return fmt.Errorf("virtual system management service is nil")
	}

	// Get VM element name for identifying the VM
	vmName, err := vm.GetPropertyElementName()
	if err != nil {
		return fmt.Errorf("failed to get VM element name: %w", err)
	}

	// Create resource settings for the network adapter
	resourceSettings := []interface{}{
		map[string]interface{}{
			"ResourceType":    uint16(10), // 10 = Network adapter
			"ResourceSubType": "Microsoft:Hyper-V:Synthetic Ethernet Port",
			"ElementName":     adapterName,
			"Connection":      []string{switchName},
		},
	}

	// Get the VM path or system name to use in AddResourceSettings
	systemName := fmt.Sprintf("\\\\%s\\root\\virtualization\\v2:Msvm_ComputerSystem.CreationClassName=\"Msvm_ComputerSystem\",Name=\"%s\"",
		v.host.HostName, vm.InstanceID)

	// Add the resource settings to the VM
	result, err := v.vmManagementService.InvokeMethod("AddResourceSettings", []interface{}{systemName, resourceSettings})
	if err != nil {
		return fmt.Errorf("failed to add network adapter: %w", err)
	}

	// Check return value
	if len(result) == 0 {
		return fmt.Errorf("empty result from AddResourceSettings")
	}

	// The first element should be a map with ReturnValue
	resultMap, ok := result[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected result type from AddResourceSettings")
	}

	// Check return value
	returnValue, ok := resultMap["ReturnValue"].(uint32)
	if !ok {
		return fmt.Errorf("unexpected ReturnValue type from AddResourceSettings")
	}

	if returnValue != 0 {
		return fmt.Errorf("AddResourceSettings failed with error: %s", ErrorCodeMeaning(returnValue))
	}

	logger.Infof("[INFO] Successfully added network adapter %s to VM %s connected to switch %s", adapterName, vmName, switchName)
	return nil
}
