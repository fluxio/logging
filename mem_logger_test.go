package logging

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMemLogger(t *testing.T) {
	Convey("MemLogger", t, func() {
		memlogger := NewMemLogger()

		Convey("Records trace message and level", func() {
			memlogger.Trace("Bad News Bears")
			msgs := memlogger.ExtractMsgs()
			So(msgs[0], ShouldResemble, LogMessage{TraceLevel, "Bad News Bears"})
		})
		Convey("Records debug message and level", func() {
			memlogger.Debug("Bad News Bears")
			msgs := memlogger.ExtractMsgs()
			So(msgs[0], ShouldResemble, LogMessage{DebugLevel, "Bad News Bears"})
		})
		Convey("Records info message and level", func() {
			memlogger.Info("Bad News Bears")
			msgs := memlogger.ExtractMsgs()
			So(msgs[0], ShouldResemble, LogMessage{InfoLevel, "Bad News Bears"})
		})
		Convey("Records error message and level", func() {
			memlogger.Error("Bad News Bears")
			msgs := memlogger.ExtractMsgs()
			So(msgs[0], ShouldResemble, LogMessage{ErrorLevel, "Bad News Bears"})
		})

		Convey("Records formatted trace message and level", func() {
			memlogger.Tracef("Bad News Bears: %s", "Ruxbin")
			msgs := memlogger.ExtractMsgs()
			So(msgs[0], ShouldResemble, LogMessage{TraceLevel, "Bad News Bears: Ruxbin"})
		})
		Convey("Records formatted debug message and level", func() {
			memlogger.Debugf("Bad News Bears: %s", "Ruxbin")
			msgs := memlogger.ExtractMsgs()
			So(msgs[0], ShouldResemble, LogMessage{DebugLevel, "Bad News Bears: Ruxbin"})
		})
		Convey("Records formatted info message and level", func() {
			memlogger.Infof("Bad News Bears: %s", "Ruxbin")
			msgs := memlogger.ExtractMsgs()
			So(msgs[0], ShouldResemble, LogMessage{InfoLevel, "Bad News Bears: Ruxbin"})
		})
		Convey("Records formatted error message and level", func() {
			memlogger.Errorf("Bad News Bears: %s", "Ruxbin")
			msgs := memlogger.ExtractMsgs()
			So(msgs[0], ShouldResemble, LogMessage{ErrorLevel, "Bad News Bears: Ruxbin"})
		})
	})
}
