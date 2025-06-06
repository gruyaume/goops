package charm_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
	"github.com/gruyaume/goops/internal/charm"
)

func TestCharm(t *testing.T) {
	ctx := goopstest.Context{
		Charm: charm.Configure,
	}

	stateIn := goopstest.State{}

	stateOut := ctx.Run("start", stateIn)

	if stateOut.UnitStatus != string(goops.StatusActive) {
		t.Errorf("Expected unit status to be ActiveStatus, got %v", stateOut.UnitStatus)
	}
}
