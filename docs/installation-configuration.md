---
title: Hyper-V (preview) Installation & Configuration
meta_desc: Information on how to install and configure the Pulumi Hyper-V provider.
layout: package
---

## Installation

The Pulumi Hyper-V provider is available as a package in all Pulumi languages:

* JavaScript/TypeScript: [`@pulumi/hyperv`](https://www.npmjs.com/package/@pulumi/hyperv)
* Python: [`pulumi-hyperv`](https://pypi.org/project/pulumi-hyperv/)
* Go: [`github.com/pulumi/pulumi-hyperv-provider/sdk/go/hyperv`](https://pkg.go.dev/github.com/pulumi/pulumi-hyperv-provider/sdk/go/hyperv)
* .NET: [`Pulumi.Hyperv`](https://www.nuget.org/packages/Pulumi.Hyperv)
* Java: [`com.pulumi/hyperv`](https://central.sonatype.com/artifact/com.pulumi/hyperv)

## Requirements

### System Requirements

The Hyper-V provider requires:

* Windows 10/11 or Windows Server 2016 or later with Hyper-V enabled
* The provider will automatically detect if Hyper-V is available on your system during initialization
* If Hyper-V is not available or not enabled, the provider will log a warning message, but will still initialize (operations that require Hyper-V will fail)

### Enabling Hyper-V

#### Windows 10/11

**Method 1: Using Windows Features**
1. Press **Windows + R**, type `appwiz.cpl`, and press Enter
2. Click **Turn Windows features on or off** in the sidebar
3. Check the **Hyper-V** box (this will select all Hyper-V components)
4. Click **OK** and restart your computer when prompted

**Method 2: Using PowerShell (Administrator)**
```powershell
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All
```

**Method 3: Using DISM (Administrator)**
```cmd
dism /Online /Enable-Feature /All /FeatureName:Microsoft-Hyper-V
```

#### Windows Server 2016/2019/2022

**Method 1: Using Server Manager**
1. Open **Server Manager**
2. Click **Add roles and features**
3. Click **Next** until you reach the **Server Roles** page
4. Check **Hyper-V** and click **Next**
5. Complete the wizard and restart when prompted

**Method 2: Using PowerShell (Administrator)**
```powershell
Install-WindowsFeature -Name Hyper-V -IncludeManagementTools -Restart
```

#### Prerequisites

For Hyper-V to work properly, your system must have:
- A 64-bit processor with Second Level Address Translation (SLAT)
- CPU support for VM Monitor Mode Extension (VT-x on Intel)
- Minimum of 4 GB RAM
- BIOS-level virtualization support enabled

To check if your system supports Hyper-V, run in PowerShell:
```powershell
Get-ComputerInfo -Property "HyperVRequirementVirtualizationFirmwareEnabled", "HyperVRequirementVMMonitorModeExtensions"
```

### Hyper-V Detection

The Hyper-V provider includes automatic detection of Hyper-V availability on your system:

1. When the provider initializes, it performs checks to determine if:
   - You're running on a Windows operating system
   - The Hyper-V WMI namespace is accessible
   - The Hyper-V Virtual System Management Service is available

2. If any of these checks fail, the provider logs warning messages indicating:
   - The specific reason Hyper-V is not available
   - A reminder that Hyper-V operations will fail
   - Instructions that Hyper-V must be enabled on Windows for the provider to function properly

This detection helps identify configuration issues early, before attempting to create or manage Hyper-V resources.
