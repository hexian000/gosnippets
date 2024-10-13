package slog

import (
	"bytes"
	"fmt"
	"runtime"
	"unicode"

	runewidth "github.com/mattn/go-runewidth"
)

const (
	LevelSilence Level = iota
	LevelFatal
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug
	LevelVerbose
	LevelVeryVerbose
)

type Level int

var levelChar = [...]byte{
	'-', 'F', 'E', 'W', 'I', 'I', 'D', 'V', 'V',
}

var std = NewLogger(LevelVerbose)

func Default() *Logger {
	return std
}

func Checkf(cond bool, format string, v ...interface{}) {
	if cond {
		return
	}
	s := fmt.Sprintf(format, v...)
	std.Output(2, LevelFatal, []byte(s))
	panic(s)
}

func Check(cond bool, v ...interface{}) {
	if cond {
		return
	}
	s := fmt.Sprint(v...)
	std.Output(2, LevelFatal, []byte(s))
	panic(s)
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Output(2, LevelFatal, []byte(s))
	panic(s)
}

func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(2, LevelFatal, []byte(s))
	panic(s)
}

// Serious problems that are likely to cause the program to exit.
func Fatalf(format string, v ...interface{}) {
	std.Output(2, LevelFatal, []byte(fmt.Sprintf(format, v...)))
}

// Serious problems that are likely to cause the program to exit.
func Fatal(v ...interface{}) {
	std.Output(2, LevelFatal, []byte(fmt.Sprint(v...)))
}

// Issues that shouldn't be ignored.
func Errorf(format string, v ...interface{}) {
	std.Output(2, LevelError, []byte(fmt.Sprintf(format, v...)))
}

// Issues that shouldn't be ignored.
func Error(v ...interface{}) {
	std.Output(2, LevelError, []byte(fmt.Sprint(v...)))
}

// Issues that may be ignored.
func Warningf(format string, v ...interface{}) {
	std.Output(2, LevelWarning, []byte(fmt.Sprintf(format, v...)))
}

// Issues that may be ignored.
func Warning(v ...interface{}) {
	std.Output(2, LevelWarning, []byte(fmt.Sprint(v...)))
}

// Important status changes. The prefix is 'I'.
func Noticef(format string, v ...interface{}) {
	std.Output(2, LevelNotice, []byte(fmt.Sprintf(format, v...)))
}

// Important status changes. The prefix is 'I'.
func Notice(v ...interface{}) {
	std.Output(2, LevelNotice, []byte(fmt.Sprint(v...)))
}

// Normal work reports.
func Infof(format string, v ...interface{}) {
	std.Output(2, LevelInfo, []byte(fmt.Sprintf(format, v...)))
}

// Normal work reports.
func Info(v ...interface{}) {
	std.Output(2, LevelInfo, []byte(fmt.Sprint(v...)))
}

// Extra information for debugging.
func Debugf(format string, v ...interface{}) {
	std.Output(2, LevelDebug, []byte(fmt.Sprintf(format, v...)))
}

// Extra information for debugging.
func Debug(v ...interface{}) {
	std.Output(2, LevelDebug, []byte(fmt.Sprint(v...)))
}

// Details for inspecting specific issues.
func Verbosef(format string, v ...interface{}) {
	std.Output(2, LevelVerbose, []byte(fmt.Sprintf(format, v...)))
}

// Details for inspecting specific issues.
func Verbose(v ...interface{}) {
	std.Output(2, LevelVerbose, []byte(fmt.Sprint(v...)))
}

// More details that may significantly impact performance. The prefix is 'V'.
func VeryVerbosef(format string, v ...interface{}) {
	std.Output(2, LevelVeryVerbose, []byte(fmt.Sprintf(format, v...)))
}

// More details that may significantly impact performance. The prefix is 'V'.
func VeryVerbose(v ...interface{}) {
	std.Output(2, LevelVeryVerbose, []byte(fmt.Sprint(v...)))
}

const (
	indent   = "  "
	hardWrap = 70
	tabWidth = 4
)

func Textf(level Level, text string, format string, v ...interface{}) {
	buf := bytes.NewBufferString(fmt.Sprintf(format, v...))
	line := 1
	wrap := 0
	var width int
	for _, r := range text {
		if wrap == 0 {
			buf.WriteString(fmt.Sprintf("\n%s%4d ", indent, line))
		}
		switch r {
		case '\n':
			/* soft wrap */
			line++
			wrap = 0
			continue
		case '\t':
			width = tabWidth - wrap%tabWidth
		default:
			if !(unicode.IsPrint(r) || unicode.IsSpace(r)) {
				r = '?'
			}
			width = runewidth.RuneWidth(r)
		}
		if wrap+width > hardWrap {
			/* hard wrap */
			buf.WriteString(fmt.Sprintf(" +\n%s     ", indent))
			wrap = 0
		}
		if r == '\t' {
			for i := 0; i < width; i++ {
				buf.WriteRune(' ')
			}
			wrap += width
			continue
		}
		_, _ = buf.WriteRune(r)
		wrap += width
	}
	std.Output(2, level, buf.Bytes())
}

func Binaryf(level Level, bin []byte, format string, v ...interface{}) {
	// TODO
}

func Stackf(level Level, format string, v ...interface{}) {
	var buf [16384]byte
	b := append(buf[:0], []byte(fmt.Sprintf(format, v...))...)
	b = append(b, '\n')
	n := runtime.Stack(buf[len(b):], false)
	b = buf[:len(b)+n]
	if b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	std.Output(2, level, b)
}

func Write(calldepth int, level Level, msg []byte) {
	std.Output(calldepth+1, level, msg)
}
