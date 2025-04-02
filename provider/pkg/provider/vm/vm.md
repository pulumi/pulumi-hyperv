# Hyper-V Virtual Machine Management Service (VMMS)

## Overview

The Virtual Machine Management Service (VMMS) is a core component of Hyper-V that manages virtual machine operations on a Windows Server or Windows Client system. This document provides information about the VMMS as implemented in the Pulumi Hyper-V provider.

## Features

- Virtual machine lifecycle management (create, start, stop, pause, resume, delete)
- Resource allocation and monitoring
- Snapshot management
- Virtual device configuration

## Architecture

The VMMS runs as a Windows service (`vmms.exe`) and acts as the interface between the management tools and the virtualization infrastructure. The Pulumi Hyper-V provider communicates with VMMS through the Hyper-V WMI provider and Windows PowerShell cmdlets.

## Usage in Pulumi

When using the Pulumi Hyper-V provider, the VMMS is accessed indirectly through various resource types:

- `hyperv:index:VirtualMachine`
- `hyperv:index:Snapshot`
- `hyperv:index:NetworkAdapter`
- `hyperv:index:VirtualDisk`

## Authentication and Security

The VMMS requires appropriate permissions to manage Hyper-V objects. When using the Pulumi Hyper-V provider, ensure that:

1. The user running Pulumi commands has administrative privileges on the Hyper-V host
2. Required firewall rules are configured if managing a remote Hyper-V host
3. Proper credentials are provided when connecting to remote systems

## Related Documentation

- [Microsoft Hyper-V Documentation](https://docs.microsoft.com/en-us/windows-server/virtualization/hyper-v/hyper-v-on-windows-server)
- [Pulumi Hyper-V Provider Documentation](https://www.pulumi.com/registry/packages/hyperv/)