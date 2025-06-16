package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestJujuLogStatusSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: nil,
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	goops.LogDebugf("my message")

	if fakeRunner.Command != "juju-log" {
		t.Errorf("Expected command %q, got %q", "juju-log", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--log-level=DEBUG" {
		t.Errorf("Expected argument %q, got %q", "--log-level=DEBUG", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "my message" {
		t.Errorf("Expected argument %q, got %q", "my message", fakeRunner.Args[1])
	}

	if fakeRunner.Output != nil {
		t.Errorf("Expected no output, got %q", string(fakeRunner.Output))
	}
}
