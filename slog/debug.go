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
	line := 1
	wrap := 0
	var width int
	for _, r := range txt {
		if wrap == 0 {
			if _, err := w.Write([]byte(fmt.Sprintf("%s%4d ", indent, line))); err != nil {
				return err
			}
		}
		switch r {
		case '\n':
			/* soft wrap */
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
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
			if _, err := w.Write([]byte(fmt.Sprintf(" +\n%s     ", indent))); err != nil {
				return err
			}
			wrap = 0
		}
		if r == '\t' {
			if _, err := w.Write([]byte(strings.Repeat(" ", width))); err != nil {
				return err
			}
			wrap += width
			continue
		}
		if _, err := w.Write([]byte(string(r))); err != nil {
			return err
		}
		wrap += width
	}
	if wrap > 0 {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
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
	wrap := 16
	for i := 0; i < len(bin); i += wrap {
		if _, err := w.Write([]byte(fmt.Sprintf("%s%p: ", indent, bin[i:]))); err != nil {
			return err
		}
		for j := 0; j < wrap; j++ {
			if (i + j) < len(bin) {
				if _, err := w.Write([]byte(fmt.Sprintf("%02X ", bin[i+j]))); err != nil {
					return err
				}
			} else {
				if _, err := w.Write([]byte("   ")); err != nil {
					return err
				}
			}
		}
		if _, err := w.Write([]byte(" ")); err != nil {
			return err
		}
		for j := 0; j < wrap; j++ {
			r := ' '
			if (i + j) < len(bin) {
				r = rune(bin[i+j])
				if r > unicode.MaxASCII || !unicode.IsPrint(r) {
					r = '.'
				}
			}
			if _, err := w.Write([]byte(string(r))); err != nil {
				return err
			}
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
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
