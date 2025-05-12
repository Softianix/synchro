// Colorful logger implementation.
package logging

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

// Level represents the severity of the log message.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// ANSI color codes for different log levels.
var levelColors = map[Level]string{
	LevelDebug: "\x1b[36m", // Cyan
	LevelInfo:  "\x1b[32m", // Green
	LevelWarn:  "\x1b[33m", // Yellow
	LevelError: "\x1b[31m", // Red
}

// Human-readable names for each level.
var levelNames = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO ",
	LevelWarn:  "WARN ",
	LevelError: "ERROR",
}

// resetColor resets terminal color.
const resetColor = "\x1b[0m"

// Logger is a simple colorful logger.
type Logger struct{}

// New creates a new Logger instance.
func NewLogger() *Logger {
	return &Logger{}
}

// defaultLogger is the package-level singleton instance.
var defaultLogger = NewLogger()

// GetLogger returns the singleton Logger instance.
func GetLogger() *Logger {
	return defaultLogger
}

// log prints a formatted log message with color, timestamp, file, and line number.
func (l *Logger) log(level Level, format string, args ...interface{}) {
	// Timestamp in ISO-like format
	ts := time.Now().Format("2006-01-02 15:04:05.000")

	// Retrieve caller info
	// Skip 3 frames: runtime.Caller -> Logger.log -> Info/Warn/etc.
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "?"
		line = 0
	}
	file = filepath.Base(file)

	// Prepare the message
	msg := fmt.Sprintf(format, args...)

	// Assemble and print the log line
	color := levelColors[level]
	name := levelNames[level]
	fmt.Printf("%s%s%s | %s | %s:%d | %s\n", color, name, resetColor, ts, file, line, msg)
}

// Debug logs a message at DEBUG level.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info logs a message at INFO level.
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn logs a message at WARN level.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error logs a message at ERROR level.
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}
