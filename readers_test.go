package logging

import (
	"io"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLogReader(t *testing.T) {
	logContent := `
I0101 00:00:00.000Z std_logger_test.go:00 (Flow=f1): Creating 20 blocks
I0101 00:00:00.000Z std_logger_test.go:00 (Flow=f1): Here's some JSON:
    {
      "A": "some val",
      "B": "other val"
    }
    Err:<nil>
I0101 00:00:00.000Z std_logger_test.go:00 (Block=addFoo): Computing 5 things
`[1:]

	Convey("LogReader should read each log message.", t, func() {
		r := NewLogReader(strings.NewReader(logContent))
		entry, err := r.Next()
		So(err, ShouldBeNil)
		So(entry.header, ShouldEqual, "I0101 00:00:00.000Z std_logger_test.go:00 (Flow=f1): ")
		So(entry.msg, ShouldEqual, "Creating 20 blocks")

		entry, err = r.Next()
		So(err, ShouldBeNil)
		So(entry.header, ShouldEqual, "I0101 00:00:00.000Z std_logger_test.go:00 (Flow=f1): ")
		// Note that the continuation prefix is removed.
		So(entry.msg, ShouldEqual, `Here's some JSON:
{
  "A": "some val",
  "B": "other val"
}
Err:<nil>`)

		entry, err = r.Next()
		So(err, ShouldBeNil)
		So(entry.header, ShouldEqual, "I0101 00:00:00.000Z std_logger_test.go:00 (Block=addFoo): ")
		So(entry.msg, ShouldEqual, "Computing 5 things")

		entry, err = r.Next()
		So(err, ShouldEqual, io.EOF)
	})

}
