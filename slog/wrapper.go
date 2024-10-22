// gosnippets (c) 2023-2024 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package slog

import "log"

type wrapper struct {
	l     *Logger
	level Level
}

func (w *wrapper) Write(p []byte) (n int, err error) {
	const calldepth = 4
	w.l.Output(calldepth, w.level, nil, string(p))
	return len(p), nil
}

func Wrap(logger *Logger, level Level) *log.Logger {
	return log.New(&wrapper{logger, level}, "", 0)
}
