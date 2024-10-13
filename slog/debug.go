package slog

import (
	"bytes"
	"fmt"
	"runtime"
	"unicode"

	"github.com/mattn/go-runewidth"
)

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
