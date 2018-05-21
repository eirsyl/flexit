package log

import (
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	log *logrus.Entry
}

var _ Logger = NewLogrusLogger(true)

func NewLogrusLogger(human bool) Logger {
	logger := logrus.New()
	return &logrusLogger{
		log: logrus.NewEntry(logger),
	}
}

func (l *logrusLogger) SetLevel(level Level) {
	var logrusLevel logrus.Level
	switch level {
	case DebugLevel:
		logrusLevel = logrus.DebugLevel
	case InfoLevel:
		logrusLevel = logrus.InfoLevel
	case WarnLevel:
		logrusLevel = logrus.WarnLevel
	case ErrorLevel:
		logrusLevel = logrus.ErrorLevel
	case FatalLevel:
		logrusLevel = logrus.FatalLevel
	case PanicLevel:
		logrusLevel = logrus.PanicLevel
	default:
		logrusLevel = logrus.InfoLevel
	}
	l.log.Logger.SetLevel(logrusLevel)
}

func (l *logrusLogger) WithField(field string, value interface{}) Logger {
	return &logrusLogger{
		log: l.log.WithField(field, value),
	}
}

func (l *logrusLogger) WithFields(fields *Fields) Logger {
	var loggerFields = logrus.Fields{}
	for key, value := range *fields {
		loggerFields[key] = value
	}
	return &logrusLogger{
		log: l.log.WithFields(loggerFields),
	}
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *logrusLogger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}
