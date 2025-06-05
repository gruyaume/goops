package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestJujuReboot_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetRunner(fakeRunner)

	err := goops.Reboot(true)
	if err != nil {
		t.Fatalf("Reboot returned an error: %v", err)
	}

	if fakeRunner.Command != "juju-reboot" {
		t.Errorf("Expected command %q, got %q", "juju-reboot", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--now" {
		t.Errorf("Expected argument %q, got %q", "--now", fakeRunner.Args[0])
	}
}
