# goops

**Develop Reliable, Portable, and Fast Juju Charms in Go**

`goops` is a Go library for developing robust Juju charms. Go charms compile to a single, self-contained binary, ensuring greater reliability and consistent behavior accross different charm bases.

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
	"fmt"

	"github.com/gruyaume/goops"
)

const (
	APIPort = 2111
)

func main() {
	env := goops.ReadEnv()

	err := Configure()
	if err != nil {
		goops.LogErrorf("Error handling %s hook: %v", env.HookName, err)
		os.Exit(1)
	}
}

func Configure() error {
	isLeader, err := goops.IsLeader()
	if err != nil {
		return fmt.Errorf("could not check if unit is leader: %w", err)
	}

	if !isLeader {
		goops.SetUnitStatus(goops.StatusBlocked, "Unit is not leader")
		return nil
	}

	err = goops.SetPorts([]*goops.Port{
		{
			Port:     APIPort,
			Protocol: "tcp",
		},
	})
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	goops.LogInfof("Port is set: %d/tcp", APIPort)

	goops.SetUnitStatus(goops.StatusActive, "")

	return nil
}
```

You can then pack and deploy your charm as usual with `charmcraft pack` and `juju deploy`.

## Reference

### API Documentation

The API documentation for `goops` is available at [pkg.go.dev/github.com/gruyaume/goops](https://pkg.go.dev/github.com/gruyaume/goops).

### Example Charms

The following charms use `goops` and can be used as reference implementations:
- [Certificates](https://github.com/gruyaume/certificates-operator)
- [Notary K8s](https://github.com/gruyaume/notary-k8s-operator)

### Charm Libraries

Charm Libraries are maintained centrally at [github.com/gruyaume/charm-libraries](https://github.com/gruyaume/charm-libraries).

### Unit Testing

[`goopstest`](goopstest/README.md) is a unit testing framework for `goops` charms. It allows you to simulate Juju environments and test your charm logic without needing a live Juju controller.

### Design principles

- **Reliability**: Our top priority is building predictable and robust charms.
- **Simplicity**: `goops` serves as a minimal mapping between Juju concepts and Go constructs. It is not a framework; it does not impose charm design patterns.

### Contributing

Contributions to `goops` are welcome! For details on how to contribute, please refer to the [Contributing Guide](CONTRIBUTING.md).

### Discussion

If you have questions, feedback, announcements, or ideas, please [Open a New Discussion](https://github.com/gruyaume/goops/discussions).

### Juju compatibility

`goops` is compatible with Juju 3.6 and later.
