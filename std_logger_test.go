package logging

import (
	"bytes"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func args(v ...interface{}) []interface{} { return v }

type captureWriter Entry

func (c *captureWriter) Write(e Entry) error { *c = captureWriter(e); return nil }

func TestStdLogger(t *testing.T) {
	Convey("Standard Logger", t, func() {
		var c captureWriter
		log := StdLogger{"test", &c, TraceLevel}
		Convey("should capture the logging time", func() {
			t0 := time.Now()
			log.Info("")
			t1 := time.Now()
			So(c.Time, ShouldHappenOnOrBetween, t0, t1)
		})
		Convey("should capture the logging args", func() {
			log.Debugf("asdf", 1, 2, "x", log)
			So(c.Fmt, ShouldEqual, "asdf")
			So(c.Args, ShouldResemble, []interface{}{1, 2, "x", log})
		})
		Convey("should capture the context", func() {
			log.Info("")
			So(c.Context, ShouldEqual, "test")
		})
		Convey("should capture the log level", func() {
			log.Debug("")
			So(c.Level, ShouldEqual, DebugLevel)
			log.Info("")
			So(c.Level, ShouldEqual, InfoLevel)
			log.Error("")
			So(c.Level, ShouldEqual, ErrorLevel)
			log.Trace("")
			So(c.Level, ShouldEqual, TraceLevel)
		})
		Convey("should accurately log the call point", func() {
			log.Info("hi")
			So(c.File, ShouldContainSubstring, "std_logger_test.go")
		})
		Convey("should filter out log levels", func() {
			log.Tracef("a")
			So(c.Fmt, ShouldEqual, "a")

			log.MinLevel = InfoLevel

			log.Tracef("b")
			So(c.Fmt, ShouldEqual, "a") // unchanged
			log.Debugf("c")
			So(c.Fmt, ShouldEqual, "a") // unchanged
			log.Infof("d")
			So(c.Fmt, ShouldEqual, "d") // picked it up!
			log.Errorf("e")
			So(c.Fmt, ShouldEqual, "e") // picked it up!
		})
		Convey("should allow change of log level", func() {
			log.Tracef("a")
			So(c.Fmt, ShouldEqual, "a")

			log.SetLogLevel(InfoLevel)

			log.Tracef("b")
			So(c.Fmt, ShouldEqual, "a") // unchanged
			log.Debugf("c")
			So(c.Fmt, ShouldEqual, "a") // unchanged
			log.Infof("d")
			So(c.Fmt, ShouldEqual, "d") // picked it up!
			log.Errorf("e")
			So(c.Fmt, ShouldEqual, "e") // picked it up!

			So(log.LogLevel(), ShouldEqual, InfoLevel)
		})

	})
}

func TestNewTextLogger(t *testing.T) {
	Convey("NewTextLogger", t, func() {
		Convey("should make Loggers", func() {
			var _ Logger = NewTextLogger(nil, "ctx", TraceLevel)
		})
		Convey("write to the provided io.Writer", func() {
			var buf bytes.Buffer
			var log Logger = NewTextLogger(&buf, "Flow=f1", TraceLevel)
			log.Infof("Creating %d blocks", 20)
			So(buf.String(), ShouldContainSubstring, "Flow=f1")
			So(buf.String(), ShouldContainSubstring, "Creating 20 blocks")
		})
		Convey("and set the min level", func() {
			var buf bytes.Buffer
			log := NewTextLogger(&buf, "ctx", DebugLevel)
			log.Trace("xxx")
			log.Debug("yyy")
			log.Info("zzz")
			So(buf.String(), ShouldNotContainSubstring, "xxx")
			So(buf.String(), ShouldContainSubstring, "yyy")
			So(buf.String(), ShouldContainSubstring, "zzz")
		})
	})
}

func TestSystemLogger(t *testing.T) {
	var c captureWriter
	var saved Logger
	System, saved = &StdLogger{"fake system", &c, TraceLevel}, System
	defer func() { System = saved }()

	Convey("The global system logger", t, func() {
		Convey("should accurately log the call point", func() {
			Info("hi there")
			So(c.File, ShouldContainSubstring, "std_logger_test.go")
		})
	})
}

func TestWrappedLoggerTracing(t *testing.T) {
	Convey("Wrapped loggers should report a useful stack frame", t, func() {
		buf := &bytes.Buffer{}
		wrapped := CancellableLogger{Logger: NewTextLogger(buf, "dummy", TraceLevel)}
		wrapped.Trace("hello")
		s := buf.String()
		So(s, ShouldNotContainSubstring, "cancellable_logger.go")
		So(s, ShouldContainSubstring, "std_logger_test.go")
	})
}
