package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

// NewTextLogger returns a Logger that saves all log output to the specified writer
// as text.  Each line is logged with the following format:
//
//    LMMDD HH:MM:SSSZ filename.go:## (context): msg...
//
// where:
//    L is the log level (D = debug, I = Info, E = Error)
//    MMDD HH:MM:SS.SSSZ = timestamp: Month/Day/Hour/Minute/Second/Timezone (Z = UTC)
//    filename.go:## is the filename and line number where log was called from.
//    context is the user-specified context string associated with the logger
//
// The context allows adding arbitrary additional context data to the log
// entries, for example which block sourced a particular message.
//
// The log level and timestamp are fixed-width fields.  The filename, line
// number, and context are floating-width fields.
func NewTextLogger(dest io.Writer, context string, minLevel Level) Logger {
	if dest == nil {
		dest = os.Stderr
	}
	return &StdLogger{context, &TextWriter{Writer: dest}, minLevel}
}

// StdLogger is a simple implementation that writes to the specified io.Writer.
type StdLogger struct {
	// Context is a user-specified string that is included with all log
	// messages to describe application-specific state.
	Context  string
	Writer   Writer
	MinLevel Level
}

func (s *StdLogger) Pos() (file string, line int) {
	// 3 is the number of stack frames to skip:
	//   0: origin (this func)
	//   1: logf
	//   2: Debug/Info/Error
	//   3: The caller
	numframes := 3
	// In the case of the system logger, we have one more:
	//   3: The global Debug/Info/Error calls in logger.go.
	//   4: The caller
	if s == System {
		numframes = 4
	}
	_, file, line, ok := runtime.Caller(numframes)

	// Our initial guess might still land us in the logging subsystem.  This
	// can occur when StdLogger is wrapped in another logger (for example,
	// CancellableLogger).  In this case, search up a few more frames.
	for ok && strings.HasSuffix(file, "_logger.go") && numframes < 10 {
		numframes++
		_, file, line, ok = runtime.Caller(numframes)
	}

	if !ok {
		return "???", -1
	}

	return file, line
}

// stdlogf logs a formatted message at the given level.
// This function is named stdlogf, rather than the more natural logf, because go
// vet's -printfuncs flag is case-insensitive (!) and confuses this function
// with Logf if we pass go vet -printfuncs=logf:2 (i.e., to state that the 3rd
// argument to logf is the format string).
func (s *StdLogger) stdlogf(level Level, fmtstr string, vals ...interface{}) {
	if level < s.LogLevel() {
		return // skip it!
	}

	file, line := s.Pos()
	entry := Entry{
		Level:   level,
		Time:    time.Now(),
		File:    file,
		Line:    line,
		Context: s.Context,
		Fmt:     fmtstr,
		Args:    vals,
	}
	err := s.Writer.Write(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Log write failed: %v\nEntry: %#v", err, entry)
	}
}

const kNO_FORMAT = ""

func (l *StdLogger) Trace(vals ...interface{}) { l.stdlogf(TraceLevel, kNO_FORMAT, vals...) }
func (l *StdLogger) Debug(vals ...interface{}) { l.stdlogf(DebugLevel, kNO_FORMAT, vals...) }
func (l *StdLogger) Info(vals ...interface{})  { l.stdlogf(InfoLevel, kNO_FORMAT, vals...) }
func (l *StdLogger) Error(vals ...interface{}) { l.stdlogf(ErrorLevel, kNO_FORMAT, vals...) }

func (l *StdLogger) Tracef(fmt string, params ...interface{}) { l.stdlogf(TraceLevel, fmt, params...) }
func (l *StdLogger) Debugf(fmt string, params ...interface{}) { l.stdlogf(DebugLevel, fmt, params...) }
func (l *StdLogger) Infof(fmt string, params ...interface{})  { l.stdlogf(InfoLevel, fmt, params...) }
func (l *StdLogger) Errorf(fmt string, params ...interface{}) { l.stdlogf(ErrorLevel, fmt, params...) }

func (l *StdLogger) LogLevel() Level { return Level(atomic.LoadInt32((*int32)(&l.MinLevel))) }
func (l *StdLogger) SetLogLevel(newLevel Level) {
	atomic.StoreInt32((*int32)(&l.MinLevel), int32(newLevel))
}
