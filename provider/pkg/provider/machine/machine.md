# Hyper-V Machine Resource

## Overview

The Machine resource in the Pulumi Hyper-V provider allows you to create, manage, and delete virtual machines on a Hyper-V host. This resource interacts with the Virtual Machine Management Service (VMMS) to perform virtual machine operations.

## Features

- Create and delete Hyper-V virtual machines
- Configure VM hardware properties including:
  - Memory allocation
  - Processor count
  - VM generation (Gen 1 or Gen 2)
- Unique VM identification with automatic ID generation

## Implementation Details

### Resource Structure

The Machine resource implementation consists of multiple files:
- `machine.go` - Core resource type definition, input/output models, and annotations
- `machineController.go` - Implementation of CRUD operations
- `machineOutputs.go` - Output-specific methods

### Virtual Machine Creation

The `Create` method performs the following steps:

1. **Initialize Connection**: Establishes a connection to the Hyper-V host using WMI
2. **Configure VM Settings**:
   - Sets the virtual machine generation (defaults to Generation 2)
   - Configures memory settings (defaults to 1024 MB)
   - Sets processor count (defaults to 1 vCPU)
3. **Create VM**: Calls the Hyper-V API to create a new virtual machine with the specified settings

### Virtual Machine Read

The `Read` method retrieves the current state of a virtual machine by:
1. Connecting to the Hyper-V host
2. Getting the VM by name
3. Retrieving VM properties including:
   - VM ID
   - Memory settings
   - Processor configuration  
   - Generation

### Virtual Machine Update

The `Update` method currently provides a minimal implementation that preserves the VM's state while updating its metadata.

### Virtual Machine Delete

The `Delete` method:
1. Connects to the Hyper-V host
2. Gets the virtual machine by name
3. Starts the VM (to ensure it's in a state that can be properly deleted)
4. Gracefully stops the VM
5. Deletes the virtual machine

## Available Properties

| Property | Type | Description | Default |
|----------|------|-------------|---------|
| `machineName` | string | Name of the Virtual Machine | (required) |
| `generation` | int | Generation of the Virtual Machine (1 or 2) | 2 |
| `processorCount` | int | Number of processors to allocate | 1 |
| `memorySize` | int | Memory size in MB | 1024 |
| `triggers` | array | Values that trigger resource replacement when changed | (optional) |

## Future Extensions

The code includes scaffolding for future enhancements including:
- Network adapter configuration
- Hard drive attachments
- Key protector for secure boot
- Additional system settings

## Related Documentation

- [Microsoft Hyper-V Documentation](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/hyper-v-on-windows-server)
- [Pulumi Hyper-V Provider Documentation](https://www.pulumi.com/registry/packages/hyperv/)
