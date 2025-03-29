package commands_test

import (
	"testing"

	"github.com/gruyaume/go-operator/commands"
)

func TestStatusSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: nil,
		Err:    nil,
	}

	err := commands.StatusSet(fakeRunner, commands.StatusActive)
	if err != nil {
		t.Fatalf("StatusSet returned an error: %v", err)
	}

	if fakeRunner.Command != commands.StatusSetCommand {
		t.Errorf("Expected command %q, got %q", commands.StatusSetCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != string(commands.StatusActive) {
		t.Errorf("Expected argument %q, got %q", string(commands.StatusActive), fakeRunner.Args[0])
	}
	if fakeRunner.Output != nil {
		t.Errorf("Expected no output, got %q", string(fakeRunner.Output))
	}
}
