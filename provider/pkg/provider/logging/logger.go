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
)

// Logger provides logging methods.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// GetLogger returns a logger for the current context.
func GetLogger(context.Context) Logger {
	return &stdoutLogger{}
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