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

type Logger struct {
	out        output
	outMu      sync.Mutex
	level      atomic.Int32
	filePrefix atomic.Pointer[string]
}

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

func (l *Logger) SetOutput(t OutputType, v ...interface{}) {
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
	_, file, line, ok := runtime.Caller(calldepth)
	if ok {
		file = strings.TrimPrefix(file, *l.filePrefix.Load())
	} else {
		file, line = "???", 0
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

func (l *Logger) Output(calldepth int, level Level, s string, writeExtra func(io.Writer) error) error {
	return l.output(calldepth+1, level, func(b []byte) []byte {
		return append(b, s...)
	}, writeExtra)
}

func (l *Logger) SetLevel(level Level) {
	l.level.Store(int32(level))
}

func (l *Logger) Level() Level {
	return Level(l.level.Load())
}

func (l *Logger) SetFilePrefix(prefix string) {
	l.filePrefix.Store(&prefix)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if LevelFatal > l.Level() {
		return
	}
	l.output(2, LevelFatal, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

func (l *Logger) Fatal(v ...interface{}) {
	if LevelFatal > l.Level() {
		return
	}
	l.output(2, LevelFatal, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if LevelError > l.Level() {
		return
	}
	l.output(2, LevelError, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

func (l *Logger) Error(v ...interface{}) {
	if LevelError > l.Level() {
		return
	}
	l.output(2, LevelError, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	if LevelWarning > l.Level() {
		return
	}
	l.output(2, LevelWarning, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

func (l *Logger) Warning(v ...interface{}) {
	if LevelWarning > l.Level() {
		return
	}
	l.output(2, LevelWarning, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	if LevelNotice > l.Level() {
		return
	}
	l.output(2, LevelNotice, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

func (l *Logger) Notice(v ...interface{}) {
	if LevelNotice > l.Level() {
		return
	}
	l.output(2, LevelNotice, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if LevelInfo > l.Level() {
		return
	}
	l.output(2, LevelInfo, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

func (l *Logger) Info(v ...interface{}) {
	if LevelInfo > l.Level() {
		return
	}
	l.output(2, LevelInfo, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if LevelDebug > l.Level() {
		return
	}
	l.output(2, LevelDebug, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

func (l *Logger) Debug(v ...interface{}) {
	if LevelDebug > l.Level() {
		return
	}
	l.output(2, LevelDebug, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

func (l *Logger) Verbosef(format string, v ...interface{}) {
	if LevelVerbose > l.Level() {
		return
	}
	l.output(2, LevelVerbose, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

func (l *Logger) Verbose(v ...interface{}) {
	if LevelVerbose > l.Level() {
		return
	}
	l.output(2, LevelVerbose, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

func (l *Logger) VeryVerbosef(format string, v ...interface{}) {
	if LevelVeryVerbose > l.Level() {
		return
	}
	l.output(2, LevelVeryVerbose, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

func (l *Logger) VeryVerbose(v ...interface{}) {
	if LevelVeryVerbose > l.Level() {
		return
	}
	l.output(2, LevelVeryVerbose, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}
