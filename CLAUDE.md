# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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

- Run linter: `.\make.ps1 lint`
- Run all tests: `.\make.ps1 test`
- Provider tests only: `.\make.ps1 test_provider`
- Run a single test: `cd provider && GOOS=windows go test -v -count=1 ./... -run TestName`
- Format Go code: `gofmt -w .`
- Run Go linter: `cd provider && GOOS=windows golangci-lint run --path-prefix=provider -c ../.golangci.yml`

## Post-Code Change Commands

Run these commands after every change to Go files:

- Format code: `gofmt -w .`
- Lint code: `cd provider && GOOS=windows golangci-lint run --path-prefix=provider -c ../.golangci.yml`
- Build provider: `cd $GOPATH/src/github.com/pulumi/pulumi-hyperv-provider/provider; GOOS=windows go build .`

For Markdown files:

- Lint Markdown: `npx markdownlint "**/*.md"`
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
  - vm.GetPath(): For getting the path to a VM (don't use GetID())
  - vmmsClient.GetVirtualizationConn().QueryInstances(): For WMI queries
  - wmiInstance.SetProperty(): For setting properties on WMI objects

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

3. **Virtual System Management Service**
   - Critical service with enhanced error handling
   - Provides core functionality when other services are unavailable
   - Used as a fallback for disk operations

#### Implementation Details

- Each service connection is wrapped in recovery blocks to prevent panics
- Careful null checking before accessing any WMI objects
- Clear logging with different levels (INFO, WARN, ERROR) to help diagnose issues
- Alternative code paths for key operations when primary methods fail
- Handles common Windows 10/11 permission limitations with administrator privilege reminders

## Approved Commands

- `make`, `.\make.ps1`: All build and test targets
- `git status`, `git diff`: Repository status operations
- `go test`, `go build`, `gofmt`, `golangci-lint`: Go tooling
- `cd`, `ls`: Navigation commands

## Common Issues and Solutions

- **TypeScript import compatibility**: Both import styles are supported:
  - Namespaced (recommended): `hyperv.machine.Machine`, `hyperv.virtualswitch.VirtualSwitch`
  - Direct (legacy): `hyperv.Machine`, `hyperv.VirtualSwitch`
- **Error "Property does not exist on type"**: Check correct namespace in imports (machine, vhdfile, etc.)
- **Property is required**: Check if property is marked as optional in the Go code
- **Package.json issues**: Use `"main": "index"` not `"main": "index.js"`
- **Resource ordering**: Create VMs before attaching network adapters to them
- **NetworkAdapter creation patterns**: Two supported approaches:
  - Standalone pattern: `new hyperv.NetworkAdapter()` with `vmName` property to attach to a VM
  - Embedded pattern: Include in `Machine` resource's `networkAdapters` array property
  - Reference pattern: Create standalone adapter (without vmName) and reference in Machine (see simple-all-four example)
- **VhdFile differencing disk**: Requires parentPath but not sizeBytes with `diskType: "Differencing"`
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
    - Clear error messages will indicate which service is unavailable
- **Windows 10/11 vs Windows Server differences**:
  - Client Windows editions (10/11) often have more limited WMI service access
  - The provider includes specific handling for client Windows editions
  - Error messages are tailored based on the detected OS version
  - Most operations will use fallback methods when primary methods fail
  - Running as administrator is more critical on Windows 10/11 than Server editions
- **Nil pointer dereference crashes**:
  - If you encounter these, ensure you're using the latest version with enhanced error handling
  - The provider now includes panic recovery and extensive null checks
  - Detailed logs will help identify which service is causing issues
  - The new code handles these issues gracefully with fallback methods
