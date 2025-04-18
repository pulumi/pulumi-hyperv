// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/pulumi/pulumi-hyperv/provider/pkg/provider/util/testutil"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

const PULUMI_COMMAND_STDOUT = "PULUMI_COMMAND_STDOUT"
const PULUMI_COMMAND_STDERR = "PULUMI_COMMAND_STDERR"

func LogOutput(ctx context.Context, r io.Reader, doneCh chan<- struct{}, severity diag.Severity) {
	defer close(doneCh)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		msg := scanner.Text()
		l := p.GetLogger(ctx)
		switch severity {
		case diag.Info:
			l.InfoStatus(msg)
		case diag.Warning:
			l.WarningStatus(msg)
		case diag.Error:
			l.ErrorStatus(msg)
		default:
			l.DebugStatus(msg)
		}

		if testCtx, ok := ctx.(*testutil.TestContext); ok {
			testCtx.Log(severity, msg)
		}
	}
}

// NoopLogger satisfies the expected logger shape but doesn't actually log.
// It reads from the provided reader until EOF, discarding the output, then closes the channel.
func NoopLogger(r io.Reader, done chan struct{}) {
	defer close(done)
	_, _ = io.Copy(io.Discard, r)
}

type ConcurrentWriter struct {
	Writer io.Writer
	mu     sync.Mutex
}

func (w *ConcurrentWriter) Write(bs []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Writer.Write(bs)
}

// GetOSVersion returns the Windows version string, with Azure edition detection
// Returns a string like "Windows 10", "Windows 11", "Windows Server 2022", etc.
// For Azure editions, returns the full string like "Windows Server 2022-datacenter-azure-edition"
func GetOSVersion() (string, error) {
	// Find PowerShell executable first
	powershellExe, err := FindPowerShellExe()
	if err != nil {
		return "", fmt.Errorf("failed to find PowerShell executable: %v", err)
	}

	// Run PowerShell command to get OS version info
	cmd := exec.Command(powershellExe, "-Command",
		"(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion' -Name ProductName).ProductName")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get OS version: %v, output: %s", err, string(output))
	}

	version := strings.TrimSpace(string(output))
	log.Printf("[INFO] Detected OS version: %s", version)

	return version, nil
}
