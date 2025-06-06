package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ConfigureActive() error {
	_ = goops.SetUnitStatus(goops.StatusActive, "Charm is active")
	return nil
}

func ConfigureBlocked() error {
	_ = goops.SetUnitStatus(goops.StatusBlocked, "This is a test message")
	return nil
}

func ConfigureWaiting() error {
	_ = goops.SetUnitStatus(goops.StatusWaiting, "Waiting for something")
	return nil
}

func ConfigureMaintenance() error {
	_ = goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance")
	return nil
}

func TestCharm(t *testing.T) {
	tests := []struct {
		name      string
		configure func() error
		want      string
	}{
		{
			name:      "ActiveStatus",
			configure: ConfigureActive,
			want:      string(goops.StatusActive),
		},
		{
			name:      "BlockedStatus",
			configure: ConfigureBlocked,
			want:      string(goops.StatusBlocked),
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			ctx := goopstest.Context{
				Charm: tc.configure,
			}

			stateIn := goopstest.State{}

			stateOut, err := ctx.Run("start", stateIn)
			if err != nil {
				t.Fatalf("Run returned an error: %v", err)
			}

			if stateOut.UnitStatus != tc.want {
				t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, tc.want)
			}
		})
	}
}
