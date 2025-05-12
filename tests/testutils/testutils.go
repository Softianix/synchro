// Helpers / utils for testing purposes.
package testutils

import (
	"io"
	"log"
	"log/slog"
	"os"
)

// CaptureOutput runs fn function, capturing anything it writes to:
// stdout, stderr, the standard log.Logger, and the default slog.Logger.
// It returns (stdout, stderr, logs) strings.
func CaptureOutput(fn func()) (string, string, string) {
	// Backup real outputs & loggers
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	oldLogWriter := log.Writer()
	oldSlog := slog.Default()

	// Create pipes
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	rLog, wLog, _ := os.Pipe()

	// Redirect os.Stdout / os.Stderr
	os.Stdout = wOut
	os.Stderr = wErr

	// Redirect standard log to wLog
	log.SetOutput(wLog)

	// Swap default slog.Logger to write into wLog
	handler := slog.NewTextHandler(wLog, nil)
	slog.SetDefault(slog.New(handler))

	// Run the user function
	fn()

	// Restore everything
	wOut.Close()
	wErr.Close()
	wLog.Close()

	os.Stdout = oldStdout
	os.Stderr = oldStderr
	log.SetOutput(oldLogWriter)
	slog.SetDefault(oldSlog)

	// Read pipes
	outBytes, _ := io.ReadAll(rOut)
	errBytes, _ := io.ReadAll(rErr)
	logBytes, _ := io.ReadAll(rLog)

	// Return captured output in the form of strings
	return string(outBytes), string(errBytes), string(logBytes)
}
