package logging

import (
	"context"
	"fmt"
	"reflect"
)

// DebugLogger extends the Logger with advanced debugging methods
type DebugLogger struct {
	logger Logger
}

// NewDebugLogger creates a new debug logger
func NewDebugLogger(ctx context.Context) *DebugLogger {
	return &DebugLogger{
		logger: GetLogger(ctx),
	}
}

// Debugf logs a debug message
func (l *DebugLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Infof logs an info message
func (l *DebugLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warnf logs a warning message
func (l *DebugLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Errorf logs an error message
func (l *DebugLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// LogWmiParams logs detailed information about WMI parameters
func (l *DebugLogger) LogWmiParams(methodName string, vmPath string, params interface{}) {
	l.Debugf("----------------------------------------------------------")
	l.Debugf("DETAILED DEBUG FOR %s CALL:", methodName)
	l.Debugf("  VM path: %s (type: %T)", vmPath, vmPath)
	l.Debugf("  Parameters: %+v (type: %T)", params, params)

	// If it's a slice, log each element
	if reflect.TypeOf(params).Kind() == reflect.Slice {
		s := reflect.ValueOf(params)
		l.Debugf("  Parameter slice contains %d elements", s.Len())
		for i := 0; i < s.Len(); i++ {
			elem := s.Index(i).Interface()
			l.Debugf("    Element %d: %+v (type: %T)", i, elem, elem)

			// If the element is a map, log its keys and values
			if reflect.TypeOf(elem).Kind() == reflect.Map {
				m := reflect.ValueOf(elem)
				iter := m.MapRange()
				l.Debugf("    Element %d is a map with %d entries:", i, m.Len())
				for iter.Next() {
					key := iter.Key().Interface()
					val := iter.Value().Interface()
					l.Debugf("      %v: %+v (type: %T)", key, val, val)
				}
			}
		}
	}

	// Log the combined parameter array that will be passed to InvokeMethod
	invokeParams := []interface{}{vmPath, params}
	l.Debugf("  Full InvokeMethod parameters: %+v (type: %T)", invokeParams, invokeParams)
	l.Debugf("  InvokeMethod param[0]: %+v (type: %T)", invokeParams[0], invokeParams[0])
	l.Debugf("  InvokeMethod param[1]: %+v (type: %T)", invokeParams[1], invokeParams[1])
	l.Debugf("----------------------------------------------------------")
}

// LogAddResourceSettings logs detailed information specifically for AddResourceSettings calls
func (l *DebugLogger) LogAddResourceSettings(vmPath string, diskPath string, resourceSettings interface{}) {
	l.Debugf("=================== ADD RESOURCE SETTINGS DEBUG ===================")
	l.Debugf("VM Path: %s (type: %T)", vmPath, vmPath)
	l.Debugf("Disk Path: %s", diskPath)
	l.Debugf("Resource Settings: %+v (type: %T)", resourceSettings, resourceSettings)

	// Convert to a string representation for deep inspection
	l.Debugf("String representation: %s", fmt.Sprintf("%#v", resourceSettings))

	// Recreate the parameter structure that will be passed to InvokeMethod for comparison
	recreatedParams := []interface{}{vmPath, resourceSettings}
	l.Debugf("Recreated InvokeMethod parameters: %+v", recreatedParams)
	l.Debugf("Parameter[0] type: %T, value: %v", recreatedParams[0], recreatedParams[0])
	l.Debugf("Parameter[1] type: %T, value: %v", recreatedParams[1], recreatedParams[1])
	l.Debugf("=================================================================")
}

// DumpObject dumps detailed information about an object
func (l *DebugLogger) DumpObject(prefix string, obj interface{}) {
	l.Debugf("%s: %+v", prefix, obj)

	// If it's nil, just return
	if obj == nil {
		l.Debugf("%s is nil", prefix)
		return
	}

	// Get the type and kind
	value := reflect.ValueOf(obj)
	l.Debugf("%s type: %T, kind: %s", prefix, obj, value.Kind())

	// For maps, print each key and value
	if value.Kind() == reflect.Map {
		keys := value.MapKeys()
		l.Debugf("%s contains %d keys:", prefix, len(keys))
		for _, key := range keys {
			l.Debugf("  %s[%v] = %+v", prefix, key, value.MapIndex(key))
		}
	}

	// For slices and arrays, print each element
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		l.Debugf("%s contains %d elements:", prefix, value.Len())
		for i := 0; i < value.Len(); i++ {
			l.Debugf("  %s[%d] = %+v", prefix, i, value.Index(i).Interface())

			// If the element is a map, print its content too
			element := value.Index(i).Interface()
			elementValue := reflect.ValueOf(element)
			if elementValue.Kind() == reflect.Map {
				keys := elementValue.MapKeys()
				for _, key := range keys {
					l.Debugf("    %s[%d][%v] = %+v", prefix, i, key, elementValue.MapIndex(key))
				}
			}
		}
	}
}
