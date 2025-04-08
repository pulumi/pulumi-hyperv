# Contributing to the Pulumi HyperV Provider

First, thank you for your interest in contributing to the Pulumi HyperV Provider! This document provides guidelines and instructions for contributing to this repository.

## Code of Conduct

Please read our [Code of Conduct](CODE-OF-CONDUCT.md) before participating in this project.

## Development Environment Setup

### Prerequisites

- [Go 1.17 or later](https://golang.org/dl/)
- [NodeJS 10.X.X or later](https://nodejs.org/en/download/)
- [Python 3.6 or later](https://www.python.org/downloads/)
- [.NET Core 3.1 or later](https://dotnet.microsoft.com/download)
- [PowerShell 7 or later](https://github.com/PowerShell/PowerShell/releases) (required for Windows builds)
- [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)

### Building the Provider

#### Windows Development

Windows is the primary development platform for the HyperV provider, as HyperV is a Windows-only technology.

To build the provider on Windows:

```powershell
# Build everything
.\make.ps1 build

# Build only the provider binary
.\make.ps1 provider

# Build a debug version of the provider
.\make.ps1 provider_debug

# Generate SDKs only
.\make.ps1 codegen

# Build a specific language SDK
.\make.ps1 [dotnet|go|nodejs|python|java]_sdk

# Note: When building a specific SDK, the schema file will be generated
# only if needed, and the provider binary will only be built if it doesn't exist
```

#### Linux/macOS Development

While the provider implementation requires Windows, you can still develop and build the SDKs on Linux or macOS:

```bash
# Build everything
make build

# Build only the provider binary
make provider

# Generate SDKs only
make codegen

# Build a specific language SDK
make [dotnet|go|nodejs|python|java]_sdk

# Note: When building a specific SDK, the schema file will be generated
# only if needed, and the provider binary will only be built if it doesn't exist
```

### Installing Locally

To use your locally built provider:

```powershell
# Windows
.\make.ps1 install

# Linux/macOS
make install
```

## Development Workflow

### Code Style

The codebase follows these style guidelines:

- **License Headers**: All Go files must include the Apache 2.0 license header
- **Import Ordering**: Standard library first, external dependencies second, project imports last
- **Resource Pattern**: Resources use separate files for main type, inputs, outputs, and controller
- **Types**: Use pointer types for optional values
- **Error Handling**: Return errors with contextual information using `fmt.Errorf`
- **Documentation**: All resources should be documented in separate markdown files that are embedded
- **Naming**: Use CamelCase for exported types, unexportedCamelCase for unexported types
- **Windows Compatibility**: Use path utilities that work across platforms

### Testing

```powershell
# Run linter
.\make.ps1 lint

# Run all tests
.\make.ps1 test

# Run provider tests only
.\make.ps1 test_provider

# Run a specific test
cd provider && go test -v -count=1 ./... -run TestName
```

### Adding New Resources

1. Create a new folder in `provider/pkg/provider` for your resource
2. Create the main resource file, controller file, and outputs file
3. Document your resource in a separate markdown file
4. Update the provider schema in `provider.go`
5. Regenerate SDKs with `.\make.ps1 codegen`
6. Add tests for your resource
7. Add an example of your resource to the examples directory

## Submitting Pull Requests

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

## Windows-Specific Development Notes

Since this provider interfaces with HyperV on Windows, most testing and implementation requires a Windows environment.

For local development:
- PowerShell 7+ is required to run the build scripts
- Use `.\make.ps1 -ForceWindowsMode $true [target]` to force Windows mode when needed
- HyperV must be enabled on your host for integration testing

## Additional Resources

- [Pulumi Provider Documentation](https://www.pulumi.com/docs/guides/pulumi-packages/)
- [Pulumi SDK Reference](https://pkg.go.dev/github.com/pulumi/pulumi-hyperv-provider/sdk/go)