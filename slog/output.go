// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package slog

import (
	"io"
	"strconv"
	"time"
)

var newSyslogWriter func(string) output

type message struct {
	timestamp     time.Time
	level         Level
	file          string
	line          int
	appendMessage func([]byte) []byte
	writeExtra    func(io.Writer) error
}

type output interface {
	Write(m *message) error
}

const bufSize = 4096

type discardWriter struct{}

func newDiscardWriter() output {
	return &discardWriter{}
}

func (w *discardWriter) Write(*message) error {
	return nil
}

type flusher interface {
	Flush() error
}

type textWriter struct {
	out   io.Writer
	flush func() error
}

func newTextWriter(out io.Writer) output {
	flush := func() error {
		return nil
	}
	if f, ok := out.(flusher); ok {
		flush = func() error {
			return f.Flush()
		}
	}
	return &textWriter{
		out:   out,
		flush: flush,
	}
}

type termWriter textWriter

func newTermWriter(out io.Writer) output {
	flush := func() error {
		return nil
	}
	if f, ok := out.(flusher); ok {
		flush = func() error {
			return f.Flush()
		}
	}
	return &termWriter{
		out:   out,
		flush: flush,
	}
}

/* TimeLayout is a fixed-length layout conforming to both ISO 8601 and RFC 3339 */
const TimeLayout = "2006-01-02T15:04:05-07:00"

func (w *textWriter) Write(m *message) error {
	buf := make([]byte, 0, bufSize)
	buf = append(buf, levelChar[m.level], ' ')
	buf = m.timestamp.AppendFormat(buf, TimeLayout)
	buf = append(buf, ' ')
	buf = append(buf, m.file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(m.line), 10)
	buf = append(buf, ' ')
	buf = m.appendMessage(buf)
	buf = append(buf, '\n')
	if _, err := w.out.Write(buf); err != nil {
		return err
	}
	if m.writeExtra != nil {
		if err := m.writeExtra(w.out); err != nil {
			return err
		}
	}
	return w.flush()
}

func (w *termWriter) Write(m *message) error {
	buf := make([]byte, 0, bufSize)
	buf = append(buf, "\x1b["...) // ESC [
	buf = append(buf, levelColor[m.level]...)
	buf = append(buf, levelChar[m.level], ' ')
	buf = m.timestamp.AppendFormat(buf, TimeLayout)
	buf = append(buf, ' ')
	buf = append(buf, m.file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(m.line), 10)
	buf = append(buf, ' ')
	buf = m.appendMessage(buf)
	buf = append(buf, "\x1b[0m"...)
	buf = append(buf, '\n')
	if _, err := w.out.Write(buf); err != nil {
		return err
	}
	if m.writeExtra != nil {
		if err := m.writeExtra(w.out); err != nil {
			return err
		}
	}
	return w.flush()
}
