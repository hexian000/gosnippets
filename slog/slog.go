// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package slog

import (
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

// Default returns the default logger.
func Default() *Logger {
	return std
}

// Outputf is the low-level interface to write arbitary log messages.
func Outputf(calldepth int, level Level, extra func(io.Writer) error, format string, v ...interface{}) error {
	return std.output(calldepth+1, level, func(b []byte) []byte {
		return AppendMsgf(b, format, v...)
	}, extra)
}

// Output is the low-level interface to write arbitary log messages.
func Output(calldepth int, level Level, extra func(io.Writer) error, v ...interface{}) error {
	return std.output(calldepth+1, level, func(b []byte) []byte {
		return AppendMsg(b, v...)
	}, extra)
}

// CheckLevel checks whether the given level is enabled.
func CheckLevel(level Level) bool {
	return level <= std.Level()
}

// Fatalf logs serious problems that are likely to cause the program to exit.
func Fatalf(format string, v ...any) {
	if !CheckLevel(LevelFatal) {
		return
	}
	std.output(1, LevelFatal, func(b []byte) []byte {
		return AppendMsgf(b, format, v...)
	}, nil)
}

// Fatal logs serious problems that are likely to cause the program to exit.
func Fatal(v ...any) {
	if !CheckLevel(LevelFatal) {
		return
	}
	std.output(1, LevelFatal, func(b []byte) []byte {
		return AppendMsg(b, v...)
	}, nil)
}

// Errorf logs issues that shouldn't be ignored.
func Errorf(format string, v ...any) {
	if !CheckLevel(LevelError) {
		return
	}
	std.output(1, LevelError, func(b []byte) []byte {
		return AppendMsgf(b, format, v...)
	}, nil)
}

// Error logs issues that shouldn't be ignored.
func Error(v ...any) {
	if !CheckLevel(LevelError) {
		return
	}
	std.output(1, LevelError, func(b []byte) []byte {
		return AppendMsg(b, v...)
	}, nil)
}

// Warningf logs issues that may be ignored.
func Warningf(format string, v ...any) {
	if !CheckLevel(LevelWarning) {
		return
	}
	std.output(1, LevelWarning, func(b []byte) []byte {
		return AppendMsgf(b, format, v...)
	}, nil)
}

// Warning logs issues that may be ignored.
func Warning(v ...any) {
	if !CheckLevel(LevelWarning) {
		return
	}
	std.output(1, LevelWarning, func(b []byte) []byte {
		return AppendMsg(b, v...)
	}, nil)
}

// Noticef logs important status changes. The prefix is 'I'.
func Noticef(format string, v ...any) {
	if !CheckLevel(LevelNotice) {
		return
	}
	std.output(1, LevelNotice, func(b []byte) []byte {
		return AppendMsgf(b, format, v...)
	}, nil)
}

// Notice logs important status changes. The prefix is 'I'.
func Notice(v ...any) {
	if !CheckLevel(LevelNotice) {
		return
	}
	std.output(1, LevelNotice, func(b []byte) []byte {
		return AppendMsg(b, v...)
	}, nil)
}

// Infof logs normal work reports.
func Infof(format string, v ...any) {
	if !CheckLevel(LevelInfo) {
		return
	}
	std.output(1, LevelInfo, func(b []byte) []byte {
		return AppendMsgf(b, format, v...)
	}, nil)
}

// Info logs normal work reports.
func Info(v ...any) {
	if !CheckLevel(LevelInfo) {
		return
	}
	std.output(1, LevelInfo, func(b []byte) []byte {
		return AppendMsg(b, v...)
	}, nil)
}

// Debugf logs extra information for debugging.
func Debugf(format string, v ...any) {
	if !CheckLevel(LevelDebug) {
		return
	}
	std.output(1, LevelDebug, func(b []byte) []byte {
		return AppendMsgf(b, format, v...)
	}, nil)
}

// Debug logs extra information for debugging.
func Debug(v ...any) {
	if !CheckLevel(LevelDebug) {
		return
	}
	std.output(1, LevelDebug, func(b []byte) []byte {
		return AppendMsg(b, v...)
	}, nil)
}

// Verbosef logs details for inspecting specific issues.
func Verbosef(format string, v ...any) {
	if !CheckLevel(LevelVerbose) {
		return
	}
	std.output(1, LevelVerbose, func(b []byte) []byte {
		return AppendMsgf(b, format, v...)
	}, nil)
}

// Verbose logs details for inspecting specific issues.
func Verbose(v ...any) {
	if !CheckLevel(LevelVerbose) {
		return
	}
	std.output(1, LevelVerbose, func(b []byte) []byte {
		return AppendMsg(b, v...)
	}, nil)
}

// VeryVerbosef logs more details that may significantly impact performance. The prefix is 'V'.
func VeryVerbosef(format string, v ...any) {
	if !CheckLevel(LevelVeryVerbose) {
		return
	}
	std.output(1, LevelVeryVerbose, func(b []byte) []byte {
		return AppendMsgf(b, format, v...)
	}, nil)
}

// VeryVerbose logs more details that may significantly impact performance. The prefix is 'V'.
func VeryVerbose(v ...any) {
	if !CheckLevel(LevelVeryVerbose) {
		return
	}
	std.output(1, LevelVeryVerbose, func(b []byte) []byte {
		return AppendMsg(b, v...)
	}, nil)
}
