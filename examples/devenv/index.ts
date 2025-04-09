import * as hyperv from "@pulumi/hyperv";
import * as pulumi from "@pulumi/pulumi";

// Create a virtual switch for isolated networking
const devSwitch = new hyperv.virtualswitch.VirtualSwitch("dev-switch", {
    name: "dev-network",
    switchType: "Private",
    notes: "Development network",
});

// Create multiple development VMs
const vmCount = 3;
const vms = [];

for (let i = 0; i < vmCount; i++) {
    const vmName = `dev-vm-${i + 1}`;

    // Create a differencing disk that uses our base VHD
    const vmDisk = new hyperv.vhdfile.VhdFile(`${vmName}-disk`, {
        path: `c:\\vms\\${vmName}\\${vmName}-disk.vhdx`,
        sizeBytes: 2949672960, // 2.75GB
        blockSize: 1048576,     // 1MB
        diskType: "Dynamic",
    });

    // Create the VM 
    const vm = new hyperv.machine.Machine(vmName, {
        machineName: vmName,
        generation: 2,
        processorCount: 4,
        memorySize: 4096, // 4GB
        dynamicMemory: true,
        minimumMemory: 2048, // 2GB
        maximumMemory: 8192, // 8GB
        autoStartAction: "StartIfRunning",
        autoStopAction: "Save",
        networkAdapters: [{
            name: `${vmName}-nic`,
            switchName: devSwitch.name,
        }],
        hardDrives: [{
            path: vmDisk.path,
            controllerType: "SCSI",
            controllerNumber: 0,
            controllerLocation: 0,
        }]
    });

    vms.push(vm);
}

export const switchName = devSwitch.name;
export const vmNames = vms.map(vm => vm.machineName);