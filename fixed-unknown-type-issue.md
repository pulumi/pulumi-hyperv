# Fixed "unknown type" Issue in AddResourceSettings

## Problem
The error was occurring in the AddResourceSettings method when adding a hard drive to a VM:
```
2025/04/17 11:25:28 [ERROR] Recovered from panic in AddResourceSettings: unknown type
2025/04/17 11:25:28 [ERROR] Failed to add hard drive: recovered from panic in AddResourceSettings: unknown type
```

## Root Cause
The issue was that we were using the wrong ResourceSubType for the hard drive resource. We were using the controller type ("Microsoft:Hyper-V:Synthetic SCSI Controller") for the ResourceSubType of the disk, but we needed to use "Microsoft:Hyper-V:Virtual Hard Disk" instead.

In WMI, there is a strict requirement that the ResourceType and ResourceSubType must match:
- ResourceType 31 (disk drive) must be paired with "Microsoft:Hyper-V:Virtual Hard Disk"
- Not with a controller type like "Microsoft:Hyper-V:Synthetic SCSI Controller"

## Changes Made

1. Changed the ResourceSubType for hard drive resources:
   ```go
   resourceSettings := []interface{}{
       map[string]interface{}{
           "ResourceType":       uint16(31), // 31 = Disk drive
           "Path":               *hd.Path,
           "ResourceSubType":    "Microsoft:Hyper-V:Virtual Hard Disk", // Changed from controller type
           "ControllerNumber":   uint32(controllerNumber),
           "ControllerLocation": uint32(controllerLocation),
       },
   }
   ```

2. Removed the unused controllerSubType variable to avoid confusion, replacing it with logging.

3. Fixed a bug where vmName was used instead of id in a couple of log messages.

4. Added detailed logging to help diagnose similar issues in the future.

5. Updated CLAUDE.md with information about the proper ResourceSubType usage for AddResourceSettings.

## Testing
The provider now builds successfully and should properly add hard drives to VMs without the "unknown type" error.

## Additional Notes
This type of issue is common when working with WMI interfaces in Hyper-V, as the ResourceType and ResourceSubType combinations must be precisely correct. The same would apply for other resources like network adapters, where ResourceType and ResourceSubType must match according to the Hyper-V WMI schema.