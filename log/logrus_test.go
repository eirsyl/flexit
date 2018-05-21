package log

import (
	"testing"
)

func testLogger(l Logger) {
	l.Debug("Debug")
	l.Info("Info")
	l.Warn("Warn")
	l.Error("Error")

	fieldLogger := l.WithField("field", "value")
	fieldLogger.Info("With field")

	fieldsLogger := fieldLogger.WithFields(&Fields{
		"another field": "value",
	})
	fieldsLogger.Info("With two fields")

	l.Info("Old logger")
}

func TestLogrusLogger(t *testing.T) {
	// Smoke test logrus logger
	var logger Logger = NewLogrusLogger(true)
	testLogger(logger)
}
