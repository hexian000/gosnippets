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
	Write(m *message)
}

const lineBufSize = 8192

type discardWriter struct{}

func newDiscardWriter() writer {
	return &discardWriter{}
}

func (w *discardWriter) Write(*message) {}

type lineWriter struct {
	out io.Writer
}

func newLineWriter(out io.Writer) writer {
	return &lineWriter{
		out: out,
	}
}

func (w *lineWriter) Write(m *message) {
	buf := make([]byte, 0, lineBufSize)
	buf = append(buf, levelChar[m.level], ' ')
	buf = m.timestamp.AppendFormat(buf, time.RFC3339)
	buf = append(buf, ' ')
	buf = append(buf, m.file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(m.line), 10)
	buf = append(buf, ' ')
	buf = append(buf, m.msg...)
	buf = append(buf, '\n')
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
