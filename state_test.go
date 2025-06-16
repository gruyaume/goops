package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestStateDelete_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.DeleteState("key")
	if err != nil {
		t.Fatalf("StateDelete returned an error: %v", err)
	}

	if fakeRunner.Command != "state-delete" {
		t.Errorf("Expected command %q, got %q", "state-delete", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "key" {
		t.Errorf("Expected key arg %q, got %q", "key", fakeRunner.Args[0])
	}
}

func TestStateGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"value"`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	state, err := goops.GetState("key")
	if err != nil {
		t.Fatalf("StateGet returned an error: %v", err)
	}

	if fakeRunner.Command != "state-get" {
		t.Errorf("Expected command %q, got %q", "state-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "key" {
		t.Errorf("Expected key arg %q, got %q", "key", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected format arg %q, got %q", "--format=json", fakeRunner.Args[1])
	}

	if state != "value" {
		t.Errorf("Expected state %q, got %q", "value", state)
	}
}
