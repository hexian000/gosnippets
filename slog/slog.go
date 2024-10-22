// gosnippets (c) 2023-2024 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package slog

import (
	"fmt"
	"io"
)

type Level int32

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

var levelChar = [...]byte{
	'-', 'F', 'E', 'W', 'I', 'I', 'D', 'V', 'V',
}

var std = NewLogger()

func Default() *Logger {
	return std
}

func Outputf(calldepth int, level Level, extra func(io.Writer) error, format string, v ...interface{}) error {
	return std.output(calldepth+1, level, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, extra)
}

func Output(calldepth int, level Level, extra func(io.Writer) error, v ...interface{}) error {
	return std.output(calldepth+1, level, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, extra)
}

func CheckLevel(level Level) bool {
	return level <= std.Level()
}

// Serious problems that are likely to cause the program to exit.
func Fatalf(format string, v ...interface{}) {
	if !CheckLevel(LevelFatal) {
		return
	}
	std.output(1, LevelFatal, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Serious problems that are likely to cause the program to exit.
func Fatal(v ...interface{}) {
	if !CheckLevel(LevelFatal) {
		return
	}
	std.output(1, LevelFatal, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Issues that shouldn't be ignored.
func Errorf(format string, v ...interface{}) {
	if !CheckLevel(LevelError) {
		return
	}
	std.output(1, LevelError, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Issues that shouldn't be ignored.
func Error(v ...interface{}) {
	if !CheckLevel(LevelError) {
		return
	}
	std.output(1, LevelError, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Issues that may be ignored.
func Warningf(format string, v ...interface{}) {
	if !CheckLevel(LevelWarning) {
		return
	}
	std.output(1, LevelWarning, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Issues that may be ignored.
func Warning(v ...interface{}) {
	if !CheckLevel(LevelWarning) {
		return
	}
	std.output(1, LevelWarning, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Important status changes. The prefix is 'I'.
func Noticef(format string, v ...interface{}) {
	if !CheckLevel(LevelNotice) {
		return
	}
	std.output(1, LevelNotice, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Important status changes. The prefix is 'I'.
func Notice(v ...interface{}) {
	if !CheckLevel(LevelNotice) {
		return
	}
	std.output(1, LevelNotice, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Normal work reports.
func Infof(format string, v ...interface{}) {
	if !CheckLevel(LevelInfo) {
		return
	}
	std.output(1, LevelInfo, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Normal work reports.
func Info(v ...interface{}) {
	if !CheckLevel(LevelInfo) {
		return
	}
	std.output(1, LevelInfo, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Extra information for debugging.
func Debugf(format string, v ...interface{}) {
	if !CheckLevel(LevelDebug) {
		return
	}
	std.output(1, LevelDebug, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Extra information for debugging.
func Debug(v ...interface{}) {
	if !CheckLevel(LevelDebug) {
		return
	}
	std.output(1, LevelDebug, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Details for inspecting specific issues.
func Verbosef(format string, v ...interface{}) {
	if !CheckLevel(LevelVerbose) {
		return
	}
	std.output(1, LevelVerbose, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Details for inspecting specific issues.
func Verbose(v ...interface{}) {
	if !CheckLevel(LevelVerbose) {
		return
	}
	std.output(1, LevelVerbose, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// More details that may significantly impact performance. The prefix is 'V'.
func VeryVerbosef(format string, v ...interface{}) {
	if !CheckLevel(LevelVeryVerbose) {
		return
	}
	std.output(1, LevelVeryVerbose, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// More details that may significantly impact performance. The prefix is 'V'.
func VeryVerbose(v ...interface{}) {
	if !CheckLevel(LevelVeryVerbose) {
		return
	}
	std.output(1, LevelVeryVerbose, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}
