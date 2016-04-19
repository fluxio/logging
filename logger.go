// Package logging defines an injectable logging interface.
package logging

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var os_Exit = os.Exit

func init() {
	// Set the default log library to prefix output with the name of the current
	// binary and output the file & line producing the log messages.
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("» [%s] ", binaryName()))
}

func binaryName() string { return path.Base(os.Args[0]) }

// Logger defines a standard logging interface that may be backed by different
// implementations.  For example, this can easily defer to pre-existing logging
// libraries such as seelog or the standard Go log lib.  Using this shim allows
// us to inject special log handling in different parts of the application.  For
// example, each block holds a Logger interface, but the actual log output can
// be recorded to disk indepdently or unified, individually streamed to the
// client, filtered, etc.
type Logger interface {
	Trace(vals ...interface{})
	Tracef(fmt string, args ...interface{})

	Debug(vals ...interface{})
	Debugf(fmt string, args ...interface{})

	Info(vals ...interface{})
	Infof(fmt string, args ...interface{})

	Error(vals ...interface{})
	Errorf(fmt string, args ...interface{})

	LogLevel() Level
	SetLogLevel(Level)
}

// System is a single global logger for convenience.  By default, it prints
// to stderr and uses the current binary name as the context.
var System Logger = NewTextLogger(os.Stderr, binaryName(), TraceLevel)

// Convenience accessors to the system logger.
func Trace(vals ...interface{})              { System.Trace(vals...) }
func Debug(vals ...interface{})              { System.Debug(vals...) }
func Info(vals ...interface{})               { System.Info(vals...) }
func Error(vals ...interface{})              { System.Error(vals...) }
func Tracef(fmt string, args ...interface{}) { System.Tracef(fmt, args...) }
func Debugf(fmt string, args ...interface{}) { System.Debugf(fmt, args...) }
func Infof(fmt string, args ...interface{})  { System.Infof(fmt, args...) }
func Errorf(fmt string, args ...interface{}) { System.Errorf(fmt, args...) }

// Fatal is a package-level only log function that logs to Error and then exits
// the process.
func Fatal(vals ...interface{}) {
	System.Error(vals...)
	System.Error("Failed at:\n" + stack())
	os_Exit(-1)
}

// Fatalf is a package-level only log function that logs to Errorf and then
// exits the process.
func Fatalf(fmt string, args ...interface{}) {
	System.Errorf(fmt, args...)
	System.Error("Failed at:\n" + stack())
	os_Exit(-1)
}

func FatalOnErr(err error) {
	if err == nil {
		return
	}
	Fatal(err)
}

// Writer is the interface for types that can format and write log entries.
type Writer interface {
	Write(e Entry) error
}

// Entry is a single log entry.
type Entry struct {
	Level   Level
	Time    time.Time
	File    string // empty string indicates unknown
	Line    int    // -1 indicates unknown
	Context string
	Fmt     string
	Args    []interface{}
}

// Level describes the log level.
type Level int32

const (
	TraceLevel = Level(1)
	DebugLevel = Level(2)
	InfoLevel  = Level(3)
	ErrorLevel = Level(4)
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "T"
	case DebugLevel:
		return "D"
	case InfoLevel:
		return "I"
	case ErrorLevel:
		return "E"
	}
	panic(fmt.Errorf("No such log level: %d", l))
}

// ParseLevelOrDie parses the string into a Level.  If the string is invalid, it
// panics.
func ParseLevelOrDie(levelstr string) Level {
	lvl, err := ParseLevel(levelstr)
	if err != nil {
		panic(err)
	}
	return lvl
}

// ParseLevel parses the string into a Level.  If the string is invalid, it
// returns a nice error object.
func ParseLevel(levelstr string) (Level, error) {
	switch strings.ToLower(levelstr) {
	case "t", "trace":
		return TraceLevel, nil
	case "d", "debug":
		return DebugLevel, nil
	case "i", "info":
		return InfoLevel, nil
	case "e", "error":
		return ErrorLevel, nil
	}
	return ErrorLevel, fmt.Errorf("Unknown level: %q", levelstr)
}

// Return the current stack as a string.
func stack() string {
	var stack [64]uintptr
	n := runtime.Callers(3, stack[:])
	var buf bytes.Buffer
	for _, pc := range stack[:n] {
		f := runtime.FuncForPC(pc)
		file, line := f.FileLine(pc)
		fmt.Fprintf(&buf, "  %-30s \t %s:%d\n", f.Name(), file, line)
	}
	return buf.String()
}

var trackTimeFormat string = "%s spent performing operation: "

// Logs the time since the starttime occured. This works best as a defer statement:
// e.g. defer TraceTrackTime(time.Now())
// Put at the top of a function for most clarity.
func TraceTrackTime(start time.Time) {
	pc, _, _, _ := runtime.Caller(1)
	caller := runtime.FuncForPC(pc).Name()
	TraceTrackTimef(start, "%s", caller)
}

// Logs the time since the starttime occured. This works best as a defer statement:
// e.g. defer TraceTrackTimef(time.Now(), "Loading Project %d", pid)
// Put at the top of a function for most clarity.
func TraceTrackTimef(start time.Time, fmt string, args ...interface{}) {
	elapsed := time.Since(start)
	args = append([]interface{}{elapsed}, args...)
	System.Tracef(trackTimeFormat+fmt, args...)
}

// Logs the time since the starttime occured. This works best as a defer statement:
// e.g. defer DebugTrackTime(time.Now())
// Put at the top of a function for most clarity.
func DebugTrackTime(start time.Time) {
	pc, _, _, _ := runtime.Caller(1)
	caller := runtime.FuncForPC(pc).Name()
	DebugTrackTimef(start, "%s", caller)
}

// Logs the time since the starttime occured. This works best as a defer statement:
// e.g. defer DebugTrackTime(time.Now(), "Loading Project %d", pid)
// Put at the top of a function for most clarity.
func DebugTrackTimef(start time.Time, fmt string, args ...interface{}) {
	elapsed := time.Since(start)
	args = append([]interface{}{elapsed}, args...)
	System.Debugf(trackTimeFormat+fmt, args...)
}

// Logs the time since the starttime occured. This works best as a defer statement:
// e.g. defer InfoTrackTime(time.Now())
// Put at the top of a function for most clarity.
func InfoTrackTime(start time.Time) {
	pc, _, _, _ := runtime.Caller(1)
	caller := runtime.FuncForPC(pc).Name()
	InfoTrackTimef(start, "%s", caller)
}

// Logs the time since the starttime occured. This works best as a defer statement:
// e.g. defer InfoTrackTimef(time.Now(), "Loading Project %d", pid)
// Put at the top of a function for most clarity.
func InfoTrackTimef(start time.Time, fmt string, args ...interface{}) {
	elapsed := time.Since(start)
	args = append([]interface{}{elapsed}, args...)
	System.Infof(trackTimeFormat+fmt, args...)
}
