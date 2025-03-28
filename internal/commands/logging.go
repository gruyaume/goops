package commands

import (
	"log"
)

// JujuLogCommand is the external command used to log messages.
const JujuLogCommand = "juju-log"

// Level represents the log severity.
type Level int

const (
	Debug Level = iota
	Info
	Warning
	Error
)

var levelStrings = []string{
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
}

// JujuLog sends a log message to juju-log using the given runner,
// appending any extra arguments provided.
func JujuLog(runner CommandRunner, message string, logLevel Level, extraArgs ...string) {
	args := []string{"--log-level", levelStrings[logLevel], message}
	args = append(args, extraArgs...)
	_, err := runner.Run(JujuLogCommand, args...)
	if err != nil {
		log.Println("failed to run juju-log command:", err)
	}
}

// Logger wraps a CommandRunner to provide logging methods.
type Logger struct {
	runner CommandRunner
}

// NewLogger creates a new Logger instance with the given CommandRunner.
func NewLogger(runner CommandRunner) *Logger {
	return &Logger{runner: runner}
}

// Debug logs a message at the DEBUG level with extra arguments.
func (l *Logger) Debug(message string, extraArgs ...string) {
	JujuLog(l.runner, message, Debug, extraArgs...)
}

// Info logs a message at the INFO level with extra arguments.
func (l *Logger) Info(message string, extraArgs ...string) {
	JujuLog(l.runner, message, Info, extraArgs...)
}

// Warning logs a message at the WARNING level with extra arguments.
func (l *Logger) Warning(message string, extraArgs ...string) {
	JujuLog(l.runner, message, Warning, extraArgs...)
}

// Error logs a message at the ERROR level with extra arguments.
func (l *Logger) Error(message string, extraArgs ...string) {
	JujuLog(l.runner, message, Error, extraArgs...)
}
