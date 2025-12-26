// gosnippets (c) 2023-2025 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package slog

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Logger represents a logger instance.
type Logger struct {
	out        output
	outMu      sync.Mutex
	level      atomic.Int32
	filePrefix atomic.Pointer[string]
}

// NewLogger creates and returns a new Logger instance.
func NewLogger() *Logger {
	return &Logger{
		out: newDiscardWriter(),
	}
}

type OutputType int

const (
	OutputDiscard OutputType = iota
	OutputWriter
	OutputSyslog
)

// SetOutput sets the output type and parameters for the logger.
func (l *Logger) SetOutput(t OutputType, v ...any) {
	var w output
	switch t {
	case OutputDiscard:
		w = newDiscardWriter()
	case OutputWriter:
		w = newTextWriter(v[0].(io.Writer))
	case OutputSyslog:
		if newSyslogWriter == nil {
			w = newDiscardWriter()
		} else {
			w = newSyslogWriter(v[0].(string))
		}
	}
	l.outMu.Lock()
	defer l.outMu.Unlock()
	l.out = w
}

func (l *Logger) output(calldepth int, level Level, appendMessage func([]byte) []byte, writeExtra func(io.Writer) error) error {
	now := time.Now()
	_, file, line, ok := runtime.Caller(calldepth + 1)
	if !ok {
		file, line = "???", 0
	} else if filePrefix := l.filePrefix.Load(); filePrefix != nil {
		file = strings.TrimPrefix(file, *filePrefix)
	}

	l.outMu.Lock()
	defer l.outMu.Unlock()
	return l.out.Write(&message{
		timestamp:     now,
		level:         level,
		file:          file,
		line:          line,
		appendMessage: appendMessage,
		writeExtra:    writeExtra,
	})
}

// Outputf is the low-level interface to write arbitary log messages.
func (l *Logger) Outputf(calldepth int, level Level, extra func(io.Writer) error, format string, v ...interface{}) error {
	return l.output(calldepth+1, level, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, extra)
}

// Output is the low-level interface to write arbitary log messages.
func (l *Logger) Output(calldepth int, level Level, extra func(io.Writer) error, v ...any) error {
	return l.output(calldepth+1, level, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, extra)
}

// SetLevel sets the logging level for the logger.
func (l *Logger) SetLevel(level Level) {
	l.level.Store(int32(level))
}

// Level returns the current logging level of the logger.
func (l *Logger) Level() Level {
	return Level(l.level.Load())
}

// SetFilePrefix sets the file prefix to be stripped from file paths in log messages.
func (l *Logger) SetFilePrefix(prefix string) {
	l.filePrefix.Store(&prefix)
}

// Fatalf logs serious problems that are likely to cause the program to exit.
func (l *Logger) Fatalf(format string, v ...any) {
	if LevelFatal > l.Level() {
		return
	}
	l.output(1, LevelFatal, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Fatal logs serious problems that are likely to cause the program to exit.
func (l *Logger) Fatal(v ...any) {
	if LevelFatal > l.Level() {
		return
	}
	l.output(1, LevelFatal, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Errorf logs issues that shouldn't be ignored.
func (l *Logger) Errorf(format string, v ...any) {
	if LevelError > l.Level() {
		return
	}
	l.output(1, LevelError, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Error logs issues that shouldn't be ignored.
func (l *Logger) Error(v ...any) {
	if LevelError > l.Level() {
		return
	}
	l.output(1, LevelError, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Warningf logs issues that may be ignored.
func (l *Logger) Warningf(format string, v ...any) {
	if LevelWarning > l.Level() {
		return
	}
	l.output(1, LevelWarning, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Warning logs issues that may be ignored.
func (l *Logger) Warning(v ...any) {
	if LevelWarning > l.Level() {
		return
	}
	l.output(1, LevelWarning, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Noticef logs important status changes. The prefix is 'I'.
func (l *Logger) Noticef(format string, v ...any) {
	if LevelNotice > l.Level() {
		return
	}
	l.output(1, LevelNotice, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Notice logs important status changes. The prefix is 'I'.
func (l *Logger) Notice(v ...any) {
	if LevelNotice > l.Level() {
		return
	}
	l.output(1, LevelNotice, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Infof logs normal work reports.
func (l *Logger) Infof(format string, v ...any) {
	if LevelInfo > l.Level() {
		return
	}
	l.output(1, LevelInfo, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Info logs normal work reports.
func (l *Logger) Info(v ...any) {
	if LevelInfo > l.Level() {
		return
	}
	l.output(1, LevelInfo, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Debugf logs extra information for debugging.
func (l *Logger) Debugf(format string, v ...any) {
	if LevelDebug > l.Level() {
		return
	}
	l.output(1, LevelDebug, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Debug logs extra information for debugging.
func (l *Logger) Debug(v ...any) {
	if LevelDebug > l.Level() {
		return
	}
	l.output(1, LevelDebug, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Verbosef logs details for inspecting specific issues.
func (l *Logger) Verbosef(format string, v ...any) {
	if LevelVerbose > l.Level() {
		return
	}
	l.output(1, LevelVerbose, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Verbose logs details for inspecting specific issues.
func (l *Logger) Verbose(v ...any) {
	if LevelVerbose > l.Level() {
		return
	}
	l.output(1, LevelVerbose, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// VeryVerbosef logs more details that may significantly impact performance. The prefix is 'V'.
func (l *Logger) VeryVerbosef(format string, v ...any) {
	if LevelVeryVerbose > l.Level() {
		return
	}
	l.output(1, LevelVeryVerbose, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// VeryVerbose logs more details that may significantly impact performance. The prefix is 'V'.
func (l *Logger) VeryVerbose(v ...any) {
	if LevelVeryVerbose > l.Level() {
		return
	}
	l.output(1, LevelVeryVerbose, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}
