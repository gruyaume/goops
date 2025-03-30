package goops

import (
	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/goops/environment"
)

type HookContext struct {
	Commands    *commands.Command
	Environment *environment.Environment
}

func NewHookContext() *HookContext {
	return &HookContext{
		Commands:    &commands.Command{},
		Environment: &environment.Environment{},
	}
}
