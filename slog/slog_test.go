// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package slog_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
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

func TestNewLogger(t *testing.T) {
	logger := slog.NewLogger()
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}
	// Default level should be LevelSilence (0)
	if logger.Level() != slog.LevelSilence {
		t.Errorf("expected default level %d, got %d", slog.LevelSilence, logger.Level())
	}
}

func TestDefault(t *testing.T) {
	logger := slog.Default()
	if logger == nil {
		t.Fatal("Default returned nil")
	}
}

func TestSetLevel(t *testing.T) {
	logger := slog.NewLogger()
	levels := []slog.Level{
		slog.LevelSilence,
		slog.LevelFatal,
		slog.LevelError,
		slog.LevelWarning,
		slog.LevelNotice,
		slog.LevelInfo,
		slog.LevelDebug,
		slog.LevelVerbose,
		slog.LevelVeryVerbose,
	}
	for _, level := range levels {
		logger.SetLevel(level)
		if logger.Level() != level {
			t.Errorf("expected level %d, got %d", level, logger.Level())
		}
	}
}

func TestSetOutput(t *testing.T) {
	logger := slog.NewLogger()
	var buf bytes.Buffer

	// Test OutputDiscard
	logger.SetOutput(slog.OutputDiscard)
	logger.SetLevel(slog.LevelVeryVerbose)
	logger.Debug("test discard")

	// Test OutputWriter
	logger.SetOutput(slog.OutputWriter, &buf)
	logger.Debug("test writer")
	if buf.Len() == 0 {
		t.Error("expected output to buffer, got nothing")
	}
}

func TestLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.NewLogger()
	logger.SetOutput(slog.OutputWriter, &buf)

	tests := []struct {
		setLevel    slog.Level
		logLevel    slog.Level
		shouldLog   bool
		description string
	}{
		{slog.LevelError, slog.LevelFatal, true, "Fatal should log when level is Error"},
		{slog.LevelError, slog.LevelError, true, "Error should log when level is Error"},
		{slog.LevelError, slog.LevelWarning, false, "Warning should not log when level is Error"},
		{slog.LevelWarning, slog.LevelWarning, true, "Warning should log when level is Warning"},
		{slog.LevelSilence, slog.LevelFatal, false, "Fatal should not log when level is Silence"},
		{slog.LevelVeryVerbose, slog.LevelVeryVerbose, true, "VeryVerbose should log when level is VeryVerbose"},
	}

	for _, tt := range tests {
		buf.Reset()
		logger.SetLevel(tt.setLevel)

		switch tt.logLevel {
		case slog.LevelFatal:
			logger.Fatal("test")
		case slog.LevelError:
			logger.Error("test")
		case slog.LevelWarning:
			logger.Warning("test")
		case slog.LevelNotice:
			logger.Notice("test")
		case slog.LevelInfo:
			logger.Info("test")
		case slog.LevelDebug:
			logger.Debug("test")
		case slog.LevelVerbose:
			logger.Verbose("test")
		case slog.LevelVeryVerbose:
			logger.VeryVerbose("test")
		}

		hasOutput := buf.Len() > 0
		if hasOutput != tt.shouldLog {
			t.Errorf("%s: expected shouldLog=%v, got hasOutput=%v", tt.description, tt.shouldLog, hasOutput)
		}
	}
}

func TestCheckLevel(t *testing.T) {
	logger := slog.NewLogger()
	var buf bytes.Buffer
	logger.SetOutput(slog.OutputWriter, &buf)

	logger.SetLevel(slog.LevelWarning)

	// These should be true (level <= current level)
	// Note: CheckLevel is for the default logger, so we test via output behavior
	if slog.LevelFatal > logger.Level() {
		t.Error("Fatal level should be enabled when level is Warning")
	}
	if slog.LevelError > logger.Level() {
		t.Error("Error level should be enabled when level is Warning")
	}
	if slog.LevelWarning > logger.Level() {
		t.Error("Warning level should be enabled when level is Warning")
	}

	// These should be false
	if slog.LevelInfo <= logger.Level() {
		t.Error("Info level should be disabled when level is Warning")
	}
}

func TestLoggerOutputFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.NewLogger()
	logger.SetOutput(slog.OutputWriter, &buf)
	logger.SetLevel(slog.LevelDebug)

	logger.Debug("hello world")

	output := buf.String()
	// Check format: "D YYYY-MM-DDTHH:MM:SS+ZZ:ZZ file:line message\n"
	if !strings.HasPrefix(output, "D ") {
		t.Errorf("expected output to start with 'D ', got: %s", output)
	}
	if !strings.Contains(output, "hello world") {
		t.Errorf("expected output to contain 'hello world', got: %s", output)
	}
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("expected output to end with newline, got: %s", output)
	}
}

func TestLoggerFormatFunctions(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.NewLogger()
	logger.SetOutput(slog.OutputWriter, &buf)
	logger.SetLevel(slog.LevelVeryVerbose)

	tests := []struct {
		name     string
		logFunc  func()
		expected string
	}{
		{"Fatalf", func() { logger.Fatalf("test %d", 42) }, "test 42"},
		{"Fatal", func() { logger.Fatal("test", 42) }, "test 42"},
		{"Errorf", func() { logger.Errorf("error %s", "msg") }, "error msg"},
		{"Error", func() { logger.Error("error", "msg") }, "error msg"},
		{"Warningf", func() { logger.Warningf("warn %v", true) }, "warn true"},
		{"Warning", func() { logger.Warning("warn", true) }, "warn true"},
		{"Noticef", func() { logger.Noticef("notice %d", 1) }, "notice 1"},
		{"Notice", func() { logger.Notice("notice", 1) }, "notice 1"},
		{"Infof", func() { logger.Infof("info %s", "test") }, "info test"},
		{"Info", func() { logger.Info("info", "test") }, "info test"},
		{"Debugf", func() { logger.Debugf("debug %d", 99) }, "debug 99"},
		{"Debug", func() { logger.Debug("debug", 99) }, "debug 99"},
		{"Verbosef", func() { logger.Verbosef("verbose %s", "x") }, "verbose x"},
		{"Verbose", func() { logger.Verbose("verbose", "x") }, "verbose x"},
		{"VeryVerbosef", func() { logger.VeryVerbosef("vv %d", 0) }, "vv 0"},
		{"VeryVerbose", func() { logger.VeryVerbose("vv", 0) }, "vv 0"},
	}

	for _, tt := range tests {
		buf.Reset()
		tt.logFunc()
		if !strings.Contains(buf.String(), tt.expected) {
			t.Errorf("%s: expected output to contain '%s', got: %s", tt.name, tt.expected, buf.String())
		}
	}
}

func TestSetFilePrefix(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.NewLogger()
	logger.SetOutput(slog.OutputWriter, &buf)
	logger.SetLevel(slog.LevelDebug)

	// Without prefix, full path is shown
	buf.Reset()
	logger.Debug("test1")
	output1 := buf.String()

	// With prefix, path should be trimmed
	logger.SetFilePrefix("/home/hex/source/gosnippets/")
	buf.Reset()
	logger.Debug("test2")
	output2 := buf.String()

	// output2 should have shorter file path (or at least different)
	// Both should contain the message
	if !strings.Contains(output1, "test1") {
		t.Errorf("expected output1 to contain 'test1', got: %s", output1)
	}
	if !strings.Contains(output2, "test2") {
		t.Errorf("expected output2 to contain 'test2', got: %s", output2)
	}
}

func TestOutputWithExtra(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.NewLogger()
	logger.SetOutput(slog.OutputWriter, &buf)
	logger.SetLevel(slog.LevelDebug)

	extraData := "extra data line\n"
	logger.Println(0, slog.LevelDebug, func(w io.Writer) error {
		_, err := w.Write([]byte(extraData))
		return err
	}, "main message")

	output := buf.String()
	if !strings.Contains(output, "main message") {
		t.Errorf("expected output to contain 'main message', got: %s", output)
	}
	if !strings.Contains(output, extraData) {
		t.Errorf("expected output to contain extra data, got: %s", output)
	}
}

func TestConcurrentLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.NewLogger()
	logger.SetOutput(slog.OutputWriter, &buf)
	logger.SetLevel(slog.LevelDebug)

	var wg sync.WaitGroup
	numGoroutines := 10
	numLogs := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numLogs; j++ {
				logger.Debugf("goroutine %d log %d", id, j)
			}
		}(i)
	}

	wg.Wait()

	// Count newlines to verify all logs were written
	lineCount := strings.Count(buf.String(), "\n")
	expected := numGoroutines * numLogs
	if lineCount != expected {
		t.Errorf("expected %d log lines, got %d", expected, lineCount)
	}
}

func TestWrap(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.NewLogger()
	logger.SetOutput(slog.OutputWriter, &buf)
	logger.SetLevel(slog.LevelDebug)

	stdLogger := slog.Wrap(logger, slog.LevelDebug)
	if stdLogger == nil {
		t.Fatal("Wrap returned nil")
	}

	stdLogger.Print("wrapped message")

	output := buf.String()
	if !strings.Contains(output, "wrapped message") {
		t.Errorf("expected output to contain 'wrapped message', got: %s", output)
	}
}

func TestPackageLevelFunctions(t *testing.T) {
	// Save current state
	origLevel := slog.Default().Level()
	defer slog.Default().SetLevel(origLevel)

	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelVeryVerbose)

	tests := []struct {
		name    string
		logFunc func()
		prefix  string
	}{
		{"Fatal", func() { slog.Fatal("fatal msg") }, "F "},
		{"Fatalf", func() { slog.Fatalf("fatal %s", "msg") }, "F "},
		{"Error", func() { slog.Error("error msg") }, "E "},
		{"Errorf", func() { slog.Errorf("error %s", "msg") }, "E "},
		{"Warning", func() { slog.Warning("warn msg") }, "W "},
		{"Warningf", func() { slog.Warningf("warn %s", "msg") }, "W "},
		{"Notice", func() { slog.Notice("notice msg") }, "I "},
		{"Noticef", func() { slog.Noticef("notice %s", "msg") }, "I "},
		{"Info", func() { slog.Info("info msg") }, "I "},
		{"Infof", func() { slog.Infof("info %s", "msg") }, "I "},
		{"Debug", func() { slog.Debug("debug msg") }, "D "},
		{"Debugf", func() { slog.Debugf("debug %s", "msg") }, "D "},
		{"Verbose", func() { slog.Verbose("verbose msg") }, "V "},
		{"Verbosef", func() { slog.Verbosef("verbose %s", "msg") }, "V "},
		{"VeryVerbose", func() { slog.VeryVerbose("vv msg") }, "V "},
		{"VeryVerbosef", func() { slog.VeryVerbosef("vv %s", "msg") }, "V "},
	}

	for _, tt := range tests {
		buf.Reset()
		tt.logFunc()
		output := buf.String()
		if !strings.HasPrefix(output, tt.prefix) {
			t.Errorf("%s: expected prefix '%s', got: %s", tt.name, tt.prefix, output)
		}
	}
}

func TestCheckLevelPackageFunc(t *testing.T) {
	origLevel := slog.Default().Level()
	defer slog.Default().SetLevel(origLevel)

	slog.Default().SetLevel(slog.LevelWarning)

	if !slog.CheckLevel(slog.LevelFatal) {
		t.Error("CheckLevel(Fatal) should be true when level is Warning")
	}
	if !slog.CheckLevel(slog.LevelWarning) {
		t.Error("CheckLevel(Warning) should be true when level is Warning")
	}
	if slog.CheckLevel(slog.LevelInfo) {
		t.Error("CheckLevel(Info) should be false when level is Warning")
	}
}

func TestText(t *testing.T) {
	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelDebug)

	testText := "line1\nline2\nline3"
	slog.Text(slog.LevelDebug, testText, "text message")

	output := buf.String()
	if !strings.Contains(output, "text message") {
		t.Errorf("expected output to contain 'text message', got: %s", output)
	}
	if !strings.Contains(output, "line1") {
		t.Errorf("expected output to contain 'line1', got: %s", output)
	}
}

