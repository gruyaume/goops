---
description: Handling hooks with `goops` charms.
---

# Handling Hooks

In Juju, a hook is a notification to the charm that the internal representation of Juju has changed in a way that requires a reaction from the charm so that the unit’s state and the controller’s state can be reconciled. You can find the list of available hooks and more hook-related information in the [Juju documentation](https://documentation.ubuntu.com/juju/3.6/reference/hook/).

## Hook execution

This notification is sent in the form of Juju executing the charm binary with a specific set of environment variables. The charm binary is expected to handle the hook by comparing its existing state with the intended state and performing the necessary actions to reconcile them.

Juju expects the charm binary to be named `dispatch` and to be located in the charm root directory. This is why we inform charmcraft to rename the binary to `dispatch` in the `charmcraft.yaml` file.

```yaml
parts:
  charm:
    source: .
    plugin: go
    build-snaps:
      - go
    organize:
      bin/example: dispatch
```

## Handling hooks with `goops`

During the hook execution `goops` provides access to the following:

- **Hook Commands**: `goops` exposes every Juju [hook commands](https://documentation.ubuntu.com/juju/3.6/reference/hook-command/list-of-hook-commands/), as a Go function.
- **Environment Variables**: `goops` provides access to every Juju-defined [environment variables](https://documentation.ubuntu.com/juju/3.6/reference/hook/#hook-execution).
- **Charm metadata**: `goops` provides access to the charm metadata as defined in `charmcraft.yaml`.
- **Pebble**: `goops` provides access to the Pebble API, allowing you to manage services and containers for Kubernetes charms.

### Independent hook handling

The hook name is made available by Juju through an environment variable. In your charm, you can access it using the `ReadEnv()` function and handle it accordingly.

```go
package main

import (
	"os"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/internal/charm"
)

func main() {
	env := goops.ReadEnv()

	switch env.HookName {
	case "install":
		err := charm.Install()
		if err != nil {
			goops.LogErrorf("Error handling install hook: %s", err.Error())
			os.Exit(1)
		}
	case "remove":
		err := charm.Remove()
		if err != nil {
			goops.LogErrorf("Error handling remove hook: %s", err.Error())
			os.Exit(1)
		}
	case "":
		goops.LogInfof("No hook name provided, running default configuration.")
	default:
		err := charm.Configure()
		if err != nil {
			goops.LogErrorf("Error handling default hook: %s", err.Error())
			os.Exit(1)
		}
	}
}
```
