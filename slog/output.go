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
	file      []byte
	line      int
	msg       []byte
}

type writer interface {
	Write(m message)
}

type discardWriter struct{}

func newDiscardWriter() writer {
	return &discardWriter{}
}

func (w *discardWriter) Write(_ message) {}

type lineWriter struct {
	buf []byte
	out io.Writer
}

func newLineWriter(out io.Writer) writer {
	return &lineWriter{
		buf: make([]byte, 0),
		out: out,
	}
}

func (w *lineWriter) Write(m message) {
	buf := w.buf[:0]
	buf = append(buf, levelChar[m.level], ' ')
	buf = m.timestamp.AppendFormat(buf, time.RFC3339)
	buf = append(buf, ' ')
	buf = append(buf, m.file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(m.line), 10)
	buf = append(buf, ' ')
	buf = append(buf, m.msg...)
	buf = append(buf, '\n')
	w.buf = buf
	_, _ = w.out.Write(buf)
}

var builtinOutput map[string]func(tag string) (writer, error)

func init() {
	builtinOutput = map[string]func(tag string) (writer, error){
		"discard": func(string) (writer, error) {
			return newDiscardWriter(), nil
		},
		"stdout": func(string) (writer, error) {
			return newLineWriter(os.Stdout), nil
		},
		"stderr": func(string) (writer, error) {
			return newLineWriter(os.Stderr), nil
		},
	}
}
