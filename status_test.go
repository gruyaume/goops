package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestStatusSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: nil,
		Err:    nil,
	}

	goops.SetRunner(fakeRunner)

	err := goops.SetUnitStatus(goops.StatusActive)
	if err != nil {
		t.Fatalf("StatusSet returned an error: %v", err)
	}

	if fakeRunner.Command != "status-set" {
		t.Errorf("Expected command %q, got %q", "status-set", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != string(goops.StatusActive) {
		t.Errorf("Expected argument %q, got %q", string(goops.StatusActive), fakeRunner.Args[0])
	}

	if fakeRunner.Output != nil {
		t.Errorf("Expected no output, got %q", string(fakeRunner.Output))
	}
}
