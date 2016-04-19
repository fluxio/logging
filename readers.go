package logging

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type FilteringLogReader struct {
	reader        *bufio.Reader
	contextRegexp *regexp.Regexp
	buf           []byte
	skipEntry     bool
}

func (f *FilteringLogReader) Read(buf []byte) (n int, err error) {
	if len(f.buf) > 0 {
		n = copy(buf, f.buf)
		f.buf = f.buf[n:]
		return n, nil
	}

	var prefix bool
	var haveNewLine bool
	for f.skipEntry {
		f.buf, prefix, err = f.reader.ReadLine()
		if err != nil {
			return 0, err
		}

		if haveNewLine {
			// When we've read a new line, check to see if it's a continuation.
			// If so, keep skipping it.
			f.skipEntry = bytes.HasPrefix(f.buf, []byte(continuation))
		}

		haveNewLine = !prefix
	}
	// If we haven't read anything in yet, read in some data.
	if !haveNewLine {
		f.buf, prefix, err = f.reader.ReadLine()
	}

	return 0, io.EOF
}

type LogReader struct {
	reader     *bufio.Reader
	linebuffer bytes.Buffer
}

func NewLogReader(r io.Reader) *LogReader {
	return &LogReader{reader: bufio.NewReader(r)}
}

type logEntry struct {
	full string

	header  string
	typ     byte
	ts      string
	file    string
	line    string
	context string

	msg string
}

func (r *LogReader) Next() (entry logEntry, err error) {
	for err == nil {
		var line string
		line, err = r.reader.ReadString('\n')

		// ReadString will return io.EOF even if line is non-empty if the last
		// line of the file doesn't end in a newline.  But for this function,
		// we only return an io.EOF error if prefix and msg are empty.
		if len(line) > 0 && err == io.EOF {
			err = nil
		} else if err != nil {
			break
		}

		lineType, lineEntry := determineLineType(line)
		if lineType == unknown {
			return entry, fmt.Errorf("Malformatted log file.  Cannot parse line: %q", line)
		}
		if lineType == entryCont && len(lineEntry.header) == 0 {
			return entry, fmt.Errorf("Starting in the middle of a log file: %q", line)
		}

		if lineType == entryStart {
			entry = lineEntry
		} else {
			entry.msg += lineEntry.msg
		}

		next, _ := r.reader.Peek(1)
		if len(next) == 0 || next[0] != ' ' {
			break
		}
	}
	entry.msg = strings.TrimRight(entry.msg, "\r\n")
	return entry, err
}

const (
	unknown = iota
	entryStart
	entryCont
)

var entryStartRegexp = regexp.MustCompile(`^(I|D|E)(\d{4} [\d:\.]{12}[-+\dZ]+) ([\w.]+):(\d+) \((.*)\): `)

func determineLineType(line string) (lineType int, e logEntry) {
	if len(line) < len(continuation) {
		return unknown, e
	}
	if line[:len(continuation)] == string(continuation) {
		return entryCont, logEntry{full: line, header: string(continuation), msg: line[len(continuation):]}
	}

	headerPos := entryStartRegexp.FindStringSubmatchIndex(line)
	if len(headerPos) == 0 {
		return unknown, e
	}

	part := func(n int) string { return line[headerPos[n*2]:headerPos[n*2+1]] }

	e = logEntry{
		full:    line,
		header:  part(0),
		typ:     part(1)[0],
		ts:      part(2),
		file:    part(3),
		line:    part(4),
		context: part(5),
		msg:     line[headerPos[1]:],
	}

	return entryStart, e
}
