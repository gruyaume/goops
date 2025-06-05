package main

import (
	"os"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/internal/charm"
)

// Example charm using `goops`
func main() {
	hookContext := goops.NewHookContext()

	env := goops.ReadEnv()

	if env.ActionName != "" {
		goops.LogInfof("Action name: %s", env.ActionName)

		switch env.ActionName {
		case "get-ca-certificate":
			err := charm.HandleGetCACertificateAction(hookContext)
			if err != nil {
				goops.LogErrorf("Error handling get-ca-certificate action: %s", err.Error())
				os.Exit(0)
			}

			goops.LogInfof("Handled get-ca-certificate action successfully")
			os.Exit(0)
		default:
			goops.LogErrorf("Action '%s' not recognized, exiting", env.ActionName)
			os.Exit(0)
		}
	}

	if env.HookName != "" {
		goops.LogInfof("Hook name: %s", env.HookName)

		err := charm.HandleDefaultHook(hookContext)
		if err != nil {
			goops.LogErrorf("Error handling default hook: %s", err.Error())
			os.Exit(0)
		}
	}
}
