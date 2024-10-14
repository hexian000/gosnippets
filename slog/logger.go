package slog

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu         sync.RWMutex
	level      Level
	out        writer
	filePrefix string
}

func NewLogger(level Level) *Logger {
	return &Logger{
		out:   newDiscardWriter(),
		level: level,
	}
}

func (l *Logger) SetOutputConfig(output, tag string) error {
	if newOutput, ok := builtinOutput[output]; ok {
		o, err := newOutput(tag)
		if err != nil {
			return err
		}
		l.SetOutput(o)
		return nil
	}
	// otherwise, the string must be a url
	u, err := url.Parse(output)
	if err != nil {
		return fmt.Errorf("unsupported log output: %s", output)
	}
	conn, err := net.Dial(u.Scheme, u.Host)
	if err != nil {
		return err
	}
	l.SetOutput(newTextWriter(conn))
	return nil
}

func (l *Logger) Output(calldepth int, level Level, msg []byte) {
	now := time.Now()
	l.mu.RLock()
	out := l.out
	filePrefix := l.filePrefix
	l.mu.RUnlock()

	_, file, line, ok := runtime.Caller(calldepth)
	if ok {
		file = strings.TrimPrefix(file, filePrefix)
	} else {
		file, line = "???", 0
	}

	out.Write(&message{
		timestamp: now,
		level:     level,
		file:      file,
		line:      line,
		msg:       msg,
	})
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) Level() Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

func (l *Logger) SetFilePrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.filePrefix = prefix
}

func (l *Logger) SetOutput(w writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) SetLineOutput(w io.Writer) {
	l.SetOutput(newTextWriter(w))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if LevelFatal > l.Level() {
		return
	}
	l.Output(2, LevelFatal, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Fatal(v ...interface{}) {
	if LevelFatal > l.Level() {
		return
	}
	l.Output(2, LevelFatal, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if LevelError > l.Level() {
		return
	}
	l.Output(2, LevelError, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Error(v ...interface{}) {
	if LevelError > l.Level() {
		return
	}
	l.Output(2, LevelError, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	if LevelWarning > l.Level() {
		return
	}
	l.Output(2, LevelWarning, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Warning(v ...interface{}) {
	if LevelWarning > l.Level() {
		return
	}
	l.Output(2, LevelWarning, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	if LevelNotice > l.Level() {
		return
	}
	l.Output(2, LevelNotice, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Notice(v ...interface{}) {
	if LevelNotice > l.Level() {
		return
	}
	l.Output(2, LevelNotice, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if LevelInfo > l.Level() {
		return
	}
	l.Output(2, LevelInfo, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Info(v ...interface{}) {
	if LevelInfo > l.Level() {
		return
	}
	l.Output(2, LevelInfo, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if LevelDebug > l.Level() {
		return
	}
	l.Output(2, LevelDebug, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Debug(v ...interface{}) {
	if LevelDebug > l.Level() {
		return
	}
	l.Output(2, LevelDebug, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Verbosef(format string, v ...interface{}) {
	if LevelVerbose > l.Level() {
		return
	}
	l.Output(2, LevelVerbose, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Verbose(v ...interface{}) {
	if LevelVerbose > l.Level() {
		return
	}
	l.Output(2, LevelVerbose, []byte(fmt.Sprint(v...)))
}

func (l *Logger) VeryVerbosef(format string, v ...interface{}) {
	if LevelVeryVerbose > l.Level() {
		return
	}
	l.Output(2, LevelVeryVerbose, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) VeryVerbose(v ...interface{}) {
	if LevelVeryVerbose > l.Level() {
		return
	}
	l.Output(2, LevelVeryVerbose, []byte(fmt.Sprint(v...)))
}
