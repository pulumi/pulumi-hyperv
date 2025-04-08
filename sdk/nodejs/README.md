[![Actions Status](https://github.com/pulumi/pulumi-hyperv-provider/workflows/master/badge.svg)](https://github.com/pulumi/pulumi-hyperv-provider/actions)
[![Slack](http://www.pulumi.com/images/docs/badges/slack.svg)](https://slack.pulumi.com)
[![NPM version](https://badge.fury.io/js/%40pulumi%2Fhyperv.svg)](https://www.npmjs.com/package/@pulumi/hyperv)
[![Python version](https://badge.fury.io/py/pulumi-hyperv.svg)](https://pypi.org/project/pulumi-hyperv)
[![NuGet version](https://badge.fury.io/nu/pulumi.hyperv.svg)](https://badge.fury.io/nu/pulumi.hyperv)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/pulumi/pulumi-hyperv-provider/sdk/go)](https://pkg.go.dev/github.com/pulumi/pulumi-hyperv-provider/sdk/go)
[![License](https://img.shields.io/npm/l/%40pulumi%2Fpulumi.svg)](https://github.com/pulumi/pulumi-hyperv-provider/blob/master/LICENSE)

# Pulumi Hyper-V Provider (preview)

The Pulumi Hyper-V Provider enables you to manage Microsoft Hyper-V resources like virtual machines, virtual switches, and other virtualization components as part of your Pulumi Infrastructure as Code.

There are many scenarios where the Hyper-V provider can be useful:

* Creating and managing virtual machines on Windows Hyper-V hosts
* Setting up virtual networks and switches for VM connectivity
* Managing VM checkpoints and backups
* Configuring VM storage including virtual disks and ISO attachments
* Automating the deployment of complete virtualized environments

The Hyper-V provider is especially useful for organizations that utilize Microsoft's virtualization technology for development, testing, or production environments. It allows you to define your Hyper-V infrastructure in code, making it reproducible, version-controlled, and easier to manage at scale.

You can use the Hyper-V provider from a Pulumi program written in any Pulumi language: C#, Go, JavaScript/TypeScript, Python, and YAML.
You'll need to [install and configure the Pulumi CLI](https://pulumi.com/docs/get-started/install) if you haven't already.

> **NOTE**: The Hyper-V provider is in preview. The API design may change ahead of general availability based on [user feedback](https://github.com/pulumi/pulumi-hyperv-provider/issues).

## Examples

### Creating a simple virtual machine

```typescript
import * as hyperv from "@pulumi/hyperv";

// Create a new virtual machine
const vm = new hyperv.Machine("example-vm", {
    name: "example-vm",
    generation: 2,
    processorCount: 2,
    memoryStartupBytes: 2147483648, // 2GB
    networkAdapters: [{
        name: "Network Adapter",
        switchName: "Default Switch",
    }],
    hardDiskDrives: [{
        path: "c:\\vms\\example-vm\\disk.vhdx",
        controllerType: "Scsi",
        controllerNumber: 0,
        controllerLocation: 0,
    }],
});

export const vmName = vm.name;
```

### Creating a virtual switch

```typescript
import * as hyperv from "@pulumi/hyperv";

// Create a new private virtual switch
const vSwitch = new hyperv.Switch("example-switch", {
    name: "example-switch",
    switchType: "Private",
    notes: "Created by Pulumi",
});

// Create a VM connected to this switch
const vm = new hyperv.Machine("example-vm", {
    name: "example-vm",
    generation: 2,
    processorCount: 2,
    memoryStartupBytes: 2147483648, // 2GB
    networkAdapters: [{
        name: "Network Adapter",
        switchName: vSwitch.name,
    }],
    hardDiskDrives: [{
        path: "c:\\vms\\example-vm\\disk.vhdx",
        controllerType: "Scsi",
        controllerNumber: 0,
        controllerLocation: 0,
    }],
});

export const switchName = vSwitch.name;
export const vmName = vm.name;
```

### Setting up a complete virtual development environment

```typescript
import * as hyperv from "@pulumi/hyperv";
import * as pulumi from "@pulumi/pulumi";

// Create a virtual switch for isolated networking
const devSwitch = new hyperv.Switch("dev-switch", {
    name: "dev-network",
    switchType: "Private",
    notes: "Development network",
});

// Create a base VHD that we'll use for multiple VMs
const baseVhd = new hyperv.VhdFile("base-vhd", {
    path: "c:\\vms\\base\\base.vhdx",
    sizeBytes: 42949672960, // 40GB
    blockSize: 1048576,     // 1MB
    diskType: "Dynamic",
});

// Create multiple development VMs
const vmCount = 3;
const vms = [];

for (let i = 0; i < vmCount; i++) {
    const vmName = `dev-vm-${i+1}`;
    
    // Create a differencing disk that uses our base VHD
    const vmDisk = new hyperv.VhdFile(`${vmName}-disk`, {
        path: `c:\\vms\\${vmName}\\disk.vhdx`,
        parentPath: baseVhd.path,
        diskType: "Differencing",
    });
    
    // Create the VM with the differencing disk
    const vm = new hyperv.Machine(vmName, {
        name: vmName,
        generation: 2,
        processorCount: 4,
        memoryStartupBytes: 4294967296, // 4GB
        dynamicMemory: true,
        memoryMinimumBytes: 2147483648, // 2GB
        memoryMaximumBytes: 8589934592, // 8GB
        networkAdapters: [{
            name: "Network Adapter",
            switchName: devSwitch.name,
        }],
        hardDiskDrives: [{
            path: vmDisk.path,
            controllerType: "Scsi",
            controllerNumber: 0,
            controllerLocation: 0,
        }],
        // Set the VM to automatically start with the host
        autoStartAction: "StartIfRunning",
        autoStopAction: "Save",
    });
    
    vms.push(vm);
}

export const switchName = devSwitch.name;
export const vmNames = vms.map(vm => vm.name);
```

## Building

### Dependencies

- Go 1.17
- NodeJS 10.X.X or later
- Python 3.6 or later
- .NET Core 3.1

Please refer to [Contributing to Pulumi](https://github.com/pulumi/pulumi/blob/master/CONTRIBUTING.md) for installation
guidance.

### Building locally

Run the following commands to install Go modules, generate all SDKs, and build the provider:

```
$ make ensure
$ make build
$ make install
```

Add the `bin` folder to your `$PATH` or copy the `bin/pulumi-resource-hyperv` file to another location in your `$PATH`.

### Running an example

Navigate to the simple example and run Pulumi:

```
$ cd examples/simple
$ yarn link @pulumi/hyperv
$ yarn install
$ pulumi up
```

