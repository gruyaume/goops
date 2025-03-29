package commands

import (
	"log"
)

const JujuLogCommand = "juju-log"

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

func JujuLog(runner CommandRunner, message string, logLevel Level, extraArgs ...string) {
	args := []string{"--log-level=" + levelStrings[logLevel], message}
	args = append(args, extraArgs...)
	_, err := runner.Run(JujuLogCommand, args...)
	if err != nil {
		log.Println("failed to run juju-log command:", err)
	}
}

type Logger struct {
	runner CommandRunner
}

func NewLogger(runner CommandRunner) *Logger {
	return &Logger{runner: runner}
}

func (l *Logger) Debug(message string, extraArgs ...string) {
	JujuLog(l.runner, message, Debug, extraArgs...)
}

func (l *Logger) Info(message string, extraArgs ...string) {
	JujuLog(l.runner, message, Info, extraArgs...)
}

func (l *Logger) Warning(message string, extraArgs ...string) {
	JujuLog(l.runner, message, Warning, extraArgs...)
}

func (l *Logger) Error(message string, extraArgs ...string) {
	JujuLog(l.runner, message, Error, extraArgs...)
}
