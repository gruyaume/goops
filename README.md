# goops

**Develop Reliable, Portable, and Fast Juju Charms in Go**

`goops` is a Go library for developing robust Juju charms. While charm developers traditionally use the [ops Python framework](https://github.com/canonical/operator), Python's dynamic typing and interpreter-based execution often lead to runtime errors and portability issues across different bases. In contrast, Go compiles to a single, self-contained binary, ensuring greater reliability and consistent behavior in any environment.

## Getting Started

### 1. Use the Charmcraft `go` plugin

Use the `go` plugin to build your charm in `charmcraft.yaml`:

```yaml
parts:
  charm:
    source: .
    plugin: go
    build-snaps:
      - go
    organize:
      bin/<your-charm-name>: dispatch
```

### 2. Write your charm

Create a `main.go` file under the `cmd/<your-charm-name>/` directory in your charm's root directory. This file will contain the main logic of your charm. Import the `goops` library and use its functions to interact with Juju. For example:

```go
package main

import (
	"os"

	"github.com/gruyaume/goops"
)

func main() {
	env := goops.ReadEnv()

	goops.LogInfof("Hook name: %s", env.HookName)

	err := goops.SetUnitStatus(goops.StatusActive, "A happy charm")
	if err != nil {
		goops.LogErrorf("Could not set status: %v", err)
		os.Exit(0)
	}

	goops.LogInfof("Status set to active")
	os.Exit(0)
}
```

The following charms use `goops` and can be used as examples:
- [Certificates](https://github.com/gruyaume/certificates-operator)
- [Notary K8s](https://github.com/gruyaume/notary-k8s-operator)

## Reference

### Design principles

- **Reliability**: Our top priority is building predictable and robust charms.
- **Simplicity**: `goops` serves as a minimal, one-to-one mapping between Juju concepts and Go constructs. It is not a framework; it does not impose charm design patterns.

### Unit Testing

[`goopstest`](goopstest/README.md) is a unit testing framework for `goops` charms. It allows you to simulate Juju environments and test your charm logic without needing a live Juju controller.

### Juju compatibility

`goops` is compatible with Juju 3.6 and later.
