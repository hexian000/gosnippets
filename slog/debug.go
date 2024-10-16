package slog

import (
	"fmt"
	"io"
	"runtime"
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
	hardWrap = 110
	tabSpace = "    "
)

func writeText(w io.Writer, txt string) error {
	var buf [256]byte
	b := buf[:]
	lineno := true
	line, column := 0, 0
	for _, r := range txt {
		if lineno {
			line++
			b = fmt.Appendf(b, "%s%4d ", indent, line)
			lineno = false
		}
		var width int
		switch r {
		case '\n':
			/* soft wrap */
			b = append(b, '\n')
			if _, err := w.Write(b); err != nil {
				return err
			}
			b = b[:0]
			column = 0
			lineno = true
			continue
		case '\t':
			width = len(tabSpace) - column%len(tabSpace)
		default:
			if !unicode.IsPrint(r) {
				r = '?'
			}
			width = runewidth.RuneWidth(r)
		}
		if column+width > hardWrap {
			/* hard wrap */
			b = fmt.Appendf(b, " +\n%s     ", indent)
			if _, err := w.Write(b); err != nil {
				return err
			}
			b = b[:0]
			column = 0
		}
		if r == '\t' {
			b = append(b, tabSpace[:width]...)
			column += width
		} else {
			b = append(b, string(r)...)
			column += width
		}
		if cap(b)-len(b) < 16 {
			if _, err := w.Write(b); err != nil {
				return err
			}
			b = b[:0]
		}
	}
	if column > 0 {
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
	var buf [256]byte
	b := buf[:]
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
