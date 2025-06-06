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
	commandRunner := GetRunner()

	message := fmt.Sprintf(format, args...)

	cmdArgs := []string{"--log-level=" + levelStrings[level], message}

	_, err := commandRunner.Run(jujuLogCommand, cmdArgs...)
	if err != nil {
		log.Println("failed to run juju-log command:", err)
	}
}

func LogDebugf(format string, args ...any) {
	logf(Debug, format, args...)
}

func LogInfof(format string, args ...any) {
	logf(Info, format, args...)
}

func LogWarningf(format string, args ...any) {
	logf(Warning, format, args...)
}

func LogErrorf(format string, args ...any) {
	logf(Error, format, args...)
}
