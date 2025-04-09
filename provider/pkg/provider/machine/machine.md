# Hyper-V Machine Resource

## Overview

The Machine resource in the Pulumi Hyper-V provider allows you to create, manage, and delete virtual machines on a Hyper-V host. This resource interacts with the Virtual Machine Management Service (VMMS) to perform virtual machine operations.

## Features

- Create and delete Hyper-V virtual machines
- Configure VM hardware properties including:
  - Memory allocation (static or dynamic with min/max)
  - Processor count
  - VM generation (Gen 1 or Gen 2)
  - Auto start/stop actions
- Attach hard drives with custom controller configuration
- Configure network adapters with virtual switch connections
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
   - Sets dynamic memory with min/max values if requested
   - Sets processor count (defaults to 1 vCPU)
   - Configures auto start/stop actions
3. **Create VM**: Calls the Hyper-V API to create a new virtual machine with the specified settings
4. **Attach Hard Drives**: Attaches any specified hard drives to the VM
5. **Configure Network Adapters**: Adds any specified network adapters to the VM

### Virtual Machine Read

The `Read` method retrieves the current state of a virtual machine by:
1. Connecting to the Hyper-V host
2. Getting the VM by name
3. Retrieving VM properties including:
   - VM ID
   - Memory settings (including dynamic memory configuration)
   - Processor configuration  
   - Generation
   - Auto start/stop actions

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
| `dynamicMemory` | bool | Enable dynamic memory for the VM | false |
| `minimumMemory` | int | Minimum memory in MB when using dynamic memory | - |
| `maximumMemory` | int | Maximum memory in MB when using dynamic memory | - |
| `autoStartAction` | string | Action on host start (Nothing, StartIfRunning, Start) | Nothing |
| `autoStopAction` | string | Action on host shutdown (TurnOff, Save, ShutDown) | TurnOff |
| `networkAdapters` | array | Network adapters to attach to the VM | [] |
| `hardDrives` | array | Hard drives to attach to the VM | [] |
| `triggers` | array | Values that trigger resource replacement when changed | (optional) |

### Network Adapter Properties

| Property | Type | Description | Default |
|----------|------|-------------|---------|
| `name` | string | Name of the network adapter | "Network Adapter" |
| `switchName` | string | Name of the virtual switch to connect to | (required) |

### Hard Drive Properties

| Property | Type | Description | Default |
|----------|------|-------------|---------|
| `path` | string | Path to the VHD/VHDX file | (required) |
| `controllerType` | string | Type of controller (IDE or SCSI) | SCSI |
| `controllerNumber` | int | Controller number | 0 |
| `controllerLocation` | int | Controller location | 0 |

## Usage Examples

```typescript
// Create a new VM with a network adapter and hard drive
const vm = new hyperv.Machine("example-vm", {
    machineName: "example-vm",
    generation: 2,
    processorCount: 4,
    memorySize: 4096,
    dynamicMemory: true,
    minimumMemory: 2048,
    maximumMemory: 8192,
    autoStartAction: "StartIfRunning",
    autoStopAction: "Save",
    hardDrives: [{
        path: "C:\\VMs\\example-vm\\disk.vhdx",
        controllerType: "SCSI",
        controllerNumber: 0,
        controllerLocation: 0
    }],
    networkAdapters: [{
        name: "Primary Network",
        switchName: "External Switch"
    }]
});
```

## Related Documentation

- [Microsoft Hyper-V Documentation](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/hyper-v-on-windows-server)
- [Pulumi Hyper-V Provider Documentation](https://www.pulumi.com/registry/packages/hyperv/)
