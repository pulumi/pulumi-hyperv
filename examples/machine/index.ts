import * as hyperv from "@pulumi/hyperv";

// Create a virtual machine
const vm = new hyperv.Machine("example-vm", {
    machineName: "example-vm",
    generation: 2,
    processorCount: 2,
    memorySize: 4096, // 4GB
});

// Export the VM ID and name
export const vmName = vm.machineName;
export const vmId = vm.vmId;
