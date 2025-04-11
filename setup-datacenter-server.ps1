# Setup script for Windows Server Datacenter Azure Edition
# This script installs and configures Hyper-V components required for the pulumi-hyperv-provider

# Ensure script is running with Administrator privileges
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Error "This script must be run as Administrator. Please restart PowerShell as Administrator and try again."
    exit 1
}

Write-Host "Starting Hyper-V setup for Windows Server Datacenter Azure Edition..." -ForegroundColor Green

# Step 1: Check current OS version
$osInfo = Get-ComputerInfo | Select-Object WindowsProductName, OsVersion, OsHardwareAbstractionLayer
Write-Host "Detected Operating System:" -ForegroundColor Cyan
Write-Host "  Product Name: $($osInfo.WindowsProductName)"
Write-Host "  Version: $($osInfo.OsVersion)"
Write-Host "  HAL: $($osInfo.OsHardwareAbstractionLayer)"

# Check if this is an Azure Edition
$isAzureEdition = $osInfo.WindowsProductName -like "*Azure*"
if ($isAzureEdition) {
    Write-Host "Confirmed: Running on Azure Edition" -ForegroundColor Green
} else {
    Write-Host "Warning: This script is intended for Azure Edition, but we're running on: $($osInfo.WindowsProductName)" -ForegroundColor Yellow
}

# Step 2: Install the Hyper-V role with management tools
Write-Host "Installing Hyper-V role and management tools..." -ForegroundColor Cyan
try {
    Install-WindowsFeature -Name Hyper-V -IncludeManagementTools -ErrorAction Stop
    Write-Host "Hyper-V role installed successfully" -ForegroundColor Green
} catch {
    Write-Error "Failed to install Hyper-V role: $_"
    exit 1
}

# Step 3: Install additional Hyper-V components
Write-Host "Installing additional Hyper-V components..." -ForegroundColor Cyan
try {
    # Install specific components needed for ImageManagementService
    $hyper_v_features = @(
        "Hyper-V-PowerShell", 
        "Hyper-V-Tools", 
        "Hyper-V-Services",
        "Hyper-V"
    )
    
    Install-WindowsFeature -Name $hyper_v_features -IncludeAllSubFeature -IncludeManagementTools -ErrorAction Stop
    Write-Host "Additional Hyper-V components installed successfully" -ForegroundColor Green
} catch {
    Write-Error "Failed to install additional Hyper-V components: $_"
    # Continue anyway as these might already be included in the main role
}

# Step 3b: Register Hyper-V WMI providers
Write-Host "Registering Hyper-V WMI providers..." -ForegroundColor Cyan
try {
    # Force registration of virtualization WMI providers
    $null = & "$env:SystemRoot\System32\mofcomp.exe" "$env:SystemRoot\System32\WindowsVirtualization.V2.mof"
    $null = & "$env:SystemRoot\System32\mofcomp.exe" "$env:SystemRoot\System32\virtualization.mof"
    Write-Host "WMI providers registered successfully" -ForegroundColor Green
} catch {
    Write-Error "Failed to register WMI providers: $_"
}

# Step 4: Restart the Hyper-V Virtual Machine Management Service
Write-Host "Restarting Hyper-V Management Service (vmms)..." -ForegroundColor Cyan
try {
    Restart-Service vmms -Force -ErrorAction Stop
    Write-Host "Hyper-V Management Service restarted successfully" -ForegroundColor Green
} catch {
    Write-Error "Failed to restart Hyper-V Management Service: $_"
    # Try starting it if it wasn't running
    try {
        Start-Service vmms -ErrorAction Stop
        Write-Host "Hyper-V Management Service started successfully" -ForegroundColor Green
    } catch {
        Write-Error "Failed to start Hyper-V Management Service: $_"
        exit 1
    }
}

# Step 5: Verify the service is running
$vmmsService = Get-Service vmms
if ($vmmsService.Status -eq "Running") {
    Write-Host "Confirmed: Hyper-V Management Service is running" -ForegroundColor Green
} else {
    Write-Error "Hyper-V Management Service is not running. Status: $($vmmsService.Status)"
    exit 1
}

# Step 6: Verify VirtualSystemManagementService is available
Write-Host "Checking for VirtualSystemManagementService..." -ForegroundColor Cyan
try {
    $vsms = Get-CimInstance -Namespace root/virtualization/v2 -ClassName Msvm_VirtualSystemManagementService -ErrorAction Stop
    if ($vsms) {
        Write-Host "Confirmed: VirtualSystemManagementService is available" -ForegroundColor Green
        Write-Host "  Found service: $($vsms | Format-List | Out-String)"
    } else {
        Write-Error "VirtualSystemManagementService was not found"
        exit 1
    }
} catch {
    Write-Error "Failed to query VirtualSystemManagementService: $_"
    exit 1
}

