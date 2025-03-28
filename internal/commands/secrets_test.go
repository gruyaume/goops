package commands

import (
	"testing"
)

// FakeRunner implements Runner and records the command name and arguments.
type FakeRunner struct {
	Command string
	Args    []string
	// Simulated output and error to return.
	Output []byte
	Err    error
}

// Run records the command name and arguments and returns preset output and error.
func (f *FakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.Command = name
	f.Args = args
	return f.Output, f.Err
}

// TestSecretAdd_Success verifies that SecretAdd builds the correct command line and returns the expected output.
func TestSecretAdd_Success(t *testing.T) {
	// Set up test input.
	content := map[string]string{
		"username": "user1",
		"password": "pass1",
	}
	description := "my secret"
	label := "my-label"

	fakeRunner := &FakeRunner{
		Output: []byte(`{"result":"success"}`),
		Err:    nil,
	}

	result, err := SecretAdd(fakeRunner, content, description, label)
	if err != nil {
		t.Fatalf("SecretAdd returned an error: %v", err)
	}

	expectedOutput := `{"result":"success"}`
	if result != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, result)
	}

	if fakeRunner.Command != SecretAddCommand {
		t.Errorf("Expected command %q, got %q", SecretAddCommand, fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 4 {
		t.Fatalf("Expected 4 arguments, got %d", len(fakeRunner.Args))
	}

	contentArg := fakeRunner.Args[0]
	if fakeRunner.Args[0] != "username=user1" {
		t.Errorf("Expected content arg %q, got %q", "username=user1", contentArg)
	}

	if fakeRunner.Args[1] != "password=pass1" {
		t.Errorf("Expected content arg %q, got %q", "password=pass1", fakeRunner.Args[1])
	}

	if fakeRunner.Args[2] != "--description=my secret" {
		t.Errorf("Expected description arg %q, got %q", "--description=my secret", fakeRunner.Args[2])
	}

	if fakeRunner.Args[3] != "--label=my-label" {
		t.Errorf("Expected label arg %q, got %q", "--label=my-label", fakeRunner.Args[3])
	}

}

// TestSecretAdd_EmptyContent verifies that calling SecretAdd with an empty content map returns an error.
func TestSecretAdd_EmptyContent(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(""),
		Err:    nil,
	}

	_, err := SecretAdd(fakeRunner, map[string]string{}, "desc", "label")
	if err == nil {
		t.Error("Expected error when content is empty, but got nil")
	}
}