func TestTextf(t *testing.T) {
	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelDebug)

	testText := "hello\nworld"
	slog.Textf(slog.LevelDebug, testText, -1, "formatted %s", "msg")

	output := buf.String()
	if !strings.Contains(output, "formatted msg") {
		t.Errorf("expected output to contain 'formatted msg', got: %s", output)
	}
}

func TestBinary(t *testing.T) {
	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelDebug)

	testData := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello"
	slog.Binary(slog.LevelDebug, testData, "binary dump")

	output := buf.String()
	if !strings.Contains(output, "binary dump") {
		t.Errorf("expected output to contain 'binary dump', got: %s", output)
	}
	// Should contain hex dump
	if !strings.Contains(output, "48") {
		t.Errorf("expected output to contain hex '48', got: %s", output)
	}
}

func TestBinaryf(t *testing.T) {
	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelDebug)

	testData := []byte{0x41, 0x42, 0x43} // "ABC"
	slog.Binaryf(slog.LevelDebug, testData, -1, "binary %d", 123)

	output := buf.String()
	if !strings.Contains(output, "binary 123") {
		t.Errorf("expected output to contain 'binary 123', got: %s", output)
	}
}

func TestStack(t *testing.T) {
	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelDebug)

	slog.Stack(slog.LevelDebug, 0, "stack trace")

	output := buf.String()
	if !strings.Contains(output, "stack trace") {
		t.Errorf("expected output to contain 'stack trace', got: %s", output)
	}
	// Should contain function name
	if !strings.Contains(output, "TestStack") {
		t.Errorf("expected output to contain 'TestStack', got: %s", output)
	}
}

func TestStackf(t *testing.T) {
	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelDebug)

	slog.Stackf(slog.LevelDebug, 0, "stack %s", "formatted")

	output := buf.String()
	if !strings.Contains(output, "stack formatted") {
		t.Errorf("expected output to contain 'stack formatted', got: %s", output)
	}
}

func TestCheck(t *testing.T) {
	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelDebug)

	// Should not panic
	slog.Check(true, "this should not panic")

	// Should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Check(false) should panic")
		}
	}()
	slog.Check(false, "this should panic")
}

func TestCheckf(t *testing.T) {
	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelDebug)

	// Should not panic
	slog.Checkf(true, "should not panic %d", 42)

	// Should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Checkf(false) should panic")
		}
	}()
	slog.Checkf(false, "should panic %d", 42)
}

func TestOutputf(t *testing.T) {
	var buf bytes.Buffer
	slog.Default().SetOutput(slog.OutputWriter, &buf)
	slog.Default().SetLevel(slog.LevelDebug)

	slog.Printf(0, slog.LevelDebug, nil, "outputf %d %s", 42, "test")

	output := buf.String()
	if !strings.Contains(output, "outputf 42 test") {
		t.Errorf("expected output to contain 'outputf 42 test', got: %s", output)
	}
}

func TestLoggerOutputf(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.NewLogger()
	logger.SetOutput(slog.OutputWriter, &buf)
	logger.SetLevel(slog.LevelDebug)

	logger.Printf(0, slog.LevelDebug, nil, "logger outputf %d", 99)

	output := buf.String()
	if !strings.Contains(output, "logger outputf 99") {
		t.Errorf("expected output to contain 'logger outputf 99', got: %s", output)
	}
}

func BenchmarkOutput(b *testing.B) {
	testData := "The quick brown fox jumps over the lazy dog.🍌 "
	for i := 1; i < b.N; i++ {
		slog.Println(0, slog.LevelFatal, nil, testData)
	}
}

func BenchmarkOutputf(b *testing.B) {
	for i := 1; i < b.N; i++ {
		slog.Printf(0, slog.LevelFatal, nil, "test %d %s", 42, "hello")
	}
}

func BenchmarkLevelFiltered(b *testing.B) {
	// Test performance when log is filtered out
	slog.Default().SetLevel(slog.LevelSilence)
	for i := 1; i < b.N; i++ {
		slog.Debug("this should be filtered")
	}
}

func BenchmarkConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			slog.Println(0, slog.LevelFatal, nil, "concurrent test")
		}
	})
}
