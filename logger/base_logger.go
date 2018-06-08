package logger

import (
	"testing"
)

type baseLogger struct {
	level LogLevel
	t     *testing.T
}

// SetLevel implements Logger.
func (l *baseLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *baseLogger) shouldLog(level LogLevel) bool {
	return level >= l.level
}
