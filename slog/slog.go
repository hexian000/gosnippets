package slog

import (
	"fmt"
)

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

type Level int

var levelChar = [...]byte{
	'-', 'F', 'E', 'W', 'I', 'I', 'D', 'V', 'V',
}

var std = NewLogger(LevelVerbose)

func Default() *Logger {
	return std
}

func Write(calldepth int, level Level, msg []byte) {
	std.Output(calldepth+1, level, msg)
}

func (l *Logger) CheckLevel(level Level) bool {
	return level <= std.Level()
}

// Serious problems that are likely to cause the program to exit.
func Fatalf(format string, v ...interface{}) {
	std.Output(2, LevelFatal, []byte(fmt.Sprintf(format, v...)))
}

// Serious problems that are likely to cause the program to exit.
func Fatal(v ...interface{}) {
	std.Output(2, LevelFatal, []byte(fmt.Sprint(v...)))
}

// Issues that shouldn't be ignored.
func Errorf(format string, v ...interface{}) {
	std.Output(2, LevelError, []byte(fmt.Sprintf(format, v...)))
}

// Issues that shouldn't be ignored.
func Error(v ...interface{}) {
	std.Output(2, LevelError, []byte(fmt.Sprint(v...)))
}

// Issues that may be ignored.
func Warningf(format string, v ...interface{}) {
	std.Output(2, LevelWarning, []byte(fmt.Sprintf(format, v...)))
}

// Issues that may be ignored.
func Warning(v ...interface{}) {
	std.Output(2, LevelWarning, []byte(fmt.Sprint(v...)))
}

// Important status changes. The prefix is 'I'.
func Noticef(format string, v ...interface{}) {
	std.Output(2, LevelNotice, []byte(fmt.Sprintf(format, v...)))
}

// Important status changes. The prefix is 'I'.
func Notice(v ...interface{}) {
	std.Output(2, LevelNotice, []byte(fmt.Sprint(v...)))
}

// Normal work reports.
func Infof(format string, v ...interface{}) {
	std.Output(2, LevelInfo, []byte(fmt.Sprintf(format, v...)))
}

// Normal work reports.
func Info(v ...interface{}) {
	std.Output(2, LevelInfo, []byte(fmt.Sprint(v...)))
}

// Extra information for debugging.
func Debugf(format string, v ...interface{}) {
	std.Output(2, LevelDebug, []byte(fmt.Sprintf(format, v...)))
}

// Extra information for debugging.
func Debug(v ...interface{}) {
	std.Output(2, LevelDebug, []byte(fmt.Sprint(v...)))
}

// Details for inspecting specific issues.
func Verbosef(format string, v ...interface{}) {
	std.Output(2, LevelVerbose, []byte(fmt.Sprintf(format, v...)))
}

// Details for inspecting specific issues.
func Verbose(v ...interface{}) {
	std.Output(2, LevelVerbose, []byte(fmt.Sprint(v...)))
}

// More details that may significantly impact performance. The prefix is 'V'.
func VeryVerbosef(format string, v ...interface{}) {
	std.Output(2, LevelVeryVerbose, []byte(fmt.Sprintf(format, v...)))
}

// More details that may significantly impact performance. The prefix is 'V'.
func VeryVerbose(v ...interface{}) {
	std.Output(2, LevelVeryVerbose, []byte(fmt.Sprint(v...)))
}
