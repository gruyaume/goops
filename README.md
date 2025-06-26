# goops

**Develop Reliable, Portable, and Fast Juju Charms in Go**

`goops` is a Go library for developing robust Juju charms. Go charms compile to a single, self-contained binary, ensuring greater reliability and consistent behavior accross different charm bases.

[Get Started Now!](https://gruyaume.github.io/goops/tutorials/getting_started/)

```go
package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
)

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
			Port:     2111,
			Protocol: "tcp",
		},
	})
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	goops.SetUnitStatus(goops.StatusActive, fmt.Sprintf("Port %d/tcp is set", 2111))

	return nil
}
```

## Quick Links

- [Documentation](https://gruyaume.github.io/goops/)
- [Discussion](https://github.com/gruyaume/goops/discussions)
- [Contributing](CONTRIBUTING.md)
