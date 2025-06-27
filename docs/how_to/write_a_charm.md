---
description: Write a Juju charm using `goops`
---

# Write a Charm

Create a new directory for your charm project and initialize a Go module. Here the charm name is `example`, replace it with your desired charm name:

```bash
mkdir example
cd example
go mod init example
```

Create a `cmd/example/main.go` file with the following content:

```go
package main

import (
	"os"

	"example/internal/charm"

	"github.com/gruyaume/goops"
)

func main() {
	env := goops.ReadEnv()

	goops.LogInfof("Hook name: %s", env.HookName)

	err := charm.Configure()
	if err != nil {
		goops.LogErrorf("Error handling hook: %s", err.Error())
		os.Exit(1)
	}
}
```

Create a `internal/charm/charm.go` file with the charm logic. For example:

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
	isLeader, err := goops.IsLeader()
	if err != nil {
		return fmt.Errorf("could not check if unit is leader: %w", err)
	}

	if !isLeader {
		err := goops.SetUnitStatus(goops.StatusBlocked, "Unit is not leader")
		if err != nil {
			return fmt.Errorf("could not set unit status: %w", err)
		}
		return nil
	}

	goops.LogInfof("Unit is leader")

	var config Config
	err = goops.GetConfig(&config)
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	if config.Username == "" {
		err := goops.SetUnitStatus(goops.StatusBlocked, "Username is not set in config")
		if err != nil {
			return fmt.Errorf("could not set unit status: %w", err)
		}
		return nil
	}

	err = goops.SetUnitStatus(goops.StatusActive, fmt.Sprintf("Username is set to '%s'", config.Username))
	if err != nil {
		return fmt.Errorf("could not set unit status: %w", err)
	}

	return nil
}
```

Install go dependencies:

```bash
go mod tidy
```
