package main

import (
	"os"
	"strings"

	"github.com/hexian000/gosnippets/slog"
)

func init() {
	std := slog.Default()
	std.SetOutput(slog.OutputTerminal, os.Stdout)
	std.SetLevel(slog.LevelVerbose)
	if dir, err := os.Getwd(); err == nil {
		std.SetFilePrefix(dir + "/")
	}
}

func testStack() {
	slog.Stackf(slog.LevelDebug, 0, "wa")
}

func testStackInlined() {
	testStack()
}

func main() {
	slog.Debug("begin")

	s := strings.Repeat("The quick brown fox jumps over the lazy dog. 🍌 😂 🐇 ☕ 🎈", 8) + `
1	2	3	4	5	6	7	8	9	10	11	12	13	14	15	16	17	18	19	20	21	22	23	
1	slog.Default().SetOutputConfig("stdout", "")
1	2	slog.Default().SetFilePrefix("")
1	2	3	slog.Default().SetLevel(slog.LevelNotice)
		slog.Default().SetLevel(slog.LevelVerbose)
	slog.Debug("wa")` + "\r\n\r\n"
	slog.Textf(slog.LevelDebug, s, "wa")
	slog.Binaryf(slog.LevelDebug, []byte(s), "wa")
	testStackInlined()
	slog.Output(0, slog.LevelDebug, nil, "test1")
	slog.Default().Output(0, slog.LevelDebug, nil, "test2")
	slog.Debug("end")
}
