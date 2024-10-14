package slog

import (
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
	std.Output(2, LevelFatal, s)
	panic(s)
}

func Check(cond bool, v ...interface{}) {
	if cond {
		return
	}
	s := fmt.Sprint(v...)
	std.Output(2, LevelFatal, s)
	panic(s)
}

const (
	indent   = "  "
	hardWrap = 70
	tabWidth = 4
)

func appendText(b []byte, txt string) []byte {
	line := 1
	wrap := 0
	var width int
	for _, r := range txt {
		if wrap == 0 {
			b = fmt.Appendf(b, "\n%s%4d ", indent, line)
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
			b = fmt.Appendf(b, " +\n%s     ", indent)
			wrap = 0
		}
		if r == '\t' {
			for i := 0; i < width; i++ {
				b = append(b, ' ')
			}
			wrap += width
			continue
		}
		b = append(b, string(r)...)
		wrap += width
	}
	return b
}

func Textf(level Level, txt string, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(2, level, func(b []byte) []byte {
		b = fmt.Appendf(b, format, v...)
		return appendText(b, txt)
	})
}

func Text(level Level, txt string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(2, level, func(b []byte) []byte {
		b = fmt.Append(b, v...)
		return appendText(b, txt)
	})
}

func appendBinary(b []byte, bin []byte) []byte {
	wrap := 16
	for i := 0; i < len(bin); i += wrap {
		b = fmt.Appendf(b, "\n%s%p: ", indent, bin[i:])
		for j := 0; j < wrap; j++ {
			if (i + j) < len(bin) {
				b = fmt.Appendf(b, "%02X ", bin[i+j])
			} else {
				b = append(b, "   "...)
			}
		}
		b = append(b, ' ')
		for j := 0; j < wrap; j++ {
			r := ' '
			if (i + j) < len(bin) {
				r = rune(bin[i+j])
				if r > unicode.MaxASCII || !unicode.IsPrint(r) {
					r = '.'
				}
			}
			b = append(b, string(r)...)
		}
	}
	return b
}

func Binaryf(level Level, bin []byte, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(2, level, func(b []byte) []byte {
		b = fmt.Appendf(b, format, v...)
		return appendBinary(b, bin)
	})
}

func Binary(level Level, bin []byte, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(2, level, func(b []byte) []byte {
		b = fmt.Append(b, v...)
		return appendBinary(b, bin)
	})
}

const stackBufSize = 4096

func Stackf(level Level, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	var stack [stackBufSize]byte
	n := runtime.Stack(stack[:], false)
	if stack[n-1] == '\n' {
		n--
	}
	std.output(2, level, func(b []byte) []byte {
		b = fmt.Appendf(b, format, v...)
		b = append(b, '\n')
		return append(b, stack[:n]...)
	})
}

func Stack(level Level, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	var stack [stackBufSize]byte
	n := runtime.Stack(stack[:], false)
	if stack[n-1] == '\n' {
		n--
	}
	std.output(2, level, func(b []byte) []byte {
		b = fmt.Append(b, v...)
		b = append(b, '\n')
		return append(b, stack[:n]...)
	})
}
