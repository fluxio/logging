package logging

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLogFatalAndFatalf(t *testing.T) {
	var exit_code int
	var buf bytes.Buffer
	os_Exit = func(code int) { exit_code = code }
	Convey("logging", t, func() {
		exit_code = 0
		buf.Reset()
		System = NewTextLogger(&buf, "die", TraceLevel)
		Convey(".Fatal", func() {
			Fatal("xyz")
			Convey("should log the content as an error", func() {
				So(buf.String(), ShouldContainSubstring, "xyz\n")
				So(buf.String(), ShouldContainSubstring, "logger_test.go")
				So(buf.String(), ShouldStartWith, ErrorLevel.String())
			})
			Convey("should exit the process with a non-zero return code", func() {
				So(exit_code, ShouldEqual, -1)
			})
		})
		Convey(".Fatalf", func() {
			Fatalf("x:%d s:%s", 17, "qrs")
			Convey("should log the content as an error", func() {
				So(buf.String(), ShouldContainSubstring, "x:17 s:qrs\n")
				So(buf.String(), ShouldContainSubstring, "logger_test.go")
				So(buf.String(), ShouldStartWith, ErrorLevel.String())
			})
			Convey("should exit the process with a non-zero return code", func() {
				So(exit_code, ShouldEqual, -1)
			})
		})
	})
}

func TestLevel(t *testing.T) {
	Convey("Level", t, func() {
		Convey("should convert to string", func() {
			So(InfoLevel.String(), ShouldEqual, "I")
		})
		Convey("should be comparable", func() {
			So(InfoLevel, ShouldBeGreaterThan, DebugLevel)
		})
		Convey("should be created via ParseLeveLOrDie", func() {
			So(ParseLevelOrDie("InFo"), ShouldEqual, InfoLevel)
			So(ParseLevelOrDie("t"), ShouldEqual, TraceLevel)
			So(ParseLevelOrDie("ERROR"), ShouldEqual, ErrorLevel)
			So(ParseLevelOrDie("debug"), ShouldEqual, DebugLevel)
			So(func() { ParseLevelOrDie("asdf") }, ShouldPanic)
			So(func() { ParseLevelOrDie("") }, ShouldPanic)
		})
	})
}
