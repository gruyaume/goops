package commands_test

import (
	"testing"

	"github.com/gruyaume/go-operator/internal/commands"
)

func TestRelationIDs_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["123", "456"]`),
		Err:    nil,
	}

	result, err := commands.RelationIDs(fakeRunner, "tls-certificates")
	if err != nil {
		t.Fatalf("RelationIDs returned an error: %v", err)
	}

	expectedOutput := []string{
		"123",
		"456",
	}
	if len(result) != len(expectedOutput) {
		t.Fatalf("Expected %d relation IDs, got %d", len(expectedOutput), len(result))
	}
	for i, id := range result {
		if id != expectedOutput[i] {
			t.Errorf("Expected relation ID %q, got %q", expectedOutput[i], id)
		}
	}
	if fakeRunner.Command != commands.RelationIDsCommand {
		t.Errorf("Expected command %q, got %q", commands.RelationIDsCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 argument, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != "tls-certificates" {
		t.Errorf("Expected argument %q, got %q", "tls-certificates", fakeRunner.Args[0])
	}
	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}
}
