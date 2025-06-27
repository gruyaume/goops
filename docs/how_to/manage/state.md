---
description: Manage state with `goops` charms.
---

# Manage state

Juju allows charms to store state in a key-value store. Here we cover how you can use `goops` to set, get, and delete state in your charms. In this simple example, we check whether a state exists for a key named `my-key`. If it exists, we delete it. If it does not exist, we set it to a value of `my-value`.

```go
package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
)

const (
	StateKey   = "my-key"
	StateValue = "my-value"
)

func Configure() error {
	_, err := goops.GetState(StateKey)
	if err != nil {
		goops.LogInfof("could not get state: %s", err.Error())

		err := goops.SetState(StateKey, StateValue)
		if err != nil {
			return fmt.Errorf("could not set state: %w", err)
		}

		goops.LogInfof("set state: %s = %s", StateKey, StateValue)

		return nil
	}

	goops.LogInfof("state already set: %s = %s", StateKey, StateValue)

	err = goops.DeleteState(StateKey)
	if err != nil {
		return fmt.Errorf("could not delete state: %w", err)
	}

	goops.LogInfof("deleted state: %s", StateKey)

	return nil
}
```

!!! info
    Learn more about state management in charms:

    - [Juju Hook commands :octicons-link-external-24:](https://documentation.ubuntu.com/juju/3.6/reference/hook-command/list-of-hook-commands/)
    - [goops API reference :octicons-link-external-24:](https://pkg.go.dev/github.com/gruyaume/goops)
