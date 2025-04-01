package main

import (
	"os"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/goops/internal/charm"
)

// Example charm using `goops`
func main() {
	hookContext := goops.NewHookContext()
	actionName := hookContext.Environment.JujuActionName()

	if actionName != "" {
		hookContext.Commands.JujuLog(commands.Info, "Action name: "+actionName)

		switch actionName {
		case "get-ca-certificate":
			err := charm.HandleGetCACertificateAction(hookContext)
			if err != nil {
				hookContext.Commands.JujuLog(commands.Error, "Error handling get-ca-certificate action: "+err.Error())
				os.Exit(0)
			}

			hookContext.Commands.JujuLog(commands.Info, "Handled get-ca-certificate action successfully")
			os.Exit(0)
		default:
			hookContext.Commands.JujuLog(commands.Error, "Action not recognized, exiting")
			os.Exit(0)
		}
	}

	hookName := hookContext.Environment.JujuHookName()
	if hookName != "" {
		hookContext.Commands.JujuLog(commands.Info, "Hook name: "+hookName)

		err := charm.HandleDefaultHook(hookContext)
		if err != nil {
			hookContext.Commands.JujuLog(commands.Error, "Error handling default hook: "+err.Error())
			os.Exit(0)
		}
	}
}
