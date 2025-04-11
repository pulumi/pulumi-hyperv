// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/microsoft/wmi/pkg/base/host"
	wmiinstance "github.com/microsoft/wmi/pkg/wmiinstance"
)

// IsHyperVAvailable checks if Hyper-V is available on the current system
func IsHyperVAvailable() (bool, string, error) {
	// First check if we're on Windows
	if runtime.GOOS != "windows" {
		return false, fmt.Sprintf("Hyper-V is only available on Windows, current OS: %s", runtime.GOOS), nil
	}

	// Create a WMI host connection
	wmiHost := host.NewWmiLocalHost()
	if wmiHost == nil {
		return false, "Failed to create WMI host connection", fmt.Errorf("host.NewWmiLocalHost returned nil")
	}
	// WmiHost doesn't need to be explicitly disposed

	// Create a session manager
	sm := wmiinstance.NewWmiSessionManager()
	defer sm.Close()
	defer sm.Dispose()

	// Try to connect to the Hyper-V virtualization namespace
	virtConn, err := sm.GetLocalSession("root\\virtualization\\v2")
	if err != nil {
		return false, fmt.Sprintf("Failed to create virtualization connection: %v", err), err
	}

	// Connect to the session - ensure proper cleanup with defers
	// even if connection fails
	defer func() {
		if virtConn != nil {
			virtConn.Close()
			virtConn.Dispose()
		}
	}()

	_, err = virtConn.Connect()
	if err != nil {
		// More specific error messages based on error type
		if strings.Contains(err.Error(), "Access is denied") {
			return false, "Access denied connecting to Hyper-V. Try running as Administrator.", err
		} else if strings.Contains(err.Error(), "Object is not connected") {
			return false, "Failed to connect to Hyper-V WMI service. Verify Hyper-V is enabled and WMI service is running.", err
		} else {
			return false, fmt.Sprintf("Failed to connect to Hyper-V virtualization namespace: %v", err), err
		}
	}

	// Check if Hyper-V service is running by querying for the Hyper-V management service
	hyperVServices, err := virtConn.QueryInstances("SELECT * FROM Msvm_VirtualSystemManagementService")
	if err != nil {
		if strings.Contains(err.Error(), "Object is not connected") {
			return false, "Lost connection to Hyper-V WMI service. Verify Hyper-V is enabled and WMI service is running.", err
		}
		return false, fmt.Sprintf("Failed to query Hyper-V management service: %v", err), err
	}

	if len(hyperVServices) == 0 {
		return false, "Hyper-V management service not found. Hyper-V might not be enabled.", nil
	}

	// Try to detect OS version to provide better guidance
	osConn, err := sm.GetLocalSession("root\\cimv2")
	if err != nil {
		return false, "Failed to create CIM connection", err
	}
	// Connect to the session - ensure proper cleanup with defers
	// even if connection fails
	defer func() {
		if osConn != nil {
			osConn.Close()
			osConn.Dispose()
		}
	}()
	_, err = osConn.Connect()
	if err != nil {
		if strings.Contains(err.Error(), "Access is denied") {
			return false, "Access denied connecting to CIM namespace. Try running as Administrator.", err
		} else if strings.Contains(err.Error(), "Object is not connected") {
			return false, "Failed to connect to CIM namespace. Verify WMI service is running.", err
		} else {
			return false, fmt.Sprintf("Failed to connect to CIM namespace: %v", err), err
		}
	}
	// Attempt to get OS version information
	osInfo, err := getOSInfo(osConn)
	if err != nil {
		// Non-critical error, just continue
		log.Printf("[DEBUG] Could not detect OS version: %v", err)
		return false, "", fmt.Errorf("[DEBUG] Could not detect OS version: %v, unsure if hyperv is available", err)
	}

	if osInfo.isServer {
		return true, fmt.Sprintf("Hyper-V is available on Windows Server %s", osInfo.version), nil
	} else {
		return true, fmt.Sprintf("Hyper-V is available on Windows %s", osInfo.version), nil
	}
}

// osInfo contains basic information about the OS
type osInfo struct {
	version  string
	isServer bool
}

