package slog

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu    sync.RWMutex
	out   writer
	ch    chan *writeOp
	wg    sync.WaitGroup
	level int
}

func NewLogger(level int) *Logger {
	l := &Logger{
		out:   newDiscardWriter(),
		ch:    make(chan *writeOp, bufSize),
		level: level,
	}
	l.wg.Add(1)
	go l.run()
	return l
}

func (l *Logger) Close() {
	close(l.ch)
	l.wg.Wait()
}

const bufSize = 16

type writeOp struct {
	out writer
	m   message
}

func (l *Logger) run() {
	defer l.wg.Done()
	for op := range l.ch {
		op.out.Write(op.m)
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
	l.SetOutput(newLineWriter(conn))
	return nil
}

func (l *Logger) Write(calldepth int, level int, msg []byte) {
	now := time.Now()
	l.mu.RLock()
	if level > l.level {
		return
	}
	out := l.out
	l.mu.RUnlock()

	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file, line = "???", 0
	} else if cwd, err := os.Getwd(); err == nil {
		file = strings.TrimPrefix(file, cwd+"/")
	}

	l.ch <- &writeOp{
		out: out,
		m: message{
			timestamp: now,
			level:     level,
			file:      []byte(file),
			line:      line,
			msg:       msg,
		},
	}
}

func (l *Logger) SetLevel(level int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) CheckLevel(level int) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return level <= l.level
}

func (l *Logger) SetOutput(w writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) SetLineOutput(w io.Writer) {
	l.SetOutput(newLineWriter(w))
}

func (l *Logger) Checkf(cond bool, format string, v ...interface{}) {
	if cond {
		return
	}
	s := fmt.Sprintf(format, v...)
	l.Write(2, LevelFatal, []byte(s))
	panic(s)
}

func (l *Logger) Check(cond bool, v ...interface{}) {
	if cond {
		return
	}
	s := fmt.Sprint(v...)
	l.Write(2, LevelFatal, []byte(s))
	panic(s)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Write(2, LevelFatal, []byte(s))
	panic(s)
}

func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Write(2, LevelFatal, []byte(s))
	panic(s)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Write(2, LevelFatal, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Write(2, LevelFatal, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Write(2, LevelError, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Error(v ...interface{}) {
	l.Write(2, LevelError, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Write(2, LevelWarning, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Warning(v ...interface{}) {
	l.Write(2, LevelWarning, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.Write(2, LevelNotice, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Notice(v ...interface{}) {
	l.Write(2, LevelNotice, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Write(2, LevelInfo, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Info(v ...interface{}) {
	l.Write(2, LevelInfo, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Write(2, LevelDebug, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Debug(v ...interface{}) {
	l.Write(2, LevelDebug, []byte(fmt.Sprint(v...)))
}

func (l *Logger) Verbosef(format string, v ...interface{}) {
	l.Write(2, LevelVerbose, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) Verbose(v ...interface{}) {
	l.Write(2, LevelVerbose, []byte(fmt.Sprint(v...)))
}

func (l *Logger) VeryVerbosef(format string, v ...interface{}) {
	l.Write(2, LevelVeryVerbose, []byte(fmt.Sprintf(format, v...)))
}

func (l *Logger) VeryVerbose(v ...interface{}) {
	l.Write(2, LevelVeryVerbose, []byte(fmt.Sprint(v...)))
}
