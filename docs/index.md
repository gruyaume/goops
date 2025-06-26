---
description: Develop reliable, portable, and fast Juju Charms in Go.
---

# goops

**goops** is a Go library for developing robust Juju charms. Go charms compile to a single, self-contained binary, ensuring greater reliability and consistent behavior accross different charm bases.

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

    **Start here**: a hands-on introduction to goops for new users. Build and deploy a Go charm in minutes.

-   [__How-to Guides__](how_to/index.md)

    ---

    **Step-by-step guides** covering key operation and common tasks.

-   [__Reference__](reference/index.md)

    ---

    **Technical information** - API, configuration, and more.

-   [__Explanation__](explanation/index.md)

    ---

    **Discussion and clarification** of key topics.


</div>
