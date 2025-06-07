# goopstest

The unit testing framework for Goops charms.

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
