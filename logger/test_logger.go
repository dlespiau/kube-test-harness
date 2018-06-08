package logger

import (
	"testing"
)

// TestLogger is a logger using testing.T.Log for its output.
type TestLogger struct {
	baseLogger
}

var _ Logger = &TestLogger{}

// ForTest implements Logger.
func (l *TestLogger) ForTest(t *testing.T) Logger {
	return &TestLogger{
		baseLogger: baseLogger{
			level: l.level,
			t:     t,
		},
	}
}

// Log implements Logger.
func (l *TestLogger) Log(level LogLevel, msg string) {
	if !l.shouldLog(level) {
		return
	}
	l.t.Helper()
	l.t.Log(msg)
}

// Logf implements Logger.
func (l *TestLogger) Logf(level LogLevel, f string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}
	l.t.Helper()
	l.t.Logf(f, args...)
}
