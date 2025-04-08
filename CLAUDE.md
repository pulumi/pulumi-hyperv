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

## Code Style Guidelines
- License header: All Go files include Apache 2.0 license header
- Imports: Standard library first, external dependencies second, project imports last
- Resource pattern: Separate files for main type, inputs, outputs, and controller
- Types: Use pointer types for optional values
- Error handling: Return errors with contextual information using `fmt.Errorf`
- Test coverage: Aim for comprehensive test coverage, especially for resources
- Documentation: Document resources in separate markdown files that are embedded
- Naming: CamelCase for exported types, unexportedCamelCase for unexported
- Windows compatibility: Use path utilities that work across platforms