---
description: Manage config with `goops` charms.
---

# Manage config

Juju users can configure charms using the `juju config` command. Here we cover how you can use `goops` to read those configuration options in your charm.

## 1. Declare configuration options

Declare the configuration options in you charm's `charmcraft.yaml` file. For example:

```yaml
config:
  options:
    username:
      type: string
      default: gruyaume
      description: >
        Example configuration option for this charm.
```

!!! note
    For more information on the `charmcraft.yaml` charm definition, read the [official charmcraft documentation](https://canonical-charmcraft.readthedocs-hosted.com/stable/reference/files/charmcraft-yaml-file/).

## 2. Read configuration options

You can read the configuration options in your charm using `GetConfig()`:

```go
package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
)

type Config struct {
	Username string `json:"username"`
}

func Configure() error {
	c := Config{}

	err := goops.GetConfig(c)
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	goops.LogInfof("Configuring charm with username: %s", c.Username)

	return nil
}
```

!!! info
    Learn more about config management in charms:

    - [Juju Hook commands :octicons-link-external-24:](https://documentation.ubuntu.com/juju/3.6/reference/hook-command/list-of-hook-commands/)
    - [goops API reference :octicons-link-external-24:](https://pkg.go.dev/github.com/gruyaume/goops)
