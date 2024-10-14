package slog

import "log"

type wrapper struct {
	l     *Logger
	level Level
}

func (w *wrapper) Write(p []byte) (n int, err error) {
	const calldepth = 4
	w.l.Output(calldepth, w.level, string(p))
	return len(p), nil
}

func Wrap(logger *Logger, level Level) *log.Logger {
	return log.New(&wrapper{logger, level}, "", 0)
}
