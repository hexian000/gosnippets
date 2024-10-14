package slog

import (
	"bytes"
	"fmt"
	"runtime"
	"unicode"

	"github.com/mattn/go-runewidth"
)

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

const (
	indent   = "  "
	hardWrap = 70
	tabWidth = 4
)

func printText(buf *bytes.Buffer, txt string) {
	line := 1
	wrap := 0
	var width int
	for _, r := range txt {
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
			if !unicode.IsPrint(r) {
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
}

func Textf(level Level, txt string, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	buf := bytes.NewBufferString(fmt.Sprintf(format, v...))
	printText(buf, txt)
	std.Output(2, level, buf.Bytes())
}

func Text(level Level, txt string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	buf := bytes.NewBufferString(fmt.Sprint(v...))
	printText(buf, txt)
	std.Output(2, level, buf.Bytes())
}

func printBinary(buf *bytes.Buffer, bin []byte) {
	wrap := 16
	for i := 0; i < len(bin); i += wrap {
		buf.WriteString(fmt.Sprintf("\n%s%p: ", indent, bin[i:]))
		for j := 0; j < wrap; j++ {
			if (i + j) < len(bin) {
				buf.WriteString(fmt.Sprintf("%02X ", bin[i+j]))
			} else {
				buf.WriteString("   ")
			}
		}
		buf.WriteRune(' ')
		for j := 0; j < wrap; j++ {
			ch := ' '
			if (i + j) < len(bin) {
				ch = rune(bin[i+j])
				if ch > unicode.MaxASCII || !unicode.IsPrint(ch) {
					ch = '.'
				}
			}
			buf.WriteRune(ch)
		}
	}
}

func Binaryf(level Level, bin []byte, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	buf := bytes.NewBufferString(fmt.Sprintf(format, v...))
	printBinary(buf, bin)
	std.Output(2, level, buf.Bytes())
}

func Binary(level Level, bin []byte, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	buf := bytes.NewBufferString(fmt.Sprint(v...))
	printBinary(buf, bin)
	std.Output(2, level, buf.Bytes())
}

func Stackf(level Level, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	buf := bytes.NewBuffer(make([]byte, 0, 8192))
	buf.WriteString(fmt.Sprintf(format, v...))
	buf.WriteRune('\n')
	b := buf.AvailableBuffer()
	n := runtime.Stack(b, false)
	_, _ = buf.Write(b[:n])
	b = buf.Bytes()
	if b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	std.Output(2, level, b)
}

func Stack(level Level, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	buf := bytes.NewBuffer(make([]byte, 0, 8192))
	buf.WriteString(fmt.Sprint(v...))
	buf.WriteRune('\n')
	b := buf.AvailableBuffer()
	n := runtime.Stack(b, false)
	_, _ = buf.Write(b[:n])
	b = buf.Bytes()
	if b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	std.Output(2, level, b)
}
