package logging

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"time"

	"github.com/fluxio/iohelpers/line"
)

// If long lines wrap multiple lines, use this prefix for each continuation line
const continuation = "    "

type TextWriter struct {
	Writer io.Writer
	mutex  sync.Mutex
}

func (t *TextWriter) fmtTimestamp(ts time.Time) string { return ts.Format("0102 15:04:05.000Z0700") }
func (t *TextWriter) fmtOrigin(file string, line int) string {
	if file == "" {
		file = "???"
	} else {
		file = filepath.Base(file)
	}
	if line == -1 {
		return file
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func (t *TextWriter) Write(e Entry) error {
	// First we construct an in-memory string of the of entry.  This
	// can be done in parallel for all calling threads.
	var buf bytes.Buffer
	w := line.PrefixWriter{&buf, []byte(continuation), true}
	// Prefix
	fmt.Fprintf(&w, "%s%s %s (%s): ", e.Level, t.fmtTimestamp(e.Time),
		t.fmtOrigin(e.File, e.Line), e.Context)
	// Content
	if e.Fmt == kNO_FORMAT {
		fmt.Fprintln(&w, e.Args...)
	} else {
		fmt.Fprintf(&w, e.Fmt+"\n", e.Args...)
	}

	// Then we lock and write to the final output writer in one go.
	t.mutex.Lock()
	_, err := buf.WriteTo(t.Writer)
	t.mutex.Unlock()

	return err
}
