# Pulumi Hyper-V Provider (preview)

[![Actions Status](https://github.com/pulumi/pulumi-hyperv/workflows/master/badge.svg)](https://github.com/pulumi/pulumi-hyperv/actions)
[![Slack](http://www.pulumi.com/images/docs/badges/slack.svg)](https://slack.pulumi.com)
[![NPM version](https://badge.fury.io/js/%40pulumi%2Fhyperv.svg)](https://www.npmjs.com/package/@pulumi/hyperv)
[![Python version](https://badge.fury.io/py/pulumi-hyperv.svg)](https://pypi.org/project/pulumi-hyperv)
[![NuGet version](https://badge.fury.io/nu/pulumi.hyperv.svg)](https://badge.fury.io/nu/pulumi.hyperv)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/pulumi/pulumi-hyperv/sdk/go)](https://pkg.go.dev/github.com/pulumi/pulumi-hyperv/sdk/go)
[![License](https://img.shields.io/npm/l/%40pulumi%2Fpulumi.svg)](https://github.com/pulumi/pulumi-hyperv/blob/master/LICENSE)

# ⚠️⚠️⚠️ EXPERIMENTAL: ABSOLUTELY DOES NOT WORK BETWEEN COMMITS AND HAS NO SUPPORT ⚠️⚠️⚠️

The Pulumi Hyper-V Provider enables you to manage Microsoft Hyper-V resources like virtual machines, virtual switches,
and other virtualization components as part of your Pulumi Infrastructure as Code.

There are many scenarios where the Hyper-V provider can be useful:

* Creating and managing virtual machines on Windows Hyper-V hosts
* Setting up virtual networks and switches for VM connectivity
* Managing VM checkpoints and backups
* Configuring VM storage including virtual disks and ISO attachments
* Automating the deployment of complete virtualized environments

The Hyper-V provider is especially useful for organizations that utilize Microsoft's virtualization technology for development,
testing, or production environments. It allows you to define your Hyper-V infrastructure in code, making it reproducible,
version-controlled, and easier to manage at scale.

You can use the Hyper-V provider from a Pulumi program written in any Pulumi language: C#, Go, JavaScript/TypeScript,
Python, and YAML.
You'll need to [install and configure the Pulumi CLI](https://pulumi.com/docs/get-started/install) if you haven't already.

> **NOTE**: The Hyper-V provider is in preview. The API design may change ahead of general availability based on [user feedback](https://github.com/pulumi/pulumi-hyperv/issues).

## Features

The provider includes robust error handling and compatibility features designed to work across different Windows environments:

* **Automatic Hyper-V Detection**: Checks and logs the availability of Hyper-V during initialization
* **Windows 10/11 Compatibility**: Enhanced support for both client and server Windows editions
* **Fallback Mechanisms**: Multi-layer fallback approaches for core operations when primary WMI services are unavailable
* **Panic Recovery**: All WMI operations are wrapped in panic recovery blocks to prevent crashes
* **Detailed Logging**: Comprehensive logging with different levels to help diagnose issues
* **Multiple Service Paths**: Alternative code paths for key operations when primary methods fail
* **PowerShell Fallbacks**: Uses PowerShell commands as a last resort for operations like VHD creation

## Implementation Details

### Multi-Layer Fallback Design

The provider implements a multi-layer fallback approach to handle various Windows environments:

1. **Primary Method**: Uses WMI interfaces through the Microsoft WMI library
2. **Secondary Methods**: Falls back to alternative WMI services when primary ones are unavailable
3. **Last Resort**: Uses PowerShell commands when WMI services are completely unavailable

### VHD Creation and Management

VHD operations are particularly robust with multiple fallback mechanisms:

1. **ImageManagementService**: First attempts to use the Hyper-V Image Management Service
2. **VirtualSystemManagementService**: Falls back to VSMS if IMS is unavailable
3. **PowerShell Fallback**: Uses the `New-VHD` PowerShell cmdlet as a last resort

The provider supports all VHD types with comprehensive input validation:

* **Fixed**: Pre-allocated storage, better performance
* **Dynamic**: Expands as needed, more efficient storage usage
* **Differencing**: Child disks that store changes relative to a parent VHD

Differencing disks are fully supported with proper handling of parent paths and validation.

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

### SDK Import Styles

The provider supports two import styles in TypeScript:

#### Namespaced Imports (Recommended)

```typescript
import * as hyperv from "@pulumi/hyperv";

// Use namespaced resources
const vSwitch = new hyperv.virtualswitch.VirtualSwitch("example-switch", {...});
const vhd = new hyperv.vhdfile.VhdFile("example-disk", {...});
const vm = new hyperv.machine.Machine("example-vm", {...});
const adapter = new hyperv.networkadapter.NetworkAdapter("example-adapter", {...});
```

#### Direct Imports (Legacy)

```typescript
import * as hyperv from "@pulumi/hyperv";

// Use direct resources 
const vSwitch = new hyperv.VirtualSwitch("example-switch", {...});
const vhd = new hyperv.VhdFile("example-disk", {...});
const vm = new hyperv.Machine("example-vm", {...});
const adapter = new hyperv.NetworkAdapter("example-adapter", {...});
```

For new code, the namespaced style is recommended for better type safety and clarity.

## Requirements

### System Requirements

* Windows 10/11 or Windows Server 2016 or later with Hyper-V enabled
* The provider will automatically detect if Hyper-V is available on your system during initialization
* If Hyper-V is not available or not enabled, the provider will log a warning message,
but will still initialize (operations that require Hyper-V will fail)

### Enabling Hyper-V

#### Windows 10/11

##### Method 1: Using Windows Features

1. Press **Windows + R**, type `appwiz.cpl`, and press Enter
2. Click **Turn Windows features on or off** in the sidebar
3. Check the **Hyper-V** box (this will select all Hyper-V components)
4. Click **OK** and restart your computer when prompted

##### Method 2: Using PowerShell (Administrator)

```powershell
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All
```

##### Method 3: Using DISM (Administrator)

```cmd
dism /Online /Enable-Feature /All /FeatureName:Microsoft-Hyper-V
```

#### Windows Server 2016/2019/2022

##### Method 1: Using Server Manager

1. Open **Server Manager**
2. Click **Add roles and features**
3. Click **Next** until you reach the **Server Roles** page
4. Check **Hyper-V** and click **Next**
5. Complete the wizard and restart when prompted

##### Method 2: Using PowerShell (Administrator)

```powershell
pwsh setup-datacenter-server.ps1
```

#### Prerequisites

For Hyper-V to work properly, your system must have:

* A 64-bit processor with Second Level Address Translation (SLAT)
* CPU support for VM Monitor Mode Extension (VT-x on Intel)
* Minimum of 4 GB RAM
* BIOS-level virtualization support enabled

To check if your system supports Hyper-V, run in PowerShell:

```powershell
Get-ComputerInfo -Property "HyperVRequirementVirtualizationFirmwareEnabled", "HyperVRequirementVMMonitorModeExtensions"
```

#### Windows 10/11 Additional Configuration

Windows 10/11 client systems may require additional configuration for the Hyper-V provider to work properly.
The provider now includes enhanced compatibility features specifically for Windows 10/11:

1. **Run as Administrator**: Always run Pulumi commands with administrator privileges
(right-click command prompt/PowerShell and select "Run as administrator")

2. **Check Hyper-V Administrator Membership**:

   ```powershell
   # Check if your user is in the Hyper-V Administrators group
   Get-LocalGroupMember "Hyper-V Administrators"
   
   # If not, add yourself
   Add-LocalGroupMember -Group "Hyper-V Administrators" -Member "$env:USERDOMAIN\$env:USERNAME"
   ```

3. **Verify Required Services are Running**:

   ```powershell
   # Verify Hyper-V Virtual Machine Management service is running
   Get-Service vmms
   
   # If not running, start it
   Start-Service vmms
   ```

4. **Restart after Enabling Hyper-V**: Always restart your computer after enabling Hyper-V features

#### Service Availability and Fallback Mechanisms

The provider now includes robust handling for common Windows 10/11 service limitations:

* **Host Guardian Service (HGS)**: HGS connection is now optional - the provider logs a warning but continues
if it can't connect
* **Image Management Service**: Multiple fallbacks when this service is unavailable:
  * First tries VirtualSystemManagementService for VHD operations
  * Falls back to PowerShell commands (`New-VHD`) if needed
* **Virtual System Management Service**: Enhanced error handling with recovery mechanisms
* **Automatic Hyper-V Detection**: Detects OS version and provides tailored guidance

#### Error Handling and Recovery

The provider has been significantly hardened with:

* **Panic Recovery**: All WMI operations are wrapped in panic recovery blocks
* **Thorough Nil Checks**: Comprehensive checking for nil pointers throughout the codebase
* **Detailed Error Logs**: More descriptive error messages with troubleshooting guidance
* **Graceful Degradation**: Services continue with limited functionality when certain components are unavailable

If you see warnings about unavailable services, they are likely expected on Windows 10/11 systems and can generally be
ignored - the provider will attempt alternative methods to complete operations.

### Software Dependencies

* Go 1.24
* NodeJS 22.X.X or later
* Python 3.12 or later
* .NET Core 8.0

Please refer to [Contributing to Pulumi](https://github.com/pulumi/pulumi/blob/master/CONTRIBUTING.md) for installation
guidance.
