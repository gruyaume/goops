# goops

*Develop reliable, simple, and fast Juju Charms in Go*

`goops` is a Go library for developing simple and robust Juju charms. Write, test, build, and deploy your charm in minutes with goops.

[**Get Started Now!**](https://gruyaume.github.io/goops/tutorials/getting_started/)

---

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
