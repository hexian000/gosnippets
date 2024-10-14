package slog

import (
	"io"
	"strconv"
	"time"
)

var newSyslogWriter func(string) output

type message struct {
	timestamp    time.Time
	level        Level
	file         string
	line         int
	appendOutput func([]byte) []byte
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

type textWriter struct {
	out io.Writer
}

func newTextWriter(out io.Writer) output {
	return &textWriter{
		out: out,
	}
}

func (w *textWriter) Write(m *message) error {
	buf := make([]byte, 0, bufSize)
	buf = append(buf, levelChar[m.level], ' ')
	buf = m.timestamp.AppendFormat(buf, time.RFC3339)
	buf = append(buf, ' ')
	buf = append(buf, m.file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(m.line), 10)
	buf = append(buf, ' ')
	buf = m.appendOutput(buf)
	buf = append(buf, '\n')
	_, err := w.out.Write(buf)
	return err
}
