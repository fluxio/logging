package logging

// Tee logger takes two loggers and exposes a single logging interface

type TeeLogger struct {
	Logger  // To panic on unimplemented methods
	loggers []Logger
}

func NewTeeLogger(loggers ...Logger) *TeeLogger {
	return &TeeLogger{loggers: loggers}
}

func (l *TeeLogger) Trace(vals ...interface{}) {
	for _, logger := range l.loggers {
		logger.Trace(vals...)
	}
}
func (l *TeeLogger) Debug(vals ...interface{}) {
	for _, logger := range l.loggers {
		logger.Debug(vals...)
	}
}
func (l *TeeLogger) Info(vals ...interface{}) {
	for _, logger := range l.loggers {
		logger.Info(vals...)
	}
}
func (l *TeeLogger) Error(vals ...interface{}) {
	for _, logger := range l.loggers {
		logger.Error(vals...)
	}
}

func (l *TeeLogger) Tracef(fmt string, params ...interface{}) {
	for _, logger := range l.loggers {
		logger.Tracef(fmt, params...)
	}
}
func (l *TeeLogger) Debugf(fmt string, params ...interface{}) {
	for _, logger := range l.loggers {
		logger.Debugf(fmt, params...)
	}
}
func (l *TeeLogger) Infof(fmt string, params ...interface{}) {
	for _, logger := range l.loggers {
		logger.Infof(fmt, params...)
	}
}
func (l *TeeLogger) Errorf(fmt string, params ...interface{}) {
	for _, logger := range l.loggers {
		logger.Errorf(fmt, params...)
	}
}

func (l *TeeLogger) SetLogLevel(newLevel Level) {
	for _, logger := range l.loggers {
		logger.SetLogLevel(newLevel)
	}
}

var _ Logger = &TeeLogger{}
