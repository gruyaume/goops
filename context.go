package goops

import (
	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/goops/environment"
)

type HookContext struct {
	Commands    *commands.Command
	Environment *environment.Environment
}

func NewHookContext() HookContext {
	hookCommand := &commands.Command{}
	environment := &environment.Environment{}
	return HookContext{
		Commands:    hookCommand,
		Environment: environment,
	}
}
