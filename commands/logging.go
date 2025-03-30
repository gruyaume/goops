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

func (command Command) JujuLog(logLevel Level, message string, extraArgs ...string) {
	args := []string{"--log-level=" + levelStrings[logLevel], message}
	args = append(args, extraArgs...)
	_, err := command.Runner.Run(JujuLogCommand, args...)
	if err != nil {
		log.Println("failed to run juju-log command:", err)
	}
}
