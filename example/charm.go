package example

import (
	"os"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/commands"
)

func Main() {
	hookContext := goops.NewHookContext()
	hookName := hookContext.Environment.JujuHookName()
	hookContext.Commands.JujuLog(commands.Info, "Hook name:", hookName)
	err := hookContext.Commands.StatusSet(commands.StatusActive, "A happy charm")
	if err != nil {
		hookContext.Commands.JujuLog(commands.Error, "Could not set status:", err.Error())
		os.Exit(0)
	}
	hookContext.Commands.JujuLog(commands.Info, "Status set to active")
	os.Exit(0)
}
