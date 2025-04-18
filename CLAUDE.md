# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Contributing to the Pulumi HyperV Provider

First, thank you for your interest in contributing to the Pulumi HyperV Provider!

### Code of Conduct

Please read our [Code of Conduct](CODE-OF-CONDUCT.md) before participating in this project.

### Development Environment Setup

#### Prerequisites

- [Go 1.24 or later](https://golang.org/dl/)
- [NodeJS 16.X.X or later](https://nodejs.org/en/download/)
- [Python 3.8 or later](https://www.python.org/downloads/)
- [.NET Core 6.0 or later](https://dotnet.microsoft.com/download)
- [PowerShell 7 or later](https://github.com/PowerShell/PowerShell/releases) (required for Windows builds)
- [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)

## Build Commands

- Windows builds use PowerShell: `.\make.ps1 [target]`
- Linux/macOS builds use Make: `make [target]`
- Build everything: `make build` or `.\make.ps1 build`
- Provider only: `make provider` or `.\make.ps1 provider`
- Debug build: `make provider_debug` or `.\make.ps1 provider_debug`
- Generate SDKs: `make codegen` or `.\make.ps1 codegen`
- Build specific SDK: `make [dotnet|go|nodejs|python|java]_sdk` or `.\make.ps1 [dotnet|go|nodejs|python|java]_sdk`
  - Note: Each *_sdk target checks if the schema file exists and generates it if needed
  - The provider binary is built only if needed
- Install: `make install` or `.\make.ps1 install`

## Windows-Specific Options

- Force Windows mode: `.\make.ps1 -ForceWindowsMode $true [target]`
- HyperV provider requires Windows for actual implementation testing
- Hyper-V detection: The provider automatically detects if Hyper-V is available at startup
  - Logs warnings if Hyper-V is not available or properly configured
  - Hyper-V detection is implemented in `provider/pkg/provider/util/hyperv_detection.go`
  - Detection checks for Windows OS, WMI connectivity, and Hyper-V service availability

## Lint/Test Commands

- Run linter: `.\make.ps1 lint` (Windows) or `make lint` (Linux/macOS)
- Run all tests: `.\make.ps1 test` (Windows) or `make test` (Linux/macOS)
- Provider tests only: `.\make.ps1 test_provider` (Windows) or `make test_provider` (Linux/macOS)
- Run a single test: `cd provider && GOOS=windows go test -v -count=1 ./... -run TestName`
- Format Go code: `gofmt -w .` or `make format`
- Run Go linter: `make lint` (ALWAYS preferred over direct golangci-lint calls)

## Post-Code Change Commands

Run these commands after every change to Go files:

- Format code: `make format` (preferred) or `gofmt -w .`
- Lint code: `make lint` (ALWAYS preferred over direct golangci-lint calls)
- Build provider: `make provider` (preferred) or `cd provider && GOOS=windows go build -o bin/pulumi-resource-hyperv ./cmd/pulumi-resource-hyperv`

For Markdown files:

- Lint Markdown: `npx markdownlint "**/*.md"` or `make lint`
- Fix most issues automatically: `npx markdownlint --fix "**/*.md"`

## Example Tests

- Examples are in the `examples/` directory
- Tests for examples use Pulumi's integration testing framework
- Advanced testing uses `github.com/pulumi/providertest/pulumitest` package
- IMPORTANT: Example tests will ONLY work on Windows systems because Hyper-V is Windows-only
- Do NOT attempt to run example tests on Linux/macOS as they will always fail
- Run example tests (Windows only): `cd examples && go test -v ./...`
- Run specific example test (Windows only): `cd examples && go test -v -run TestDevEnvironmentTypeScript`
- Example types:
  - Single-resource examples: Found in directories named after resources (machine, vhdfile, etc.)
  - Multi-resource examples:
    - `simple-all-four/`: Basic example showing all four resource types together
    - `devenv/`: Advanced example showing a development environment with multiple VMs
- Each example has:
  - `index.ts` - TypeScript implementation
  - `Pulumi.yaml` - Project configuration
  - `package.json` - Node.js dependencies

### Examples Structure and Dependencies

- When creating a new example, refer to existing examples for correct structure
- Package.json should use `"main": "index"` instead of `"main": "index.js"` to support TypeScript
- SDK Import Paths:
  - The SDK supports both direct and namespaced resource imports:
    - Namespaced imports (recommended): `hyperv.machine.Machine`, `hyperv.virtualswitch.VirtualSwitch`
    - Direct imports (legacy): `hyperv.Machine`, `hyperv.VirtualSwitch`
  - Available namespaces: `machine`, `virtualswitch`, `vhdfile`, `networkadapter`, `config`, `types`
  - Example `devenv/` uses namespaced imports while `simple-all-four/` uses direct imports
  - For new code, prefer namespaced imports for better type safety and clarity
- Resource dependencies:
  - NetworkAdapters require `vmName` property when creating network adapters
  - VhdFile as differencing disk requires `parentPath` (parent VHD) and `diskType: "Differencing"` but no `sizeBytes`
  - VhdFile creation strongly benefits from using `blockSize: 1048576` (1MB) for better compatibility
  - Machine with dynamic memory uses `dynamicMemory: true`, `minimumMemory` and `maximumMemory`
  - Machine auto start/stop behavior configured with `autoStartAction` and `autoStopAction` properties
  - Create resources in proper order (create VMs before attaching network adapters to them)
- Example in `examples/devenv/` demonstrates setting up a full development environment
- Example in `examples/simple-all-four/` shows how to create all four resource types

## External Dependencies

- WMI Library: This project depends on the Microsoft WMI library for Hyper-V interaction
  - Repository: <https://github.com/microsoft/wmi>
  - Import path: `github.com/microsoft/wmi@v0.31.1`
  - Local location (Optional): `$GOPATH/src/github.com/microsoft/wmi`
  - Key packages used:
    - `pkg/base/host`: For WMI host connection
    - `pkg/virtualization/core/service`: For VirtualSystemManagementService
    - `pkg/virtualization/core/virtualsystem`: For VirtualMachine operations
    - `pkg/wmiinstance`: For WMI object manipulation
  - References:
    - Look at test files like `pkg/virtualization/core/service/virtualmachinemanagementservice_test.go` for usage examples

- Pulumi Go Provider: This project requires specific versions of Pulumi packages
  - Required Provider Version: `github.com/pulumi/pulumi-go-provider@v0.25.0`
  - Compatible Pulumi SDK: `github.com/pulumi/pulumi/sdk/v3@v3.160.0`
  - Go version: 1.24

## Development Workflow

### Adding New Resources

1. Create a new folder in `provider/pkg/provider` for your resource
2. Create the main resource file, controller file, and outputs file
3. Document your resource in a separate markdown file
4. Update the provider schema in `provider.go`
5. Regenerate SDKs with `.\make.ps1 codegen`
6. Add tests for your resource
7. Add an example of your resource to the examples directory

### Creating Examples

Examples demonstrate real-world use cases of the provider resources:

1. Create a new directory in `examples/` for your example
2. Add three key files:
   - `Pulumi.yaml` - Project configuration
   - `index.ts` - TypeScript implementation
   - `package.json` - Node.js dependencies
3. Add tests to `examples_nodejs_test.go` using:
   - Standard Pulumi integration testing framework
   - Advanced testing with `github.com/pulumi/providertest/pulumitest`
4. Examples should be complete, executable Pulumi programs
5. Consider advanced scenarios like the "devenv" example that shows creating multiple related resources

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests to ensure your changes work as expected
5. Submit a pull request with a clear description of your changes

When submitting a PR, include:

- A clear description of what the change does
- Any new documentation needed
- Tests that validate your change works
- Examples demonstrating the use of new features

## Code Style Guidelines

- License header: All Go files include Apache 2.0 license header
- Imports: Standard library first, external dependencies second, project imports last
- Resource pattern: Separate files for main type, inputs, outputs, and controller
  - [resource].go: Core type definitions, inputs, outputs, and annotations
  - [resource]Controller.go: Implementation of CRUD operations
  - [resource]Outputs.go: Additional output-related code
  - [resource].md: Documentation that gets embedded
- Types: Use pointer types for optional values
- Error handling: Return errors with contextual information using `fmt.Errorf`
- Test coverage: Aim for comprehensive test coverage, especially for resources
- Documentation: Document resources in separate markdown files that are embedded
- Naming: CamelCase for exported types, unexportedCamelCase for unexported
- Windows compatibility: Use path utilities that work across platforms
- Code formatting: Run `gofmt -w .` before committing changes
- Linting: When running `golangci-lint run` on a single file, type errors may appear due to missing context from other files

### Resource Implementation Guidelines

- Mark input properties as `optional` in struct tags when they're not strictly required
- For related resources (like NetworkAdapter and VM), consider which properties should be optional
- NetworkAdapter resources have special handling to support both:
  - Direct creation: With `vmName` property to attach to a specific VM
  - Reference mode: Without `vmName` when used in a Machine's `networkAdapters` property
  - Machine resource uses NetworkAdapterInputs type defined in the networkadapter package
- Type Consistency and Reuse:
  - When multiple resources share compatible properties, use shared types
  - Shared types improve API consistency and make SDK usage more predictable
  - The Machine resource's networkAdapters property uses the networkadapter.NetworkAdapterInputs type
  - This makes the SDK more consistent when working with related resources
- When using differencing disks in VhdFile, only parentPath and diskType are required (sizeBytes can be optional)
- Use `[]interface{}` for generic collections that might contain various resource types
- Auto-generated SDKs require proper namespacing in TypeScript imports:
  - Correct: `hyperv.virtualswitch.VirtualSwitch`
  - Incorrect: `hyperv.VirtualSwitch`
- When modifying schema, rebuild the provider to regenerate schemas for SDKs

## Working with Hyper-V Resources

- Use the WMI library methods for interacting with Hyper-V resources
- Always check and handle return values from WMI method calls
- Results from InvokeMethod are returned as []interface{} where the first element is a map[string]interface{}
- Common methods:
  - vsms.InvokeMethod(): For invoking WMI methods on VirtualSystemManagementService
  - vmmsClient.GetVirtualizationConn().QueryInstances(): For WMI queries
  - wmiInstance.SetProperty(): For setting properties on WMI objects

### Utility Functions

The provider includes various utility functions in the `provider/pkg/provider/util` package:

#### PowerShell Utilities

The codebase includes robust PowerShell integration for fallback operations when WMI is unavailable:

```go
// FindPowerShellExe finds the PowerShell executable (powershell.exe or pwsh.exe)
// It tries powershell.exe first, then pwsh.exe, then pwsh (for Linux/macOS).
// Returns the executable name and nil if found, or an error if no PowerShell executable is found.
func FindPowerShellExe() (string, error) {
    // Check if PowerShell is available - try powershell.exe first, then pwsh.exe
    if _, err := exec.LookPath("powershell.exe"); err == nil {
        return "powershell.exe", nil
    } else if _, err := exec.LookPath("pwsh.exe"); err == nil {
        return "pwsh.exe", nil
    } else if _, err := exec.LookPath("pwsh"); err == nil {
        return "pwsh", nil
    }
    return "", fmt.Errorf("neither powershell.exe nor pwsh.exe found in PATH, PowerShell fallback cannot be used")
}

// RunPowerShellCommand is a helper function to run PowerShell commands with proper error handling
func RunPowerShellCommand(command string) (string, error) {
    // Find PowerShell executable
    powershellExe, err := FindPowerShellExe()
    if err != nil {
        return "", err
    }
    cmd := exec.Command(powershellExe, "-Command", command)
    output, err := cmd.CombinedOutput()
    if err != nil {
        // Always return the output even when there's an error, so callers can inspect it
        return string(output), fmt.Errorf("command failed: %v, output: %s", err, string(output))
    }
    return string(output), nil
}

// ParsePowerShellError attempts to parse common PowerShell error patterns and returns a more user-friendly error
func ParsePowerShellError(cmdOutput string, cmdName string, entityType string, entityName string) error {
    if cmdOutput == "" {
        return fmt.Errorf("unknown %s error: no output received from PowerShell command", entityType)
    }

    // Common PowerShell error patterns
    switch {
    case strings.Contains(cmdOutput, "ObjectNotFound"):
        return fmt.Errorf("%s not found: '%s'. Please verify it exists and you have permission to access it", entityType, entityName)
        
    case strings.Contains(cmdOutput, "Access is denied") || strings.Contains(cmdOutput, "AccessDenied"):
        return fmt.Errorf("access denied. Please verify you have administrator privileges")
        
    case strings.Contains(cmdOutput, "The parameter is incorrect"):
        return fmt.Errorf("incorrect parameter. This often happens with incompatible formats or configurations")
    
    case strings.Contains(cmdOutput, "Unable to find a default server with Active Directory"):
        return fmt.Errorf("unable to access Active Directory. This may happen if you're not connected to a domain controller")
    
    case strings.Contains(cmdOutput, "The operation failed because of a cluster validation error"):
        return fmt.Errorf("cluster validation error. This operation may require cluster administrative privileges")
    
    case strings.Contains(cmdOutput, "The operation failed because the process hosting the server process terminated unexpectedly"):
        return fmt.Errorf("the Hyper-V service may have restarted. Please try again")
    
    default:
        // If we can't identify the error, return the raw output
        return fmt.Errorf("%s operation failed: %s", entityType, cmdOutput)
    }
}
```

These functions provide:

1. Cross-platform PowerShell support (Windows PowerShell and PowerShell Core)
2. Centralized error handling for PowerShell command execution
3. Consistent output formatting for PowerShell operations
4. A critical fallback mechanism when WMI services are unavailable
5. Enhanced error handling with pattern-matching for specific PowerShell errors
6. User-friendly error messages that provide actionable guidance
7. Access to both error code and full PowerShell output for better diagnostics

**IMPORTANT**: Always use `util.RunPowerShellCommand()` for PowerShell operations, not direct `exec.Command()` calls. This ensures consistent handling, proper error messages, and automatic executable detection even in non-standard configurations.

### Hyper-V Detection

The provider includes automatic detection of Hyper-V availability:

- Implementation in `provider/pkg/provider/util/hyperv_detection.go`
- Called during provider initialization in `provider.go`
- Performs three main checks:
  1. Checks if running on Windows (uses runtime.GOOS)
  2. Verifies access to WMI virtualization namespace
  3. Confirms Hyper-V management service is available
- Logs warnings with helpful messages if Hyper-V is unavailable
- Allows provider to initialize even without Hyper-V (but operations will fail)
- Provides OS-specific guidance based on detected Windows version
  - Specific messages for Windows 10/11 vs Server editions
  - Special handling for access denied errors with administrator privilege reminders
  - Clear instructions for enabling Hyper-V if not enabled

### Windows 10/11 Compatibility

The provider includes special handling for Windows 10/11 client systems which have different service availability
compared to Windows Server:

- Robust error handling with panic recovery to prevent crashes
- Multiple fallback methods for core operations when primary services are unavailable
- OS version detection to provide tailored guidance messages

#### Key WMI Service Handling

1. **Host Guardian Service (HGS)**
   - HGS connection is optional in VMMS initialization
   - The provider logs a warning but continues if it can't connect to HGS
   - Basic Hyper-V functionality works without HGS access
   - HGS is primarily used for advanced security features like Shielded VMs
   - Not typically configured on Windows 10/11 client systems

2. **Image Management Service**
   - Made optional with fallback mechanisms in the VMMS initialization
   - Implements multiple layers of error handling:
     - Explicit null checks for WMI connections
     - Try/catch blocks with panic recovery
     - Helpful debug logging
   - VHD operations can use alternative methods via VirtualSystemManagementService when ImageManagementService is unavailable
   - Differencing disks and regular VHDs can both be created through fallback methods
   - PowerShell fallback is available as a last resort when both ImageManagementService and VirtualSystemManagementService are unavailable (implemented in `provider/pkg/provider/vhdfile/vhdfileController.go`)

3. **Virtual System Management Service**
   - Critical service with enhanced error handling
   - Provides core functionality when other services are unavailable
   - Used as a fallback for disk operations

#### Implementation Details

- Each service connection is wrapped in recovery blocks to prevent panics
- Careful null checking before accessing any WMI objects
- Clear logging with different levels (INFO, WARN, ERROR) to help diagnose issues
- Alternative code paths for key operations when primary methods fail
- PowerShell fallback for various operations when WMI services are unavailable:
  - Uses the `util.FindPowerShellExe()` function to locate the appropriate PowerShell executable
  - Works with both Windows PowerShell (powershell.exe) and PowerShell Core (pwsh.exe/pwsh)
  - Provides consistent fallback behavior across all resources
  - VHD operations specifically use the `CreateVirtualHardDiskFallback` function in `vhdfileController.go`
  - Invokes native PowerShell cmdlets like `New-VHD` when WMI services are unavailable
  - Supports all VHD types: Fixed, Dynamic, and Differencing
  - Includes thorough input validation and error handling
  - Properly escapes file paths and handles various parameter combinations
- Handles common Windows 10/11 permission limitations with administrator privilege reminders

## Approved Commands

- `make`, `.\make.ps1`: All build and test targets
  - `make lint`, `make format`, `make test`, `make build`, `make provider`
  - `.\make.ps1 lint`, `.\make.ps1 build`, etc.
- `git status`, `git diff`: Repository status operations (read-only)
- `go test`, `go build`, `gofmt`, `golangci-lint`: Go tooling
- `cd`, `ls`: Navigation commands

## Important Restrictions

- **NEVER RUN GIT COMMIT**: Do not attempt to create commits. Only the user should create commits.
- **NEVER MODIFY .git DIRECTORY**: Do not attempt to directly modify the git repository.
- **NEVER CREATE PULL REQUESTS**: Do not attempt to create or submit pull requests.

## Common Issues and Solutions

- **IMPORTANT: ALWAYS USE THE MICROSOFT WMI TEST CODE AS A REFERENCE**
  - When implementing any Hyper-V functionality, ALWAYS refer to Microsoft's WMI test code:
  - <https://github.com/microsoft/wmi/blob/master/pkg/virtualization/core/service/virtualmachinemanagementservice_test.go>
  - Use the existing helper methods like `vsms.AttachVirtualHardDisk()` and `vsms.AddVirtualNetworkAdapter()`

- **IMPORTANT: ROBUST HYPER-V INTERACTIONS REQUIRE MULTIPLE FALLBACK APPROACHES**
  - The codebase uses a multi-layered approach for reliability across different Windows environments:
  
  1. **PRIMARY: High-level WMI methods** from Microsoft's WMI library:
     - First attempt uses existing helper methods like `vsms.AttachVirtualHardDisk()` and `vsms.AddVirtualNetworkAdapter()`
     - These are the most reliable when available and should be the first approach
     - Example: `_, _, err := vsms.AttachVirtualHardDisk(vm, *hd.Path, diskType)`
  
  2. **SECONDARY: Direct WMI API calls** for when high-level methods fail:
     - Falls back to formatted WMI resource path strings and direct API calls
     - Uses properly formatted resource settings arrays
     - Example: `vmmsClient.AttachVirtualHardDiskDirectApi(vm, hdPath, controllerNumber, controllerLocation)`
  
  3. **FALLBACK: PowerShell cmdlets** for when WMI services are not available:
     - Final fallback using `Add-VMHardDiskDrive`, `Add-VMNetworkAdapter`, etc.
     - Use `util.RunPowerShellCommand()` for PowerShell operations
     - Example: `addVirtualNetworkAdapterPowerShell(vm, adapterName, switchName)`
     - Enhanced error handling with pattern-matching for common PowerShell errors
     - User-friendly error messages based on specific error patterns
     - Consistent error handling across all PowerShell operations

- **AddResourceSettings common errors**: When direct WMI methods are needed:
  - For direct API calls, use this system name format:
    ```go
    systemName := fmt.Sprintf("\\\\%s\\root\\virtualization\\v2:Msvm_ComputerSystem.CreationClassName=\"Msvm_ComputerSystem\",Name=\"%s\"", host.HostName, vm.InstanceID)
    ```
  - Ensure ResourceSubType matches ResourceType (31 → "Microsoft:Hyper-V:Virtual Hard Disk", 10 → "Microsoft:Hyper-V:Synthetic Ethernet Port")
  - Always check return values from WMI method calls and handle errors properly

- **TypeScript import compatibility**: Both import styles are supported:
  - Namespaced (recommended): `hyperv.machine.Machine`, `hyperv.virtualswitch.VirtualSwitch`
  - Direct (legacy): `hyperv.Machine`, `hyperv.VirtualSwitch`
- **Error "Property does not exist on type"**: Check correct namespace in imports (machine, vhdfile, etc.). Example: use `hyperv.virtualswitch.VirtualSwitch` instead of `hyperv.VirtualSwitch`
- **Property is required**: Check if property is marked as optional in the Go code
- **Package.json issues**: Use `"main": "index"` not `"main": "index.js"`
- **Resource ordering**: Create VMs before attaching network adapters to them
- **NetworkAdapter creation patterns**: Two supported approaches:
  - Standalone pattern: `new hyperv.NetworkAdapter()` with `vmName` property to attach to a VM
  - Embedded pattern: Include in `Machine` resource's `networkAdapters` array property
  - Reference pattern: Create standalone adapter (without vmName) and reference in Machine (see simple-all-four example)
- **VhdFile differencing disk**: Requires parentPath but not sizeBytes with `diskType: "Differencing"`
- **VhdFile block size**: Always specify `blockSize: 1048576` (1MB) for better compatibility. Errors like "The parameter is incorrect" (0x80070057) when creating VHDs often indicate block size compatibility issues
- **Dynamic memory configuration**: Requires `dynamicMemory: true` plus `minimumMemory` and `maximumMemory` properties
- **Auto start/stop behavior**: Configure with `autoStartAction` and `autoStopAction`:
  - `autoStartAction`: "Nothing", "StartIfRunning", or "Start"
  - `autoStopAction`: "TurnOff", "Save", or "ShutDown"
- **Build errors with Go**: Use `GOOS=windows go build ./...` for Windows-specific code
- **Hyper-V not detected**: The provider will log warnings if Hyper-V is not available, which can happen if:
  - Running on a non-Windows operating system
  - Hyper-V is not enabled in Windows features
  - Insufficient permissions to access Hyper-V services
  - WMI infrastructure issues
- **WMI service connectivity errors**:
  - **HGS connection errors**:
    - On Windows 10/11, these can be safely ignored
    - HGS is only needed for advanced security features like Shielded VMs
    - Basic Hyper-V functionality works without HGS
  - **"Object is not connected to server" errors**:
    - Common on Windows 10/11 or with limited permissions
    - The provider includes fallback mechanisms for most operations
    - If seeing this error repeatedly, try running as administrator
  - **ImageManagementService unavailable**:
    - Provider automatically uses VirtualSystemManagementService as fallback
    - VHD creation operations should still work through alternative methods
    - If both ImageManagementService and VirtualSystemManagementService are unavailable, a PowerShell fallback (`New-VHD` cmdlet) is used as a last resort
    - Clear error messages will indicate which service is unavailable and which fallback method is being used
- **Windows 10/11 vs Windows Server differences**:
  - Client Windows editions (10/11) often have more limited WMI service access
  - The provider includes specific handling for client Windows editions
  - Error messages are tailored based on the detected OS version
  - Most operations will use fallback methods when primary methods fail
  - Running as administrator is more critical on Windows 10/11 than Server editions
  
- **Azure Edition compatibility**:
  - Azure Edition Windows Server has service limitations similar to Windows 10/11
  - When both ImageManagementService and VirtualSystemManagementService are unavailable on Azure Edition, PowerShell fallback is used automatically
  - A specific message with Azure Edition guidance is displayed in the Pulumi window
  - All VHD operations will work via PowerShell fallback even when these services are unavailable
- **Nil pointer dereference crashes**:
  - If you encounter these, ensure you're using the latest version with enhanced error handling
  - The provider now includes panic recovery and extensive null checks
  - Detailed logs will help identify which service is causing issues
  - The new code handles these issues gracefully with fallback methods

- **"unknown type" errors in AddResourceSettings**:
  - These errors occur when WMI expects specific parameter formats for method parameters
  - The WMI API for AddResourceSettings expects two specific parameters:
    1. First parameter: The VM path (from vm.GetPath(), not just VM name)
    2. Second parameter: An array of resource settings objects
  - Correct format example:

    ```go
    resourceSettings := []interface{}{
        map[string]interface{}{
            "ResourceType":       uint16(31), // 31 = Disk drive
            "Path":               "C:\\path\\to\\disk.vhdx",
            "ResourceSubType":    "Microsoft:Hyper-V:Synthetic SCSI Controller",
            "ControllerNumber":   uint32(0),
            "ControllerLocation": uint32(0),
        },
    }
    vsms.InvokeMethod("AddResourceSettings", []interface{}{vmPath, resourceSettings})
    ```

  - Similar structure is needed for AddNetworkAdapter and other WMI methods
  - Error handling for resource attachments should save the VM state for proper cleanup:
    - When hard drive or network adapter attachment fails, still return the VM state
    - This ensures `pulumi destroy` can properly clean up the VM even if a component failed to attach
    - Use warning logs to indicate that the state is being saved despite the attachment failure
  - Valid resource types include:
    - 31 = Disk drive (used for hard drive resources)
    - 10 = Network adapter

- **Memory errors when starting VMs**:
  - When starting VMs, Windows can return various memory-related errors:
    - "Not enough memory in the system to start the virtual machine example-vm"
    - "Not enough memory resources are available to complete this operation. (0x8007000E)"
    - "could not initialize memory"
  - The provider now handles these errors with specific guidance:
    - Shows the memory amount that caused the problem
    - Provides actionable suggestions to fix the issue
    - Returns clear error messages to Pulumi users
  - Resolution options typically include:
    - Reducing memory allocation in the VM
    - Closing other applications to free host memory
    - Adding more RAM to the host system
    - Using dynamic memory with lower minimum allocation

- **SCSI controller issues for hard drives**:
  - When attaching hard drives using PowerShell, common errors include:
    - "The operation could not be completed because no available locations were found on the disk controller"
    - Controller location conflicts when multiple hard drives use the same controller/location
  - The provider implements robust error handling:
    - Suggests using a different controller number or location
    - Explains the error in user-friendly terms
    - Shows the controller type, number, and location in the error message
  - Best practices for SCSI controller configuration:
    - Default controller type is "SCSI" (most flexible and performant)
    - Default controller number is 0 (first controller)
    - Controller locations should be unique (0, 1, 2, etc.) per controller
    - Each controller supports up to 64 locations (0-63)
    - When attaching multiple drives, specify unique locations
