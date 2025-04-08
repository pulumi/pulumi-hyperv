// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../utilities";

// Export members:
export { VirtualSwitchArgs } from "./virtualSwitch";
export type VirtualSwitch = import("./virtualSwitch").VirtualSwitch;
export const VirtualSwitch: typeof import("./virtualSwitch").VirtualSwitch = null as any;
utilities.lazyLoad(exports, ["VirtualSwitch"], () => require("./virtualSwitch"));


const _module = {
    version: utilities.getVersion(),
    construct: (name: string, type: string, urn: string): pulumi.Resource => {
        switch (type) {
            case "hyperv:virtualswitch:VirtualSwitch":
                return new VirtualSwitch(name, <any>undefined, { urn })
            default:
                throw new Error(`unknown resource type ${type}`);
        }
    },
};
pulumi.runtime.registerResourceModule("hyperv", "virtualswitch", _module)
