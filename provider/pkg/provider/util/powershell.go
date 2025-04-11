package util

import (
	"fmt"
	"os/exec"
	"strings"
)

// FindPowerShellExe finds the PowerShell executable (powershell.exe or pwsh.exe)
// It tries powershell.exe first, then pwsh.exe, then pwsh (for Linux/macOS).
// Returns the executable name and nil if found, or an error if no PowerShell executable is found.
func FindPowerShellExe() (string, error) {
	// Check if PowerShell is available - try powershell.exe first, then pwsh.exe
	if _, err := exec.LookPath("powershell.exe"); err == nil {
		return "powershell.exe", nil
	} else if _, err := exec.LookPath("pwsh.exe"); err == nil {
		return "pwsh.exe", nil
	} else if _, err := exec.LookPath("pwsh"); err == nil {
		return "pwsh", nil
	}
	return "", fmt.Errorf("neither powershell.exe nor pwsh.exe found in PATH, PowerShell fallback cannot be used")
}

// RunPowerShellCommand is a helper function to run PowerShell commands with proper error handling
func RunPowerShellCommand(command string) (string, error) {
	// Find PowerShell executable
	powershellExe, err := FindPowerShellExe()
	if err != nil {
		return "", err
	}
	cmd := exec.Command(powershellExe, "-Command", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

// ParsePowerShellError attempts to parse common PowerShell error patterns and returns a more user-friendly error
func ParsePowerShellError(cmdOutput string, cmdName string, entityType string, entityName string) error {
	if cmdOutput == "" {
		return fmt.Errorf("unknown %s error: no output received from PowerShell command", entityType)
	}

	// Common PowerShell error patterns
	switch {
	case strings.Contains(cmdOutput, "ObjectNotFound"):
		return fmt.Errorf("%s not found: '%s'. Please verify it exists and you have permission to access it", entityType, entityName)

	case strings.Contains(cmdOutput, "Access is denied") || strings.Contains(cmdOutput, "AccessDenied"):
		return fmt.Errorf("access denied. Please verify you have administrator privileges")

	case strings.Contains(cmdOutput, "The parameter is incorrect"):
		return fmt.Errorf("incorrect parameter. This often happens with incompatible formats or configurations")

	case strings.Contains(cmdOutput, "Unable to find a default server with Active Directory"):
		return fmt.Errorf("unable to access Active Directory. This may happen if you're not connected to a domain controller")

	case strings.Contains(cmdOutput, "The operation failed because of a cluster validation error"):
		return fmt.Errorf("cluster validation error. This operation may require cluster administrative privileges")

	case strings.Contains(cmdOutput, "The operation failed because the process hosting the server process terminated unexpectedly"):
		return fmt.Errorf("the Hyper-V service may have restarted. Please try again")

	default:
		// If we can't identify the error, return the raw output
		return fmt.Errorf("%s operation failed: %s", entityType, cmdOutput)
	}
}
