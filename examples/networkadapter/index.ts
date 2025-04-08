import * as hyperv from "@pulumi/hyperv";

// Create a virtual switch for the network adapter to connect to
const vswitch = new hyperv.VirtualSwitch("example-switch", {
    name: "example-switch",
    switchType: "Internal",
});

// Create a virtual machine to attach the network adapter to
const vm = new hyperv.Machine("example-vm", {
    machineName: "example-vm",
    generation: 2,
    processorCount: 2,
    memorySize: 2048,
});

// Create a network adapter and attach it to the VM and switch
const nic = new hyperv.NetworkAdapter("example-nic", {
    name: "example-nic",
    vmName: vm.machineName,
    switchName: vswitch.name,
    // Optional properties
    dhcpGuard: false,
    routerGuard: false,
    vlanId: 100,
});

// Export the adapter ID and VM name
export const adapterName = nic.name;
export const adapterId = nic.adapterId;
export const vmName = vm.machineName;