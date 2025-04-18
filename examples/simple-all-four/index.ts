import * as pulumi from "@pulumi/pulumi";
import * as hyperv from "@pulumi/hyperv";

// Create a virtual switch for the VM to connect to
const vswitch = new hyperv.virtualswitch.VirtualSwitch("example-switch", {
    name: "example-switch",
    switchType: "Internal",
});

// Create a VHD file for the VM
const vhd = new hyperv.vhdfile.VhdFile("example-vhd", {
    path: "C:\\VMs\\example-vm\\disk.vhdx",
    sizeBytes: 2 * 1024 * 1024 * 1024, // 2GB
    blockSize: 1048576, // 1MB block size
    diskType: "Dynamic",
});

// Create a network adapter and attach it to the VM and switch
const nic = new hyperv.networkadapter.NetworkAdapter("example-nic", {
    name: "example-nic",
    switchName: vswitch.name,
});

// Create a virtual machine
const vm = new hyperv.machine.Machine("example-vm", {
    machineName: "example-vm",
    generation: 2,
    processorCount: 1,
    memorySize: 2048, // 4GB
    networkAdapters: [{
        name: nic.name,
        switchName: nic.switchName
    }],
    hardDrives: [{
        path: vhd.path,
        controllerType: "SCSI",
        controllerNumber: 0,
        controllerLocation: 0,
    }]
});

// Export the VM ID and name
export const vmName = vm.machineName;
export const vmId = vm.vmId;
export const switchName = vswitch.name;