# goopstest

**The unit testing framework for Goops charms.**

`goopstest` follows the same design principles as [ops-testing](https://ops.readthedocs.io/en/latest/reference/ops-testing.html#ops-testing), allowing users to write unit tests in a "state-transition" style. Each test consists of:
- A Context and an initial state (Arrange)
- An event (Act)
- An output state (Assert)

```go
package charm_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func Configure() error {
	isLeader, err := goops.IsLeader()
	if err != nil {
		return err
	}

	if !isLeader {
		_ = goops.SetUnitStatus(goops.StatusBlocked, "Unit is not a leader")
		return nil
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "Charm is active")

	return nil
}

func TestCharm(t *testing.T) {
	ctx := goopstest.Context{
		Charm: Configure,
	}

	stateIn := &goopstest.State{
		Leader: false,
	}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.UnitStatus != string(goops.StatusBlocked) {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, goops.StatusBlocked)
	}
}
```

### API Documentation

The API documentation for `goopstest` is available at [pkg.go.dev/github.com/gruyaume/goops/goopstest](https://pkg.go.dev/github.com/gruyaume/goops/goopstest).
