package log

type nopLogger struct{}

var _ Logger = NewNopLogger()

func NewNopLogger() Logger {
	return &logrusLogger{}
}

func (l *nopLogger) SetLevel(level Level) {
}

func (l *nopLogger) WithField(field string, value interface{}) Logger {
	return &nopLogger{}
}

func (l *nopLogger) WithFields(fields *Fields) Logger {
	return &nopLogger{}
}

func (l *nopLogger) Debug(args ...interface{}) {
}

func (l *nopLogger) Debugf(format string, args ...interface{}) {
}

func (l *nopLogger) Info(args ...interface{}) {
}

func (l *nopLogger) Infof(format string, args ...interface{}) {
}

func (l *nopLogger) Warn(args ...interface{}) {
}

func (l *nopLogger) Warnf(format string, args ...interface{}) {
}

func (l *nopLogger) Error(args ...interface{}) {
}

func (l *nopLogger) Errorf(format string, args ...interface{}) {
}

func (l *nopLogger) Fatal(args ...interface{}) {
}

func (l *nopLogger) Fatalf(format string, args ...interface{}) {
}

func (l *nopLogger) Panic(args ...interface{}) {
}

func (l *nopLogger) Panicf(format string, args ...interface{}) {
}
