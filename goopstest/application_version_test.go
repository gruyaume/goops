package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ApplicationVersion() error {
	err := goops.SetAppVersion("1.2.3")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmApplicationVersion(t *testing.T) {
	ctx := goopstest.NewContext(ApplicationVersion)

	stateIn := goopstest.State{}

	stateOut := ctx.Run("start", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
	}

	if stateOut.ApplicationVersion != "1.2.3" {
		t.Errorf("Expected application version '1.2.3', got '%s'", stateOut.ApplicationVersion)
	}
}

func TestCharmApplicationVersionInActionHook(t *testing.T) {
	ctx := goopstest.NewContext(ApplicationVersion)

	stateIn := goopstest.State{}

	stateOut, err := ctx.RunAction("run-action", stateIn, nil)
	if err != nil {
		t.Fatalf("RunAction returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
	}

	if stateOut.ApplicationVersion != "1.2.3" {
		t.Errorf("Expected application version '1.2.3', got '%s'", stateOut.ApplicationVersion)
	}
}
