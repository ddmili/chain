package log

import (
	"chain/config"

	"github.com/sirupsen/logrus"
)

var std Logger // 标准输出

// Trace 追踪日志
func Trace(format string, args ...interface{}) {
	std.Trace(format, args...)
}

// Debug debug 日志
func Debug(format string, args ...interface{}) {
	std.Debug(format, args...)
}

// Info 一般日志
func Info(format string, args ...interface{}) {
	std.Info(format, args...)
}

// Warn 提醒
func Warn(format string, args ...interface{}) {
	std.Warn(format, args...)
}

// Error error
func Error(format string, args ...interface{}) {
	std.Error(format, args...)
}

// Panic 异常
func Panic(format string, args ...interface{}) {
	std.Panic(format, args...)
}

// Fatal fail log
func Fatal(format string, args ...interface{}) {
	std.Fatal(format, args...)
}

// NewLogger creates a new logger
func NewLogger(cg *config.Config) Logger {
	std = NewLogrusAdapt(logrus.StandardLogger())
	return std
}
