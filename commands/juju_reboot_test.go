package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestJujuReboot_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	err := command.JujuReboot(true)
	if err != nil {
		t.Fatalf("JujuReboot returned an error: %v", err)
	}

	if fakeRunner.Command != commands.JujuRebootCommand {
		t.Errorf("Expected command %q, got %q", commands.JujuRebootCommand, fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--now" {
		t.Errorf("Expected argument %q, got %q", "--now", fakeRunner.Args[0])
	}
}
