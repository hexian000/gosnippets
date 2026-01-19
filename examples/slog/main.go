package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hexian000/gosnippets/slog"
)

func init() {
	std := slog.Default()
	std.SetOutput(slog.OutputTerminal, os.Stdout)
	std.SetLevel(slog.LevelVeryVerbose)
	if _, file, _, ok := runtime.Caller(0); ok {
		std.SetFilePrefix(filepath.Dir(file) + "/")
	}
	slog.Debug("slog initialized")
}

func testStack() {
	slog.Stackf(slog.LevelDebug, 0, "wa")
}

func testStackInlined() {
	testStack()
}

func main() {
	slog.VeryVerbose("VeryVerbose: More details that may significantly impact performance.")
	slog.Verbose("Verbose: Details for inspecting specific issues.")
	slog.Debug("Debug: Extra information for debugging.")
	slog.Info("Info: Normal work reports.")
	slog.Notice("Notice: Important status changes.")
	slog.Warning("Warning: Issues that may be ignored.")
	slog.Error("Error: Issues that shouldn't be ignored.")
	slog.Fatal("Fatal: Serious problems that are likely to cause the program to exit.")
	slog.Temporary("Temporary: Temporary logs are printed in any debug build regardless of log level.")

	s := strings.Repeat("The quick brown fox jumps over the lazy dog. 🍌 😂 🐇 ☕ 🎈", 8) + `
1	2	3	4	5	6	7	8	9	10	11	12	13	14	15	16	17	18	19	20	21	22	23	
1	slog.Default().SetOutputConfig("stdout", "")
1	2	slog.Default().SetFilePrefix("")
1	2	3	slog.Default().SetLevel(slog.LevelNotice)
		slog.Default().SetLevel(slog.LevelVerbose)
	slog.Debug("wa")` + "\r\n\r\n"
	slog.Textf(slog.LevelDebug, s, "wa")
	slog.Binaryf(slog.LevelWarning, []byte(s), "wa")
	testStackInlined()
	slog.Output(0, slog.LevelDebug, nil, "test1")
	slog.Default().Output(0, slog.LevelDebug, nil, "test2")
	slog.Debug("end")
}
