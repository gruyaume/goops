package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestIsLeader_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`true`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	result, err := command.IsLeader()
	if err != nil {
		t.Fatalf("IsLeader returned an error: %v", err)
	}

	if result != true {
		t.Errorf("Expected true, got %v", result)
	}

	if fakeRunner.Command != "is-leader" {
		t.Errorf("Expected command %q, got %q", "is-leader", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[0])
	}
}
