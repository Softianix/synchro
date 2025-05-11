package main

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

// captureOutput runs fn(), capturing anything it writes to stdout, stderr,
// the standard log.Logger, and the default slog.Logger. It returns
// (stdout, stderr, logs).
func captureOutput(fn func()) (string, string, string) {
	// 1) Backup real outputs & loggers
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	oldLogWriter := log.Writer()
	oldSlog := slog.Default()

	// 2) Create pipes
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	rLog, wLog, _ := os.Pipe()

	// 3) Redirect os.Stdout / os.Stderr
	os.Stdout = wOut
	os.Stderr = wErr

	// 4) Redirect standard log to wLog
	log.SetOutput(wLog)

	// 5) Swap default slog.Logger to write into wLog
	handler := slog.NewTextHandler(wLog, nil)
	slog.SetDefault(slog.New(handler))

	// 6) Run the user function
	fn()

	// 7) Restore everything
	wOut.Close()
	wErr.Close()
	wLog.Close()

	os.Stdout = oldStdout
	os.Stderr = oldStderr
	log.SetOutput(oldLogWriter)
	slog.SetDefault(oldSlog)

	// 8) Read pipes
	outBytes, _ := io.ReadAll(rOut)
	errBytes, _ := io.ReadAll(rErr)
	logBytes, _ := io.ReadAll(rLog)

	return string(outBytes), string(errBytes), string(logBytes)
}

func TestE2ESimpleCopying(t *testing.T) {
	// --- ARRANGE ---
	// Create a temporary source and destination directories
	sourceDir, err := os.MkdirTemp("", "sourceDir")
	if err != nil {
		t.Fatalf("failed to create source directory: %v", err)
	}

	destinationDir, err := os.MkdirTemp("", "destinationDir")
	if err != nil {
		t.Fatalf("failed to create destination directory: %v", err)
	}

	// Create a test file in the source directory
	fileName := "testfile.txt"
	filePath := filepath.Join(sourceDir, fileName)
	fileContent := []byte("Hello, World!")
	err = os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	expectedStdOut := ""
	expectedStdErr := ""
	expectedLogs := `{"level":"INFO","msg":"synchro started","subcommand":"synca"}`

	// --- ACT ---
	stdOut, stdErr, logs := captureOutput(func() {
		realMain([]string{"synchro", "synca", "--src", sourceDir, "--dst", destinationDir})
	})

	// --- ASSERT ---
	// Verify the standard output and error streams
	if stdOut != expectedStdOut {
		t.Fatalf("expected stdout: %s, got: %s", expectedStdOut, stdOut)
	}
	if stdErr != expectedStdErr {
		t.Fatalf("expected stderr: %s, got: %s", expectedStdErr, stdErr)
	}
	// Verify the logs
	if logs != expectedLogs {
		t.Fatalf("expected logs: %s, got: %s", expectedLogs, logs)
	}
	// Verify the file was copied
	copiedFilePath := filepath.Join(destinationDir, fileName)
	if _, err := os.Stat(copiedFilePath); os.IsNotExist(err) {
		t.Fatalf("expected file %s to exist, but it does not", copiedFilePath)
	}
	// Verify the content of the copied file
	copiedContent, err := os.ReadFile(copiedFilePath)
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(copiedContent) != string(fileContent) {
		t.Errorf("expected copied file content: %s, got: %s", string(fileContent), string(copiedContent))
	}

	// --- CLEANUP ---
	err = os.RemoveAll(sourceDir)
	if err != nil {
		t.Fatalf("failed to remove source directory: %v", err)
	}
	err = os.RemoveAll(destinationDir)
	if err != nil {
		t.Fatalf("failed to remove destination directory: %v", err)
	}
}
