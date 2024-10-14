package slog

import (
	"io"
	"os"
	"strconv"
	"time"
)

type message struct {
	timestamp time.Time
	level     Level
	file      string
	line      int
	msg       []byte
}

type writer interface {
	Write(m *message)
}

const bufSize = 4096

type discardWriter struct{}

func newDiscardWriter() writer {
	return &discardWriter{}
}

func (w *discardWriter) Write(*message) {}

type textWriter struct {
	out io.Writer
}

func newTextWriter(out io.Writer) writer {
	return &textWriter{
		out: out,
	}
}

func (w *textWriter) Write(m *message) {
	buf := make([]byte, 0, bufSize)
	buf = append(buf, levelChar[m.level], ' ')
	buf = m.timestamp.AppendFormat(buf, time.RFC3339)
	buf = append(buf, ' ')
	buf = append(buf, m.file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(m.line), 10)
	buf = append(buf, ' ')
	_, _ = w.out.Write(buf)
	msg := append(m.msg, '\n')
	_, _ = w.out.Write(msg)
}

var builtinOutput map[string]func(tag string) (writer, error)

func init() {
	builtinOutput = map[string]func(tag string) (writer, error){
		"discard": func(string) (writer, error) {
			return newDiscardWriter(), nil
		},
		"stdout": func(string) (writer, error) {
			return newTextWriter(os.Stdout), nil
		},
		"stderr": func(string) (writer, error) {
			return newTextWriter(os.Stderr), nil
		},
	}
}
