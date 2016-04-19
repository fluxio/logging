package logging

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func executeWithMaximumContention(N int, f func(i int)) {
	var ready, exec, done sync.WaitGroup
	ready.Add(N)
	exec.Add(1)
	done.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			ready.Done() // indicate that we're ready to go
			exec.Wait()  // wait for everyone to go at the same time
			f(i)         // execute the function
			done.Done()  // let them know we're done
		}(i)
	}
	ready.Wait() // Ensure that all goroutines are ready.
	exec.Done()  // Run them all for maximum conflict.
	done.Wait()  // Wait for them to complete.
}

func TestTextWriter(t *testing.T) {
	var buf bytes.Buffer
	var w = TextWriter{Writer: &buf}

	var ts = (time.Time{}).Add(12345678901 * time.Millisecond)

	Convey("TextWriter", t, func() {
		buf.Reset()
		Convey("should format entries with no format arg correctly", func() {
			w.Write(Entry{InfoLevel, ts, "/path/to/file.js", 12, "ctx", kNO_FORMAT, args("Hi", "there")})
			So(buf.String(), ShouldEqual, "I0523 21:21:18.901Z file.js:12 (ctx): Hi there\n")
		})
		Convey("should format entries with a format arg correctly", func() {
			w.Write(Entry{InfoLevel, ts, "/path/to/file.js", 12, "ctx", "(%s %d %s)", args("Hi", 4, "there")})
			So(buf.String(), ShouldEqual, "I0523 21:21:18.901Z file.js:12 (ctx): (Hi 4 there)\n")
		})
		Convey("should format the log level correctly", func() {
			w.Write(Entry{Level: DebugLevel})
			So(buf.String(), ShouldStartWith, "D")
			buf.Reset()
			w.Write(Entry{Level: InfoLevel})
			So(buf.String(), ShouldStartWith, "I")
			buf.Reset()
			w.Write(Entry{Level: ErrorLevel})
			So(buf.String(), ShouldStartWith, "E")
		})
		Convey("should format entries without line numbers correctly", func() {
			w.Write(Entry{File: "/path/to/file.js", Line: -1})
			So(buf.String(), ShouldContainSubstring, "Z file.js ")
		})
		Convey("should format entries without filename or line numbers correctly", func() {
			w.Write(Entry{File: "", Line: -1})
			So(buf.String(), ShouldContainSubstring, "Z ??? (")
		})
		Convey("should format entries that span multiple lines correctly", func() {
			w.Write(Entry{Fmt: "a\nb\nc"})
			So(buf.String(), ShouldContainSubstring, "a\n"+continuation+"b\n"+continuation+"c\n")
		})
		Convey("should format entries whose context spans multiple lines correctly", func() {
			w.Write(Entry{Context: "a\nb\nc"})
			So(buf.String(), ShouldContainSubstring, "(a\n"+continuation+"b\n"+continuation+"c)")
		})
		Convey("should write messages atomically when written from many threads", func() {
			runtime.GOMAXPROCS(10)
			N := 20
			executeWithMaximumContention(N, func(i int) {
				w.Write(Entry{Level: InfoLevel, Context: fmt.Sprintf("C%d", i), Fmt: fmt.Sprintf("FMT%d", i)})
			})
			lines := strings.Split(buf.String(), "\n")
			So(len(lines), ShouldEqual, N+1) // 10 newlines == 11 output lines
			// Make sure that each of formatted lines is written atomically and
			// not split.
			for i := 0; i < N; i++ {
				expected := fmt.Sprintf("I0101 00:00:00.000Z ???:0 (C%d): FMT%d\n", i, i)
				So(buf.String(), ShouldContainSubstring, expected)
			}
		})
	})
}
