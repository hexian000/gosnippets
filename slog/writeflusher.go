package slog

import "io"

type WriteFlusher interface {
	io.Writer
	Flush() error
}

type writeFlusher struct{ w io.Writer }

func (wf *writeFlusher) Write(p []byte) (n int, err error) { return wf.w.Write(p) }
func (wf *writeFlusher) Flush() error                      { return nil }

func NewWriteFlusher(w io.Writer) WriteFlusher {
	if wf, ok := w.(WriteFlusher); ok {
		return wf
	}
	return &writeFlusher{w: w}
}
