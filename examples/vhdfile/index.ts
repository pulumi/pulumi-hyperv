import * as hyperv from "@pulumi/hyperv";

// Create a fixed-size VHD file
const fixedVhd = new hyperv.VhdFile("fixed-vhd", {
    path: "C:\\VMs\\fixed\\disk.vhdx",
    sizeBytes: 20 * 1024 * 1024 * 1024, // 20GB
    diskType: "Fixed",
});

// Create a dynamically expanding VHD file
const dynamicVhd = new hyperv.VhdFile("dynamic-vhd", {
    path: "C:\\VMs\\dynamic\\disk.vhdx",
    sizeBytes: 40 * 1024 * 1024 * 1024, // 40GB
    diskType: "Dynamic",
    blockSize: 1 * 1024 * 1024, // 1MB block size (optional)
});

// Create another VHD file
const largeVhd = new hyperv.VhdFile("large-vhd", {
    path: "C:\\VMs\\large\\large.vhdx",
    sizeBytes: 100 * 1024 * 1024 * 1024, // 100GB
    diskType: "Dynamic",
});

// Export the VHD paths
export const fixedVhdPath = fixedVhd.path;
export const dynamicVhdPath = dynamicVhd.path;
export const largeVhdPath = largeVhd.path;