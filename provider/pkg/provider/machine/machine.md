# Hyper-V Virtual Machine Management Service (VMMS)

## Overview

The Virtual Machine Management Service (VMMS) is a core component of Hyper-V that manages virtual machine operations on a Windows Server or Windows Client system. This document provides information about the VMMS as implemented in the Pulumi Hyper-V provider.

## Features

- Virtual machine lifecycle management (create, start, stop, pause, resume, delete)
- Resource allocation and monitoring
- Snapshot management
- Virtual device configuration

## Implementation Details in Pulumi

### Virtual Machine Creation

The `Create` method in the `vmController` is responsible for creating a virtual machine. It performs the following steps:

1. **Generate a Unique ID**: A unique ID is generated for the virtual machine.
2. **Default Values**:
   - Memory size defaults to `1024 MB` if not specified.
   - Processor count defaults to `1` if not specified.
3. **VMMS Client Initialization**: A VMMS client is created to interact with the Hyper-V host.
4. **Virtual Machine Settings**:
   - The virtual machine is configured with `Hyper-V Generation 2`.
   - Memory and processor settings are applied based on the provided or default values.
5. **Virtual Machine Creation**: The virtual machine is created using the configured settings.

### Read Method

The `Read` method is a no-op in the current implementation. It does not perform any operations and always returns an empty state.

### Update Method

The `Update` method:

- Updates the virtual machine state if an `Update` command is provided.
- Falls back to the `Create` command if no `Update` command is specified.
- Does nothing if neither command is provided.

### Delete Method

The `Delete` method is a no-op unless a `Delete` command is explicitly specified.

### Resource Replacement with Triggers

The Machine resource supports the `triggers` property which forces resource replacement when values change. When any value in the `triggers` array changes between updates, the resource will be replaced (destroyed and recreated) rather than updated in-place.

## Available Properties

The Machine resource supports the following properties:

| Property | Type | Description | Default |
|----------|------|-------------|---------|
| `machineName` | string | Name of the Virtual Machine | (required) |
| `processorCount` | int | Number of processors to allocate | 1 |
| `memorySize` | int | Memory size in MB | 1024 |
| `create` | string | Command to run on create | (optional) |
| `update` | string | Command to run on update (falls back to create command if not specified) | (optional) |
| `delete` | string | Command to run on delete | (optional) |
| `triggers` | array | Values that trigger resource replacement when changed | (optional) |

## Default Behavior

- Outputs depend on all inputs by default.
- No explicit dependency wiring is required.

## Usage in Pulumi

When using the Pulumi Hyper-V provider, the VMMS is accessed indirectly through the `Machine` resource type.

## Authentication and Security

The VMMS requires appropriate permissions to manage Hyper-V objects. When using the Pulumi Hyper-V provider, ensure that:

1. The user running Pulumi commands has administrative privileges on the Hyper-V host.
2. Required firewall rules are configured if managing a remote Hyper-V host.
3. Proper credentials are provided when connecting to remote systems.

## Related Documentation

- [Microsoft Hyper-V Documentation](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/hyper-v-on-windows-server)
- [Pulumi Hyper-V Provider Documentation](https://www.pulumi.com/registry/packages/hyperv/)
