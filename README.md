# go-operator

> :construction: **Beta Notice**
> Go-operator is in beta. If you encounter any issues, please [report them here](https://github.com/gruyaume/go-operator/issues). 

**Develop Reliable, Portable, and Fast Juju Charms in Go**

`go-operator` is a Go library that empowers you to create robust Juju charms. While Python is traditionally used for charm development, its dynamic typing and interpreter-based execution often lead to runtime errors and portability issues across different bases. In contrast, Go compiles to a single, self-contained binary, ensuring greater reliability and consistent behavior in any environment.

## Try it now

### 1. Use the Charmcraft `go` plugin

Use the `go` plugin to build your charm. For example:

```yaml
name: go-operator
summary: Write Juju charms in Go.
description: |
  Write Juju charms in Go.

type: charm
base: ubuntu@24.04
build-base: ubuntu@24.04
platforms:
  amd64:

parts:
  charm:
    source: .
    plugin: go
    build-snaps:
      - go
    organize:
      bin/go-operator: dispatch  # replace `go-operator` with your binary name

config:
  options:
    key:
      type: string
      default: whatever value
      description: >
        Example configuration option
```

### 2. Write your charm

In your charm's root directory, create a `main.go` file under the `cmd/<your-charm-name>` directory. This file will contain the main logic of your charm.

Import the `go-operator` library and use its functions to interact with Juju. Here's a simple example of a charm that checks if it is leader, reads a configuration option, and sets its status to active:

```go
package main

import (
	"os"

	"github.com/gruyaume/go-operator/internal/commands"
	"github.com/gruyaume/go-operator/internal/environment"
)

func main() {
	commandRunner := &commands.DefaultRunner{}
	environmentGetter := &environment.DefaultEnvironment{}
	logger := commands.NewLogger(commandRunner)
	hookName := environment.JujuHookName(environmentGetter)
	logger.Info("Hook name:", hookName)

	isLeader, err := commands.IsLeader(commandRunner)
	if err != nil {
		logger.Info("Could not check if leader:", err.Error())
		os.Exit(0)
	}
	if !isLeader {
		logger.Info("not leader, exiting")
		os.Exit(0)
	}
	logger.Info("Unit is leader")

	keyConfig, err := commands.ConfigGet(commandRunner, "key")
	if err != nil {
		logger.Error("Could not get config:", err.Error())
		os.Exit(0)
	}
	if keyConfig == "" {
		logger.Error("Configuration option `key` is empty:", err.Error())
		os.Exit(0)
	}
	err = commands.StatusSet(commandRunner, commands.StatusActive)
	if err != nil {
		logger.Error("Could not set status:", err.Error())
		os.Exit(0)
	}
	logger.Info("Status set to active")
	os.Exit(0)
}
```

### 3. Build and deploy your charm

Just like for any other charm, you can build your charm using `charmcraft`:

```bash
charmcraft pack
```

```bash
juju deploy ./<your-charm-name>.charm
```

## Design principles

- **Reliability**: Building predictable, robust charms is our top priority.
- **Simplicity**: We maintain a clear, one-to-one mapping between Juju concepts and Go constructs.
