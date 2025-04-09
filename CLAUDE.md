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

## Lint/Test Commands

- Run linter: `make lint` or `.\make.ps1 lint`
- Run all tests: `make test` or `.\make.ps1 test`
- Provider tests only: `make test_provider` or `.\make.ps1 test_provider`
- Run a single test: `cd provider && go test -v -count=1 ./... -run TestName`
- Format Go code: `gofmt -w .`
- Run Go linter: `golangci-lint run`
- Building HyperV Go code: `GOOS=windows go build ./...` (requires Windows GOOS due to Windows-specific imports)

## External Dependencies

- WMI Library: This project depends on the Microsoft WMI library for Hyper-V interaction
  - Repository: https://github.com/microsoft/wmi
  - Import path: `github.com/microsoft/wmi@v0.31.1`
  - Local location (Optional): `$GOPATH/src/github.com/microsoft/wmi`
  - Key packages used:
    - `pkg/base/host`: For WMI host connection
    - `pkg/virtualization/core/service`: For VirtualSystemManagementService
    - `pkg/virtualization/core/virtualsystem`: For VirtualMachine operations
    - `pkg/wmiinstance`: For WMI object manipulation
  - References:
    - Look at test files like `pkg/virtualization/core/service/virtualmachinemanagementservice_test.go` for usage examples

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

## Working with Hyper-V Resources

- Use the WMI library methods for interacting with Hyper-V resources
- Always check and handle return values from WMI method calls
- Results from InvokeMethod are returned as []interface{} where the first element is a map[string]interface{}
- Common methods:
  - vsms.InvokeMethod(): For invoking WMI methods on VirtualSystemManagementService
  - vm.GetPath(): For getting the path to a VM (don't use GetID())
  - vmmsClient.GetVirtualizationConn().QueryInstances(): For WMI queries
  - wmiInstance.SetProperty(): For setting properties on WMI objects