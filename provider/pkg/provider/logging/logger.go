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

package logging

import (
	"context"
	"log"
	"os"
	"sync"
)

// Logger provides logging methods.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	LogAzureEditionFallback() // Logs specific message for Azure Edition missing services
}

// GetLogger returns a logger for the current context.
func GetLogger(context.Context) Logger {
	return &fileLogger{
		stdoutLogger: &stdoutLogger{},
	}
}

// stdoutLogger is a logger that logs to stdout.
type stdoutLogger struct{}

// Debugf logs a debug message.
func (l *stdoutLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

// Infof logs an info message.
func (l *stdoutLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

// Warnf logs a warning message.
func (l *stdoutLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}

// Errorf logs an error message.
func (l *stdoutLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// LogAzureEditionFallback logs a specific message for Azure Edition environments.
func (l *stdoutLogger) LogAzureEditionFallback() {
	azureMsg := "Both ImageManagementService and VirtualSystemManagementService are unavailable on Azure Edition. " +
		"Falling back to PowerShell commands. For information on enabling required services, see: " +
		"https://learn.microsoft.com/en-us/azure/virtual-machines/windows/hyper-v-container-configuration"
	log.Printf("[WARN] %s", azureMsg)
}

// fileLogger is a logger that logs to both stdout and a file.
type fileLogger struct {
	stdoutLogger *stdoutLogger
	fileLogger   *log.Logger
	once         sync.Once
}

// ensureFileLogger initializes the file logger.
func (l *fileLogger) ensureFileLogger() {
	l.once.Do(func() {
		file, err := os.OpenFile("wmi-actions.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("[ERROR] Failed to open log file: %v", err)
			return
		}
		l.fileLogger = log.New(file, "", log.LstdFlags)
	})
}

// Debugf logs a debug message to stdout and the log file.
func (l *fileLogger) Debugf(format string, args ...interface{}) {
	l.stdoutLogger.Debugf(format, args...)
	l.ensureFileLogger()
	if l.fileLogger != nil {
		l.fileLogger.Printf("[DEBUG] "+format, args...)
	}
}

// Infof logs an info message to stdout and the log file.
func (l *fileLogger) Infof(format string, args ...interface{}) {
	l.stdoutLogger.Infof(format, args...)
	l.ensureFileLogger()
	if l.fileLogger != nil {
		l.fileLogger.Printf("[INFO] "+format, args...)
	}
}

// Warnf logs a warning message to stdout and the log file.
func (l *fileLogger) Warnf(format string, args ...interface{}) {
	l.stdoutLogger.Warnf(format, args...)
	l.ensureFileLogger()
	if l.fileLogger != nil {
		l.fileLogger.Printf("[WARN] "+format, args...)
	}
}

// Errorf logs an error message to stdout and the log file.
func (l *fileLogger) Errorf(format string, args ...interface{}) {
	l.stdoutLogger.Errorf(format, args...)
	l.ensureFileLogger()
	if l.fileLogger != nil {
		l.fileLogger.Printf("[ERROR] "+format, args...)
	}
}

// LogAzureEditionFallback logs a specific message for Azure Edition environments.
func (l *fileLogger) LogAzureEditionFallback() {
	l.stdoutLogger.LogAzureEditionFallback()
	l.ensureFileLogger()
	if l.fileLogger != nil {
		azureMsg := "Both ImageManagementService and VirtualSystemManagementService are unavailable on Azure Edition. " +
			"Falling back to PowerShell commands. For information on enabling required services, see: " +
			"https://learn.microsoft.com/en-us/azure/virtual-machines/windows/hyper-v-container-configuration"
		l.fileLogger.Printf("[WARN] %s", azureMsg)
	}
}
