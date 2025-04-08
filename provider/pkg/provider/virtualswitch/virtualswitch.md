# Virtual Switch Resource Management

The `virtualswitch` package provides utilities for managing Hyper-V virtual switches.

## Overview

This package enables creating, modifying, and deleting virtual switches through the Pulumi Hyper-V provider. Virtual switches enable network connectivity for virtual machines.

## Key Components

### Types

- **VirtualSwitch**: Represents a Hyper-V virtual switch.

### Resource Lifecycle Methods

- **Create**: Creates a new virtual switch with specified properties.
- **Read**: Retrieves information about an existing virtual switch.
- **Update**: Modifies properties of an existing virtual switch.
- **Delete**: Removes a virtual switch.

## Available Properties

The VirtualSwitch resource supports the following properties:

| Property | Type | Description |
|----------|------|-------------|
| `name` | string | Name of the virtual switch |
| `switchType` | string | Type of switch: "External", "Internal", or "Private" |
| `allowManagementOs` | boolean | Allow the management OS to access the switch (External switches) |
| `netAdapterName` | string | Name of the physical network adapter to bind to (External switches) |

## Implementation Details

The package uses the WMI interface to interact with Hyper-V's virtual switch management functionality, providing a Go-based interface that integrates with the Pulumi resource model.

## Usage Examples

Virtual switches can be defined and managed through the Pulumi Hyper-V provider using the standard resource model.

### Creating an External Switch

```typescript
const externalSwitch = new hyperv.VirtualSwitch("external-switch", {
    name: "External Network",
    switchType: "External",
    allowManagementOs: true,
    netAdapterName: "Ethernet"
});
```

### Creating an Internal Switch

```typescript
const internalSwitch = new hyperv.VirtualSwitch("internal-switch", {
    name: "Internal Network",
    switchType: "Internal"
});
```