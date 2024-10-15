package slog_test

import (
	"os"
	"testing"

	"github.com/hexian000/gosnippets/slog"
)

func init() {
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	slog.Default().SetOutput(slog.OutputWriter, devNull)
	slog.Default().SetLevel(slog.LevelVeryVerbose)
}

func BenchmarkOutput(b *testing.B) {
	testData := "The quick brown fox jumps over the lazy dog.🍌 "
	for i := 1; i < b.N; i++ {
		slog.Output(2, slog.LevelFatal, nil, testData)
	}
}
