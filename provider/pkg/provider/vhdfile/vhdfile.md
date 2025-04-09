# VHD File Resource Management

The `vhdfile` package provides utilities for managing VHD (Virtual Hard Disk) files for Hyper-V virtual machines.

## Overview

This package enables creating, modifying, and deleting VHD and VHDX files through the Pulumi Hyper-V provider. It provides a clean abstraction for working with virtual disk files independent of virtual machines.

## Key Components

### Types

- **VhdFile**: Represents a VHD or VHDX file for use with Hyper-V virtual machines.

### Resource Lifecycle Methods

- **Create**: Creates a new VHD/VHDX file with specified properties.
- **Read**: Retrieves information about an existing VHD/VHDX file.
- **Update**: Modifies properties of an existing VHD/VHDX file (currently a no-op in the implementation).
- **Delete**: Removes a VHD/VHDX file.

## Available Properties

The VhdFile resource supports the following properties:

| Property | Type | Description |
|----------|------|-------------|
| `path` | string | Path where the VHD file should be created |
| `parentPath` | string | Path to parent VHD when creating differencing disks |
| `diskType` | string | Type of disk (Fixed, Dynamic, Differencing) |
| `sizeBytes` | number | Size of the disk in bytes (for Fixed and Dynamic disks) |

## Implementation Details

The package uses PowerShell commands under the hood to interact with Hyper-V's VHD management functionality, providing a Go-based interface that integrates with the Pulumi resource model.

### Update Behavior

The current implementation of the `Update` method is a no-op. Any changes to VHD properties that require modification of the underlying file structure will typically require replacing the resource rather than updating it in place.

## Usage Examples

VHD files can be defined and managed through the Pulumi Hyper-V provider using the standard resource model. These virtual disks can then be attached to virtual machines or managed independently.

### Creating a Base VHD

```typescript
const baseVhd = new hyperv.VhdFile("base-vhd", {
    path: "c:\\vms\\base\\disk.vhdx",
    sizeBytes: 40 * 1024 * 1024 * 1024, // 40GB
    diskType: "Dynamic"
});
```

### Creating a Differencing Disk

```typescript
const baseVhd = new hyperv.VhdFile("base-vhd", {
    path: "c:\\vms\\base\\disk.vhdx",
    sizeBytes: 40 * 1024 * 1024 * 1024, // 40GB
    diskType: "Dynamic"
});

const diffVhd = new hyperv.VhdFile("diff-vhd", {
    path: "c:\\vms\\vm1\\disk.vhdx",
    parentPath: baseVhd.path,
    diskType: "Differencing"
});
```

### Using with Machine Resource

The VhdFile resource can be used in conjunction with the Machine resource by attaching the VHD files to a virtual machine using the `hardDrives` array:

```typescript
// Create a base VHD
const baseVhd = new hyperv.VhdFile("base-vhd", {
    path: "c:\\vms\\base\\disk.vhdx",
    sizeBytes: 40 * 1024 * 1024 * 1024, // 40GB
    diskType: "Dynamic"
});

// Create a differencing disk based on the base VHD
const vmDisk = new hyperv.VhdFile("vm-disk", {
    path: "c:\\vms\\vm1\\disk.vhdx",
    parentPath: baseVhd.path,
    diskType: "Differencing"
});

// Create a VM and attach the differencing disk
const vm = new hyperv.Machine("example-vm", {
    machineName: "example-vm",
    generation: 2,
    processorCount: 2,
    memorySize: 2048,
    hardDrives: [{
        path: vmDisk.path,
        controllerType: "SCSI",
        controllerNumber: 0,
        controllerLocation: 0
    }]
});
```
