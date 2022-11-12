package log

// Logger 日志接口
type Logger interface {
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Panic(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}
