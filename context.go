package goops

import (
	"github.com/gruyaume/goops/commands"
)

type HookContext struct {
	Commands *commands.Command
}

func NewHookContext() *HookContext {
	return &HookContext{
		Commands: commands.NewCommand(),
	}
}
