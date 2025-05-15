package goops

import (
	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/goops/environment"
	"github.com/gruyaume/goops/metadata"
)

type HookContext struct {
	Commands    *commands.Command
	Environment *environment.Environment
	Metadata    *metadata.Metadata
}

func NewHookContext() *HookContext {
	env := environment.NewEnvironment()

	charmDir := env.JujuCharmDir()
	metadataPath := charmDir + "/metadata.yaml"

	return &HookContext{
		Commands:    commands.NewCommand(),
		Environment: env,
		Metadata:    metadata.GetCharmMetadata(metadataPath),
	}
}
