package logger

import (
	"fmt"
	"gogs.xlh/tools/util/logger/logs"
	"strings"
)

// Log levels to control the logging output.
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	Logger.SetLevel(l)
}

func SetLogFuncCall(b bool) {
	Logger.EnableFuncCallDepth(b)
	Logger.SetLogFuncCallDepth(3)
}

// logger references the used application logger.
var Logger *logs.Logger

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	err := Logger.SetLogger(adaptername, config)
	if err != nil {
		return err
	}
	return nil
}

// Trace logs a message at trace level.
func Trace(v ...interface{}) {
	Logger.Trace(generateFmtStr(len(v)), v...)
}

func Tracef(format string, v ...interface{}) {
	Logger.SetLogFuncCallDepth(4)
	Trace(fmt.Sprintf(format, v...))
	Logger.SetLogFuncCallDepth(3)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	Logger.Debug(generateFmtStr(len(v)), v...)
}

func Debugf(format string, v ...interface{}) {
	Logger.SetLogFuncCallDepth(4)
	Debug(fmt.Sprintf(format, v...))
	Logger.SetLogFuncCallDepth(3)
}

// Info logs a message at info level.
func Info(v ...interface{}) {
	Logger.Info(generateFmtStr(len(v)), v...)
}

func Infof(format string, v ...interface{}) {
	Logger.SetLogFuncCallDepth(4)
	Info(fmt.Sprintf(format, v...))
	Logger.SetLogFuncCallDepth(3)
}

// Warning logs a message at warning level.
func Warn(v ...interface{}) {
	Logger.Warn(generateFmtStr(len(v)), v...)
}

func Warnf(format string, v ...interface{}) {
	Logger.SetLogFuncCallDepth(4)
	Warn(fmt.Sprintf(format, v...))
	Logger.SetLogFuncCallDepth(3)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	Logger.Error(generateFmtStr(len(v)), v...)
}

func Errorf(format string, v ...interface{}) {
	Logger.SetLogFuncCallDepth(4)
	Error(fmt.Sprintf(format, v...))
	Logger.SetLogFuncCallDepth(3)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	Logger.Critical(generateFmtStr(len(v)), v...)
}

func Criticalf(format string, v ...interface{}) {
	Logger.SetLogFuncCallDepth(4)
	Critical(fmt.Sprintf(format, v...))
	Logger.SetLogFuncCallDepth(3)
}

func Flush() {
	Logger.Flush()
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
