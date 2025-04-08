# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands
- Build everything: `make build`
- Provider only: `make provider`
- Debug build: `make provider_debug`
- Generate SDKs: `make codegen`
- Build specific SDK: `make [dotnet|go|nodejs|python|java]_sdk`

## Lint/Test Commands
- Run linter: `make lint`
- Run all tests: `make test`
- Provider tests only: `make test_provider`
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