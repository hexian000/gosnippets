package slog

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

func Checkf(cond bool, format string, v ...interface{}) {
	if cond {
		return
	}
	s := fmt.Sprintf(format, v...)
	std.Output(2, LevelFatal, nil, s)
	panic(s)
}

func Check(cond bool, v ...interface{}) {
	if cond {
		return
	}
	s := fmt.Sprint(v...)
	std.Output(2, LevelFatal, nil, s)
	panic(s)
}

const (
	indent   = "  "
	hardWrap = 70
	tabWidth = 4
)

func writeText(w io.Writer, txt string) error {
	b := make([]byte, 0, 256)
	line := 1
	wrap := 0
	var width int
	for _, r := range txt {
		if wrap == 0 {
			b = fmt.Appendf(b, "%s%4d ", indent, line)
		}
		switch r {
		case '\n':
			/* soft wrap */
			b = append(b, '\n')
			if _, err := w.Write(b); err != nil {
				return err
			}
			b = b[:0]
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
			if _, err := w.Write(b); err != nil {
				return err
			}
			b = b[:0]
			wrap = 0
		}
		if r == '\t' {
			b = append(b, strings.Repeat(" ", width)...)
			wrap += width
			continue
		}
		b = append(b, string(r)...)
		wrap += width
	}
	if wrap > 0 {
		b = append(b, '\n')
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func Textf(level Level, txt string, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(2, level, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, func(w io.Writer) error {
		return writeText(w, txt)
	})
}

func Text(level Level, txt string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(2, level, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, func(w io.Writer) error {
		return writeText(w, txt)
	})
}

func writeBinary(w io.Writer, bin []byte) error {
	b := make([]byte, 0, 256)
	wrap := 16
	for i := 0; i < len(bin); i += wrap {
		b = fmt.Appendf(b, "%s%p: ", indent, bin[i:])
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
		b = append(b, '\n')
		if _, err := w.Write(b); err != nil {
			return err
		}
		b = b[:0]
	}
	return nil
}

func Binaryf(level Level, bin []byte, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(2, level, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, func(w io.Writer) error {
		return writeBinary(w, bin)
	})
}

func Binary(level Level, bin []byte, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(2, level, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, func(w io.Writer) error {
		return writeBinary(w, bin)
	})
}

const stackBufSize = 4096

func Stackf(level Level, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	var stack [stackBufSize]byte
	n := runtime.Stack(stack[:], false)
	std.output(2, level, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, func(w io.Writer) error {
		_, err := w.Write(stack[:n])
		return err
	})
}

func Stack(level Level, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	var stack [stackBufSize]byte
	n := runtime.Stack(stack[:], false)
	std.output(2, level, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, func(w io.Writer) error {
		_, err := w.Write(stack[:n])
		return err
	})
}
