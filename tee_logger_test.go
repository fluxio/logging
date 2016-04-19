package logging

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTeeLogger(t *testing.T) {
	Convey("TeeLogger", t, func() {
		buf1, buf2 := bytes.Buffer{}, bytes.Buffer{}
		logger1 := NewTextLogger(&buf1, "Logger", TraceLevel)
		logger2 := NewTextLogger(&buf2, "Logger", TraceLevel)
		teelogger := NewTeeLogger(logger1, logger2)

		compareBuffers := func(lev Level, str string) {
			strbuf1 := buf1.String()
			strbuf2 := buf2.String()

			So(strbuf1, ShouldContainSubstring, lev.String())
			So(strbuf2, ShouldContainSubstring, lev.String())
			So(strbuf1, ShouldContainSubstring, str)
			So(strbuf2, ShouldContainSubstring, str)
		}

		Convey("Writes traces to both", func() {
			teelogger.Trace("Bad News Bears")
			compareBuffers(TraceLevel, "Bad News Bears")

			teelogger.Tracef("Bad News Bears: %s", "Ruxbin")
			compareBuffers(TraceLevel, "Bad News Bears: Ruxbin")
		})
		Convey("Writes debugs to both", func() {
			teelogger.Debug("Bad News Bears")
			compareBuffers(DebugLevel, "Bad News Bears")
			teelogger.Debugf("Bad News Bears: %s", "Ruxbin")
			compareBuffers(DebugLevel, "Bad News Bears: Ruxbin")
		})
		Convey("Writes infos to both", func() {
			teelogger.Info("Bad News Bears")
			compareBuffers(InfoLevel, "Bad News Bears")
			teelogger.Infof("Bad News Bears: %s", "Ruxbin")
			compareBuffers(InfoLevel, "Bad News Bears: Ruxbin")
		})
		Convey("Writes errors to both", func() {
			teelogger.Error("Bad News Bears")
			compareBuffers(ErrorLevel, "Bad News Bears")
			teelogger.Errorf("Bad News Bears: %s", "Ruxbin")
			compareBuffers(ErrorLevel, "Bad News Bears: Ruxbin")
		})
	})
}
