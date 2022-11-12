package log

import (
	"github.com/sirupsen/logrus"
)

type logrusAdapt struct {
	l *logrus.Logger
}

func (s logrusAdapt) Trace(format string, args ...interface{}) {
	s.l.Tracef(format, args...)
}

func (s logrusAdapt) Debug(format string, args ...interface{}) {
	s.l.Debugf(format, args...)
}

func (s logrusAdapt) Info(format string, args ...interface{}) {
	s.l.Infof(format, args...)
}

func (s logrusAdapt) Warn(format string, args ...interface{}) {
	s.l.Warnf(format, args...)
}

func (s logrusAdapt) Error(format string, args ...interface{}) {
	s.l.Errorf(format, args...)
}

func (s logrusAdapt) Panic(format string, args ...interface{}) {
	s.l.Panicf(format, args...)
}

func (s logrusAdapt) Fatal(format string, args ...interface{}) {
	s.l.Fatalf(format, args...)
}

// NewLogrusAdapt create logrus adapt
func NewLogrusAdapt(l *logrus.Logger) Logger {
	return &logrusAdapt{
		l: l,
	}
}
