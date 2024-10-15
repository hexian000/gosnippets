package slog

import (
	"fmt"
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

func Output(calldepth int, level Level, s string) {
	std.Output(calldepth+1, level, s, nil)
}

func CheckLevel(level Level) bool {
	return level <= std.Level()
}

// Serious problems that are likely to cause the program to exit.
func Fatalf(format string, v ...interface{}) {
	if !CheckLevel(LevelFatal) {
		return
	}
	std.output(2, LevelFatal, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Serious problems that are likely to cause the program to exit.
func Fatal(v ...interface{}) {
	if !CheckLevel(LevelFatal) {
		return
	}
	std.output(2, LevelFatal, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Issues that shouldn't be ignored.
func Errorf(format string, v ...interface{}) {
	if !CheckLevel(LevelError) {
		return
	}
	std.output(2, LevelError, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Issues that shouldn't be ignored.
func Error(v ...interface{}) {
	if !CheckLevel(LevelError) {
		return
	}
	std.output(2, LevelError, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Issues that may be ignored.
func Warningf(format string, v ...interface{}) {
	if !CheckLevel(LevelWarning) {
		return
	}
	std.output(2, LevelWarning, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Issues that may be ignored.
func Warning(v ...interface{}) {
	if !CheckLevel(LevelWarning) {
		return
	}
	std.output(2, LevelWarning, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Important status changes. The prefix is 'I'.
func Noticef(format string, v ...interface{}) {
	if !CheckLevel(LevelNotice) {
		return
	}
	std.output(2, LevelNotice, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Important status changes. The prefix is 'I'.
func Notice(v ...interface{}) {
	if !CheckLevel(LevelNotice) {
		return
	}
	std.output(2, LevelNotice, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Normal work reports.
func Infof(format string, v ...interface{}) {
	if !CheckLevel(LevelInfo) {
		return
	}
	std.output(2, LevelInfo, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Normal work reports.
func Info(v ...interface{}) {
	if !CheckLevel(LevelInfo) {
		return
	}
	std.output(2, LevelInfo, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Extra information for debugging.
func Debugf(format string, v ...interface{}) {
	if !CheckLevel(LevelDebug) {
		return
	}
	std.output(2, LevelDebug, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Extra information for debugging.
func Debug(v ...interface{}) {
	if !CheckLevel(LevelDebug) {
		return
	}
	std.output(2, LevelDebug, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// Details for inspecting specific issues.
func Verbosef(format string, v ...interface{}) {
	if !CheckLevel(LevelVerbose) {
		return
	}
	std.output(2, LevelVerbose, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// Details for inspecting specific issues.
func Verbose(v ...interface{}) {
	if !CheckLevel(LevelVerbose) {
		return
	}
	std.output(2, LevelVerbose, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}

// More details that may significantly impact performance. The prefix is 'V'.
func VeryVerbosef(format string, v ...interface{}) {
	if !CheckLevel(LevelVeryVerbose) {
		return
	}
	std.output(2, LevelVeryVerbose, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	}, nil)
}

// More details that may significantly impact performance. The prefix is 'V'.
func VeryVerbose(v ...interface{}) {
	if !CheckLevel(LevelVeryVerbose) {
		return
	}
	std.output(2, LevelVeryVerbose, func(b []byte) []byte {
		return fmt.Append(b, v...)
	}, nil)
}