// getOSInfo attempts to get OS version information using compatible methods for Windows 10/11
func getOSInfo(conn *wmiinstance.WmiSession) (*osInfo, error) {
	// Query for OS information from WMI/CIM
	result := &osInfo{
		version:  "Unknown",
		isServer: false,
	}

	// Safety check for nil connection
	if conn == nil {
		return result, fmt.Errorf("WMI session is nil")
	}

	// Wrap the query in a panic recovery to prevent crashes
	var instances []*wmiinstance.WmiInstance
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic during WMI query: %v", r)
				log.Printf("[ERROR] %v", err)
			}
		}()

		// First try the modern CIM query approach
		instances, err = conn.QueryInstances("SELECT Caption, Version FROM CIM_OperatingSystem")
	}()

	// If the first query failed or returned no results, try the fallback
	if err != nil || len(instances) == 0 {
		log.Printf("[DEBUG] CIM query failed, falling back to WMI: %v", err)

		// Wrap the fallback query in a panic recovery too
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic during fallback WMI query: %v", r)
					log.Printf("[ERROR] %v", err)
				}
			}()

			instances, err = conn.QueryInstances("SELECT Caption, Version FROM Win32_OperatingSystem")
		}()

		if err != nil {
			return result, fmt.Errorf("failed to query OS info: %w", err)
		}
	}

	if len(instances) == 0 {
		return result, fmt.Errorf("no OS info found")
	}

	// Wrap property access in panic recovery too
	var captionProp interface{}
	var versionProp interface{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic during property access: %v", r)
				log.Printf("[ERROR] %v", err)
			}
		}()

		// Safely check for instances and grab Caption
		if len(instances) > 0 && instances[0] != nil {
			captionProp, err = instances[0].GetProperty("Caption")
		}
	}()

	if err == nil && captionProp != nil {
		caption, ok := captionProp.(string)
		if ok && caption != "" {
			// Check if it's a server OS
			result.isServer = strings.Contains(strings.ToLower(caption), "server")

			// Extract version from caption (simple heuristic)
			captionLower := strings.ToLower(caption)
			if strings.Contains(caption, "10") {
				result.version = "10"
			} else if strings.Contains(caption, "11") {
				result.version = "11"
			} else if strings.Contains(caption, "2016") {
				// Check for Azure Edition in 2016
				if strings.Contains(captionLower, "azure") {
					result.version = "2016-datacenter-azure-edition"
				} else {
					result.version = "2016"
				}
			} else if strings.Contains(caption, "2019") {
				// Check for Azure Edition in 2019
				if strings.Contains(strings.ToLower(caption), "azure") {
					result.version = "2019-datacenter-azure-edition"
				} else {
					result.version = "2019"
				}
			} else if strings.Contains(caption, "2022") {
				// Check for Azure Edition in 2022
				if strings.Contains(strings.ToLower(caption), "azure") {
					result.version = "2022-datacenter-azure-edition"
				} else {
					result.version = "2022"
				}
			} else if strings.Contains(caption, "2025") {
				// Future-proofing for 2025 release
				if strings.Contains(strings.ToLower(caption), "azure") {
					result.version = "2025-datacenter-azure-edition"
				} else {
					result.version = "2025"
				}
			} else if strings.Contains(caption, "2023") {
				// Support for potential 2023 version
				if strings.Contains(strings.ToLower(caption), "azure") {
					result.version = "2023-datacenter-azure-edition"
				} else {
					result.version = "2023"
				}
			}
		}
	}

	// Get more precise version from Version field if available
	// Wrap in panic recovery
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Just log and continue - version is not critical
				log.Printf("[DEBUG] Error getting Version property: %v", r)
			}
		}()

		// Safely check for instances and grab Version
		if len(instances) > 0 && instances[0] != nil {
			versionProp, _ = instances[0].GetProperty("Version")
			if versionProp != nil {
				log.Printf("[DEBUG] Version property found: %v", versionProp)
				// Version property exists but we don't need to do additional parsing
				// The simple heuristic above using Caption is sufficient for our needs
				// If future code needs to use the version parts, it can split the version string
			}
		}
	}()

	return result, nil
}

// CheckHyperVSupport checks if Hyper-V is supported and logs the result
func CheckHyperVSupport() {
	available, message, err := IsHyperVAvailable()

	if err != nil {
		log.Printf("[WARN] Failed to detect Hyper-V support: %v", err)

		// Check for common error patterns and provide more specific guidance
		errStr := err.Error()
		if strings.Contains(errStr, "Access denied") || strings.Contains(errStr, "Access is denied") {
			log.Printf("[WARN] Access denied error detected - try running with Administrator privileges")
		} else if strings.Contains(errStr, "not found") || strings.Contains(errStr, "Object is not connected") {
			log.Printf("[WARN] Hyper-V components not found or not responding - check if Hyper-V is enabled in Windows features")
			log.Printf("[WARN] You may need to restart the WMI service or your computer")
		}

		return
	}

	if available {
		log.Printf("[INFO] %s", message)

		// Additional guidance for HGS on Windows 10
		if strings.Contains(message, "Windows 10") || strings.Contains(message, "Windows 11") {
			log.Printf("[INFO] Note: HGS (Host Guardian Service) is not typically used on Windows 10/11 client systems")
			log.Printf("[INFO] Any HGS-related warnings can be safely ignored for basic Hyper-V functionality")
		}
	} else {
		log.Printf("[WARN] %s", message)
		log.Printf("[WARN] The Pulumi Hyper-V provider requires Hyper-V to be enabled on Windows")

		// Provide guidance based on the OS
		if runtime.GOOS == "windows" {
			log.Printf("[WARN] To enable Hyper-V on Windows 10/11, run in PowerShell as Administrator:")
			log.Printf("[WARN]   Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All")
			log.Printf("[WARN] Or use Windows Features in Control Panel to enable Hyper-V")
			log.Printf("[WARN] After enabling Hyper-V, you may need to restart your computer")
			log.Printf("[WARN] To check Hyper-V status, use: Get-CimInstance -Namespace root/virtualization/v2 -ClassName Msvm_VirtualSystemManagementService")
		}

		log.Printf("[WARN] Provider operations will fail without Hyper-V support")
	}
}
