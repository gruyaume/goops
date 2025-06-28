---
description: Manage action with `goops` charms.
---

# How-to manage actions

Juju users can run actions in charms using the `juju run` command. Here we cover how you can use `goops` to handle those actions in your charm.

## 1. Declare actions

Declare the actions in you charm's `charmcraft.yaml` file. For example:

```yaml
actions:
  get-password:
    description: Return the password for the specified user.
    params:
      username:
        type: string
        description: >-
            The username for which to return the password.
    required: [username]
```

!!! note
    For more information on the `charmcraft.yaml` charm definition, read the [official charmcraft documentation](https://canonical-charmcraft.readthedocs-hosted.com/stable/reference/files/charmcraft-yaml-file/).

## 2. Read the action name

You can read the action name in your charm using `ReadEnv()`:

```go
package main

import (
	"os"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/internal/charm"
)

func main() {
	env := goops.ReadEnv()

	if env.ActionName != "" {
		goops.LogInfof("Action name: %s", env.ActionName)

		switch env.ActionName {
		case "get-password":
			err := charm.HandleGetPasswordAction()
			if err != nil {
				goops.LogErrorf("Error handling get-password action: %s", err.Error())
				os.Exit(1)
			}

		default:
			goops.LogErrorf("Action '%s' not recognized, exiting", env.ActionName)
			os.Exit(1)
		}
	}
}
```

## 3. Handle the action

You can handle the action in your charm. Here the charm reads the `username` parameter from the action and returns a password based on that username:

```go
package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
)

type GetPasswordActionParams struct {
	Username string `json:"username"`
}

func HandleGetPasswordAction() error {
	params := GetPasswordActionParams{}

	err := goops.GetActionParams(&params)
	if err != nil {
		goops.FailActionf("could not get action parameters")
		return fmt.Errorf("could not get action parameters: %w", err)
	}

	if params.Username == "" {
		goops.FailActionf("Username is not set in action parameters")
		return nil
	}

	password := fmt.Sprintf("%s-12345", params.Username)

	err = goops.SetActionResults(map[string]string{
		"password": password,
	})

	if err != nil {
		return fmt.Errorf("could not set action result: %w", err)
	}

	return nil
}
```

!!! warning
    All action-related functions (ex. `GetActionParams`, `SetActionResults`) will fail if they are not called in the context of an action hook. I.e. `env.ActionName` must not be empty.

!!! info
    Learn more about action management in charms:

    - [Juju Hook commands :octicons-link-external-24:](https://documentation.ubuntu.com/juju/3.6/reference/hook-command/list-of-hook-commands/)
    - [goops API reference :octicons-link-external-24:](https://pkg.go.dev/github.com/gruyaume/goops)
