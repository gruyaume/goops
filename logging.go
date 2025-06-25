package goops

import (
	"fmt"
	"log"
)

const jujuLogCommand = "juju-log"

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

func logf(level Level, format string, args ...any) {
	commandRunner := GetCommandRunner()

	message := fmt.Sprintf(format, args...)

	cmdArgs := []string{"--log-level=" + levelStrings[level], message}

	_, err := commandRunner.Run(jujuLogCommand, cmdArgs...)
	if err != nil {
		log.Println("failed to run juju-log command:", err)
	}
}

// LogDebugf logs a debug message. Log messages can be read using `juju debug-log`.
func LogDebugf(format string, args ...any) {
	logf(Debug, format, args...)
}

// LogInfof logs an informational message. Log messages can be read using `juju debug-log`.
func LogInfof(format string, args ...any) {
	logf(Info, format, args...)
}

// LogWarningf logs a warning message. Log messages can be read using `juju debug-log`.
func LogWarningf(format string, args ...any) {
	logf(Warning, format, args...)
}

// LogErrorf logs an error message. Log messages can be read using `juju debug-log`.
func LogErrorf(format string, args ...any) {
	logf(Error, format, args...)
}
