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
	std.Output(1, LevelFatal, nil, s)
	panic(s)
}

func Check(cond bool, v ...interface{}) {
	if cond {
		return
	}
	s := fmt.Sprint(v...)
	std.Output(1, LevelFatal, nil, s)
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
	std.output(1, level, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, func(w io.Writer) error {
		return writeText(w, txt)
	})
}

func Text(level Level, txt string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(1, level, func(b []byte) []byte {
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
	std.output(1, level, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, func(w io.Writer) error {
		return writeBinary(w, bin)
	})
}

func Binary(level Level, bin []byte, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	std.output(1, level, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, func(w io.Writer) error {
		return writeBinary(w, bin)
	})
}

func writeStacktrace(w io.Writer, frames *runtime.Frames) error {
	frame, more := frames.Next()
	for index := 1; frame != (runtime.Frame{}); index++ {
		var line string
		if frame.Function != "" && frame.File != "" {
			line = fmt.Sprintf("%s#%-3d 0x%x in %s (%s:%d)\n", indent, index,
				frame.PC, frame.Function, frame.File, frame.Line)
		} else if frame.Function != "" {
			line = fmt.Sprintf("%s#%-3d 0x%x %s+0x%x\n", indent, index,
				frame.PC, frame.Function, frame.PC-frame.Entry)
		} else {
			line = fmt.Sprintf("%s#%-3d 0x%x <unknown>\n", indent, index, frame.PC)
		}
		if _, err := w.Write([]byte(line)); err != nil {
			return err
		}
		if !more {
			break
		}
		frame, more = frames.Next()
	}
	return nil
}

const stackMaxDepth = 256

func Stackf(level Level, calldepth int, format string, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	pc := make([]uintptr, stackMaxDepth)
	n := runtime.Callers(calldepth+2, pc)
	frames := runtime.CallersFrames(pc[:n])
	std.output(1, level, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, func(w io.Writer) error {
		return writeStacktrace(w, frames)
	})
}

func Stack(level Level, calldepth int, v ...interface{}) {
	if !CheckLevel(level) {
		return
	}
	pc := make([]uintptr, stackMaxDepth)
	n := runtime.Callers(calldepth+2, pc)
	frames := runtime.CallersFrames(pc[:n])
	std.output(1, level, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, func(w io.Writer) error {
		return writeStacktrace(w, frames)
	})
}
