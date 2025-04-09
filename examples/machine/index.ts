import * as hyperv from "@pulumi/hyperv";

// Create a virtual switch for the VM to connect to
const vswitch = new hyperv.VirtualSwitch("example-switch", {
    name: "example-switch",
    switchType: "Internal",
});

// Create a VHD file for the VM
const vhd = new hyperv.VhdFile("example-vhd", {
    path: "C:\\VMs\\example-vm\\disk.vhdx",
    sizeBytes: 40 * 1024 * 1024 * 1024, // 40GB
    diskType: "Dynamic",
});

// Create a virtual machine
const vm = new hyperv.Machine("example-vm", {
    machineName: "example-vm",
    generation: 2,
    processorCount: 2,
    memorySize: 4096, // 4GB
});

// Create a network adapter and attach it to the VM and switch
const nic = new hyperv.NetworkAdapter("example-nic", {
    name: "example-nic",
    vmName: vm.machineName,
    switchName: vswitch.name,
});

// Export the VM ID and name
export const vmName = vm.machineName;
export const vmId = vm.vmId;
export const switchName = vswitch.name;