package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ApplicationVersion() error {
	err := goops.SetApplicationVersion("1.2.3")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmApplicationVersion(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ApplicationVersion,
	}

	stateIn := &goopstest.State{}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.ApplicationVersion != "1.2.3" {
		t.Errorf("Expected application version '1.2.3', got '%s'", stateOut.ApplicationVersion)
	}
}
