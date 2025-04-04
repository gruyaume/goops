package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestApplicationVersionSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	version := "1.2.3"
	applicationVersionSetOptions := &commands.ApplicationVersionSetOptions{
		Version: version,
	}

	err := command.ApplicationVersionSet(applicationVersionSetOptions)
	if err != nil {
		t.Fatalf("ApplicationVersionSet returned an error: %v", err)
	}

	if fakeRunner.Command != "application-version-set" {
		t.Errorf("Expected command %q, got %q", "application-version-set", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != version {
		t.Errorf("Expected argument %q, got %q", version, fakeRunner.Args[0])
	}
}
