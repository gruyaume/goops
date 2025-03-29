# go-operator

> :construction: **Beta Notice**
> Go-operator is in beta. If you encounter any issues, please [report them here](https://github.com/gruyaume/go-operator/issues). 

**Develop Reliable, Portable, and Fast Juju Charms in Go**

`go-operator` is a Go library for developing robust Juju charms. While developers traditionally use the [ops Python framework](https://github.com/canonical/operator) for developing charms, Python's dynamic typing and interpreter-based execution often lead to runtime errors and portability issues across different bases. In contrast, Go compiles to a single, self-contained binary, ensuring greater reliability and consistent behavior in any environment.

## Try it now

### 1. Use the Charmcraft `go` plugin

Use the `go` plugin to build your charm in `charmcraft.yaml`:

```yaml
...
parts:
  charm:
    source: .
    plugin: go
    build-snaps:
      - go
    organize:
      bin/go-operator: dispatch  # replace `go-operator` with your binary name
```

### 2. Write your charm

In your charm's root directory, create a `main.go` file under the `cmd/<your-charm-name>` directory. This file will contain the main logic of your charm. Import the `go-operator` library and use its functions to interact with Juju. For example:

```go
package main

import (
	"os"

	"github.com/gruyaume/go-operator/commands"
	"github.com/gruyaume/go-operator/environment"
)

func main() {
	commandRunner := &commands.DefaultRunner{}
	environmentGetter := &environment.DefaultEnvironment{}
	logger := commands.NewLogger(commandRunner)
	hookName := environment.JujuHookName(environmentGetter)
	logger.Info("Hook name:", hookName)
	err := commands.StatusSet(commandRunner, commands.StatusActive)
	if err != nil {
		logger.Error("Could not set status:", err.Error())
		os.Exit(0)
	}
	logger.Info("Status set to active")
	os.Exit(0)
}
```

## Design principles

- **Reliability**: Building predictable, robust charms is our top priority.
- **Simplicity**: `go-operator` serves as a minimal, one-to-one mapping between Juju concepts and Go constructs. It is not a framework; it does not impose charm design patterns. The library has no dependencies.
