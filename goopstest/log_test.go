package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ContainerLog() error {
	goops.LogInfof("This is an info log message")

	return nil
}

func TestContainerLog(t *testing.T) {
	ctx := goopstest.NewContext(ContainerLog)

	stateIn := goopstest.State{}

	_, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
	}

	expectedLog := goopstest.JujuLogLine{
		Message: "This is an info log message",
		Level:   goopstest.LogLevelInfo,
	}
	found := false

	for _, logEntry := range ctx.JujuLog {
		if logEntry == expectedLog {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected log message %q not found in Juju log", expectedLog)
	}
}
