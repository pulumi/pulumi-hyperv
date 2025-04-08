# Network Adapter Resource

The Network Adapter resource allows you to create and manage network adapters for virtual machines in Hyper-V.

## Example Usage

```typescript
import * as hyperv from "@pulumi/hyperv";

// Create a virtual switch
const vSwitch = new hyperv.VirtualSwitch("example-switch", {
    name: "example-switch",
    switchType: "Internal",
});

// Create a virtual machine
const vm = new hyperv.Machine("example-vm", {
    machineName: "example-vm",
    generation: 2,
    processorCount: 2,
    memorySize: 2048,
});

// Create a network adapter for the VM
const nic = new hyperv.NetworkAdapter("example-nic", {
    name: "example-nic",
    vmName: vm.machineName,
    switchName: vSwitch.name,
    // Optional properties
    dhcpGuard: false,
    routerGuard: false,
    vlanId: 100,
});
```

## Input Properties

| Property         | Type     | Required | Description |
|------------------|----------|----------|-------------|
| name             | string   | Yes      | Name of the network adapter |
| vmName           | string   | Yes      | Name of the virtual machine to attach the network adapter to |
| switchName       | string   | Yes      | Name of the virtual switch to connect the network adapter to |
| macAddress       | string   | No       | MAC address for the network adapter. If not specified, a dynamic MAC address will be generated |
| vlanId           | number   | No       | VLAN ID for the network adapter. If not specified, no VLAN tagging is used |
| dhcpGuard        | boolean  | No       | Enable DHCP Guard. Prevents the virtual machine from broadcasting DHCP server messages |
| routerGuard      | boolean  | No       | Enable Router Guard. Prevents the virtual machine from broadcasting router advertisement and discovery messages |
| portMirroring    | string   | No       | Port mirroring mode. Valid values are None, Source, Destination, and Both. Defaults to None |
| ieeePriorityTag  | boolean  | No       | Enable IEEE Priority Tagging. Allows the virtual machine to tag outgoing network traffic with an IEEE 802.1p priority value |
| vmqWeight        | number   | No       | VMQ weight for the network adapter. A value of 0 disables VMQ |
| ipAddresses      | string   | No       | Comma-separated list of IP addresses to assign to the network adapter |

## Output Properties

| Property         | Type     | Description |
|------------------|----------|-------------|
| adapterId        | string   | The ID of the network adapter |

## Lifecycle Management

- **Create**: Creates a new network adapter and attaches it to the specified virtual machine.
- **Read**: Reads the properties of an existing network adapter.
- **Update**: Updates the properties of an existing network adapter.
- **Delete**: Removes a network adapter from a virtual machine.

## Notes

- The network adapter creation will fail if the virtual machine or virtual switch does not exist.
- Dynamic MAC addresses are automatically generated if not specified.
- IP addresses are specified as a comma-separated string (e.g., "192.168.1.10,192.168.1.11").
- When updating a network adapter, the virtual machine may need to be powered off depending on the properties being changed.