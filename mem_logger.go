package logging

import "fmt"

// MemLogger stores an array of logging messages, up to a cap.

const MemLoggerMaxMsgs = 128

type LogMessage struct {
	Loglevel Level
	Msg      string
}

type MemLogger struct {
	msgs []LogMessage
}

func NewMemLogger() *MemLogger {
	return &MemLogger{msgs: make([]LogMessage, 0, MemLoggerMaxMsgs)}
}

func (l *MemLogger) Trace(vals ...interface{}) {
	l.appendToLogIfRoom(LogMessage{TraceLevel, fmt.Sprint(vals...)})
}
func (l *MemLogger) Debug(vals ...interface{}) {
	l.appendToLogIfRoom(LogMessage{DebugLevel, fmt.Sprint(vals...)})
}
func (l *MemLogger) Info(vals ...interface{}) {
	l.appendToLogIfRoom(LogMessage{InfoLevel, fmt.Sprint(vals...)})
}
func (l *MemLogger) Error(vals ...interface{}) {
	l.appendToLogIfRoom(LogMessage{ErrorLevel, fmt.Sprint(vals...)})
}

func (l *MemLogger) Tracef(format string, params ...interface{}) {
	l.appendToLogIfRoom(LogMessage{TraceLevel, fmt.Sprintf(format, params...)})
}
func (l *MemLogger) Debugf(format string, params ...interface{}) {
	l.appendToLogIfRoom(LogMessage{DebugLevel, fmt.Sprintf(format, params...)})
}
func (l *MemLogger) Infof(format string, params ...interface{}) {
	l.appendToLogIfRoom(LogMessage{InfoLevel, fmt.Sprintf(format, params...)})
}
func (l *MemLogger) Errorf(format string, params ...interface{}) {
	l.appendToLogIfRoom(LogMessage{ErrorLevel, fmt.Sprintf(format, params...)})
}

func (l *MemLogger) appendToLogIfRoom(msg LogMessage) {
	if len(l.msgs) < MemLoggerMaxMsgs {
		l.msgs = append(l.msgs, msg)
	}
}

func (l *MemLogger) SetLogLevel(newLevel Level) {
	// Store them all.
}
func (l *MemLogger) LogLevel() Level {
	return TraceLevel
}

func (l *MemLogger) ExtractMsgs() []LogMessage {
	retVal := l.msgs
	l.msgs = make([]LogMessage, 0, MemLoggerMaxMsgs)
	return retVal
}

var _ Logger = &MemLogger{}

func WriteLogMessageArray(logger Logger, msgs []LogMessage) {
	for _, msg := range msgs {
		switch msg.Loglevel {
		case TraceLevel:
			logger.Trace(msg.Msg)
		case DebugLevel:
			logger.Debug(msg.Msg)
		case InfoLevel:
			logger.Info(msg.Msg)
		case ErrorLevel:
			logger.Error(msg.Msg)
		}
	}
}
