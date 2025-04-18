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
    // Depend on the directory creation command
    const vmDisk = new hyperv.vhdfile.VhdFile(`${vmName}-disk`, {
        path: `c:\\vms\\${vmName}\\${vmName}-disk.vhdx`,
        sizeBytes: 2949672960, // 2.75GB
        blockSize: 1048576,     // 1MB
        diskType: "Dynamic",
    });

    // Make sure the disk is created before creating the VM
    // Create the VM - explicitly wait for disk creation with apply
    const vm = new hyperv.machine.Machine(vmName, {
        machineName: vmName,
        generation: 2,
        processorCount: 1,
        dynamicMemory: true,
        minimumMemory: 1024, // 1GB
        maximumMemory: 2048, // 2GB
        autoStartAction: "StartIfRunning",
        autoStopAction: "Save",
        networkAdapters: [{
            name: `${vmName}-nic`,
            switchName: devSwitch.name,
        }],
        // Use apply to ensure we have the actual disk path and it's created before VM
        hardDrives: vmDisk.path.apply(diskPath => [{
            path: diskPath,
            controllerType: "SCSI",
            controllerNumber: 0,
            controllerLocation: 0,
        }])
    }, {
        // Set explicit dependency on disk
        dependsOn: [vmDisk, devSwitch]
    });

    vms.push(vm);
}

export const switchName = devSwitch.name;
export const vmNames = vms.map(vm => vm.machineName);