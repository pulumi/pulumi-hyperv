import * as hyperv from "@pulumi/hyperv";

// Create an internal virtual switch
const internalSwitch = new hyperv.VirtualSwitch("internal-switch", {
    name: "internal-network",
    switchType: "Internal",
});

// Create a private virtual switch
const privateSwitch = new hyperv.VirtualSwitch("private-switch", {
    name: "private-network",
    switchType: "Private",
});

// Create an external virtual switch
// Note: This requires an existing physical network adapter
const externalSwitch = new hyperv.VirtualSwitch("external-switch", {
    name: "external-network",
    switchType: "External",
    allowManagementOs: true,
    netAdapterName: "Ethernet", // Name of your physical network adapter
});

// Export the switch names
export const internalSwitchName = internalSwitch.name;
export const privateSwitchName = privateSwitch.name;
export const externalSwitchName = externalSwitch.name;