# Step 7: Verify and attempt to fix ImageManagementService
Write-Host "Checking for ImageManagementService..." -ForegroundColor Cyan
try {
    $ims = Get-CimInstance -Namespace root/virtualization/v2 -ClassName Msvm_ImageManagementService -ErrorAction Stop
    if ($ims) {
        Write-Host "Confirmed: ImageManagementService is available" -ForegroundColor Green
        Write-Host "  Found service: $($ims | Format-List | Out-String)"
    } else {
        Write-Warning "ImageManagementService was not found initially, attempting to fix..."
        
        # Fix 1: Register the providers more explicitly
        Write-Host "Re-registering WMI providers with namespace specification..." -ForegroundColor Yellow
        try {
            # Extract WMI provider DLLs to make sure they're registered
            $null = & "$env:SystemRoot\System32\wbem\mofcomp.exe" "$env:SystemRoot\System32\wbem\virtualization\microsoft-windows-hyper-v-imagemanagement.mof" 
            $null = & "$env:SystemRoot\System32\wbem\mofcomp.exe" "$env:SystemRoot\System32\wbem\virtualization\microsoft-windows-hyper-v-wmi-provider.mof"
            
            # Restart the WMI service
            Restart-Service winmgmt -Force
            Start-Sleep -Seconds 5
            
            # Restart the Hyper-V service
            Restart-Service vmms -Force
            Start-Sleep -Seconds 5
            
            # Check again
            $ims = Get-CimInstance -Namespace root/virtualization/v2 -ClassName Msvm_ImageManagementService -ErrorAction Stop
            if ($ims) {
                Write-Host "Success! ImageManagementService is now available" -ForegroundColor Green
                Write-Host "  Found service: $($ims | Format-List | Out-String)"
            } else {
                Write-Warning "ImageManagementService still not found after fix attempt"
            }
        } catch {
            Write-Warning "Error during fix attempt: $_"
        }
    }
} catch {
    Write-Warning "Failed to query ImageManagementService: $_"
    
    # Last resort fix
    Write-Host "Attempting last-resort fix for ImageManagementService..." -ForegroundColor Yellow
    try {
        # Try to reinstall Hyper-V components
        Uninstall-WindowsFeature -Name Hyper-V -IncludeManagementTools
        Start-Sleep -Seconds 5
        Install-WindowsFeature -Name Hyper-V -IncludeAllSubFeature -IncludeManagementTools
        Start-Sleep -Seconds 10
        Restart-Service vmms -Force
        
        # Check one more time
        $ims = Get-CimInstance -Namespace root/virtualization/v2 -ClassName Msvm_ImageManagementService -ErrorAction SilentlyContinue
        if ($ims) {
            Write-Host "Success! ImageManagementService is now available after reinstall" -ForegroundColor Green
        } else {
            Write-Warning "ImageManagementService still not available after reinstall"
            Write-Host "This is not critical as the provider can use VirtualSystemManagementService instead" -ForegroundColor Yellow
        }
    } catch {
        Write-Warning "Failed during last-resort fix: $_"
        Write-Host "This is not critical as the provider can use VirtualSystemManagementService instead" -ForegroundColor Yellow
    }
}

# Step 8: Check and create Hyper-V default directories
$defaultVhdPath = "$env:SystemDrive\Users\Public\Documents\Hyper-V\Virtual Hard Disks"
if (-not (Test-Path $defaultVhdPath)) {
    Write-Host "Creating default VHD directory: $defaultVhdPath" -ForegroundColor Cyan
    New-Item -Path $defaultVhdPath -ItemType Directory -Force | Out-Null
    Write-Host "Default VHD directory created" -ForegroundColor Green
} else {
    Write-Host "Default VHD directory already exists" -ForegroundColor Green
}

# Step 9: Display additional system information
Write-Host "System Information:" -ForegroundColor Cyan
Write-Host "  Computer Name: $env:COMPUTERNAME"
Write-Host "  System Directory: $env:SystemRoot\System32"
Write-Host "  PowerShell Version: $($PSVersionTable.PSVersion)"

# Step 10: Add diagnostic info specifically for ImageManagementService
Write-Host "Running ImageManagementService diagnostics..." -ForegroundColor Cyan

# Check for the specific DLLs needed
$imageDll = Test-Path "$env:SystemRoot\System32\vmimgmanagement.dll"
Write-Host "Image Management DLL exists: $imageDll" -ForegroundColor $(if ($imageDll) { "Green" } else { "Red" })

# Check if services are running
$vmCompute = Get-Service -Name vmcompute -ErrorAction SilentlyContinue
$vmms = Get-Service -Name vmms -ErrorAction SilentlyContinue
Write-Host "VM Compute Service: $($vmCompute.Status)" -ForegroundColor $(if ($vmCompute.Status -eq 'Running') { "Green" } else { "Red" })
Write-Host "VM Management Service: $($vmms.Status)" -ForegroundColor $(if ($vmms.Status -eq 'Running') { "Green" } else { "Red" })

# Try listing all classes in the virtualization namespace
try {
    $classes = Get-CimClass -Namespace root/virtualization/v2 -ErrorAction Stop | 
               Where-Object { $_.CimClassName -like '*Image*' } |
               Select-Object -ExpandProperty CimClassName
    
    Write-Host "Found image-related classes:" -ForegroundColor Green
    $classes | ForEach-Object { Write-Host "  - $_" }
    
    # Specifically check for the image management service
    if ($classes -contains 'Msvm_ImageManagementService') {
        Write-Host "Msvm_ImageManagementService class is available!" -ForegroundColor Green
    } else {
        Write-Host "Msvm_ImageManagementService class is NOT available!" -ForegroundColor Red
    }
} catch {
    Write-Host "Failed to enumerate virtualization classes: $_" -ForegroundColor Red
}

# Final message
Write-Host "`nSetup completed successfully!" -ForegroundColor Green
Write-Host "The Hyper-V environment is now configured for use with the pulumi-hyperv-provider." -ForegroundColor Green
Write-Host "You may need to restart your system for all changes to take effect." -ForegroundColor Yellow
Write-Host "`nIf ImageManagementService is still not available after restart, the provider will"
Write-Host "automatically fall back to using VirtualSystemManagementService instead." -ForegroundColor Cyan