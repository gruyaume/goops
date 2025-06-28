---
description: Develop reliable, portable, and fast Juju Charms in Go.
---

# goops

*Develop reliable, simple, and fast Juju Charms in Go*

`goops` is a Go library for developing simple and robust Juju charms. Write, test, build, and deploy your charm in minutes with goops.

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

## In this documentation

<div class="grid cards" markdown>

-   [__Tutorials__](tutorials/index.md)

    ---

    **Start here**: a hands-on introduction to goops for new users. Write, build and deploy your first charm in minutes.

-   [__How-to Guides__](how_to/index.md)

    ---

    **Step-by-step guides** covering key operations such as managing config, integrations, and workloads.

-   [__Reference__](reference/index.md)

    ---

    **Technical information** - API, example charms, charm libraries, best practices, and more.

-   [__Explanation__](explanation/index.md)

    ---

    **Discussion and clarification** of key topics like unit testing, and handling hooks.


</div>
