package commands_test

import (
	"testing"

	"github.com/gruyaume/go-operator/commands"
)

func TestSecretIDs_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["123", "456"]`),
		Err:    nil,
	}

	result, err := commands.SecretIDs(fakeRunner)
	if err != nil {
		t.Fatalf("SecretIDs returned an error: %v", err)
	}

	expectedOutput := []string{
		"123",
		"456",
	}
	if len(result) != len(expectedOutput) {
		t.Fatalf("Expected %d secret IDs, got %d", len(expectedOutput), len(result))
	}
	for i, id := range result {
		if id != expectedOutput[i] {
			t.Errorf("Expected secret ID %q, got %q", expectedOutput[i], id)
		}
	}
	if fakeRunner.Command != commands.SecredIDsCommand {
		t.Errorf("Expected command %q, got %q", commands.SecredIDsCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[0])
	}
}

func TestSecretGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"username":"user1","password":"pass1"}`),
		Err:    nil,
	}

	result, err := commands.SecretGet(fakeRunner, "123", "my-label", false, true)
	if err != nil {
		t.Fatalf("SecretGet returned an error: %v", err)
	}

	expectedOutput := map[string]string{
		"username": "user1",
		"password": "pass1",
	}
	if len(result) != len(expectedOutput) {
		t.Fatalf("Expected %d secret content keys, got %d", len(expectedOutput), len(result))
	}
	for key, value := range result {
		if value != expectedOutput[key] {
			t.Errorf("Expected secret content %q, got %q", expectedOutput[key], value)
		}
	}

	if fakeRunner.Command != commands.SecretGetCommand {
		t.Errorf("Expected command %q, got %q", commands.SecretGetCommand, fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 4 {
		t.Fatalf("Expected 4 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "123" {
		t.Errorf("Expected ID arg %q, got %q", "123", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--label=my-label" {
		t.Errorf("Expected label arg %q, got %q", "--label=my-label", fakeRunner.Args[1])
	}

	if fakeRunner.Args[2] != "--refresh" {
		t.Errorf("Expected refresh arg %q, got %q", "--refresh", fakeRunner.Args[2])
	}

	if fakeRunner.Args[3] != "--format=json" {
		t.Errorf("Expected format arg %q, got %q", "--format=json", fakeRunner.Args[3])
	}
}

func TestSecretAdd_Success(t *testing.T) {
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

	result, err := commands.SecretAdd(fakeRunner, content, description, label)
	if err != nil {
		t.Fatalf("SecretAdd returned an error: %v", err)
	}

	expectedOutput := `{"result":"success"}`
	if result != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, result)
	}

	if fakeRunner.Command != commands.SecretAddCommand {
		t.Errorf("Expected command %q, got %q", commands.SecretAddCommand, fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 4 {
		t.Fatalf("Expected 4 arguments, got %d", len(fakeRunner.Args))
	}

	contentArg := fakeRunner.Args[0]
	if fakeRunner.Args[0] != "username=user1" && fakeRunner.Args[1] != "username=user1" {
		t.Errorf("Expected content arg %q, got %q", "username=user1", contentArg)
	}

	if fakeRunner.Args[0] != "password=pass1" && fakeRunner.Args[1] != "password=pass1" {
		t.Errorf("Expected content arg %q, got %q", "password=pass1", contentArg)
	}

	if fakeRunner.Args[2] != "--description=my secret" {
		t.Errorf("Expected description arg %q, got %q", "--description=my secret", fakeRunner.Args[2])
	}

	if fakeRunner.Args[3] != "--label=my-label" {
		t.Errorf("Expected label arg %q, got %q", "--label=my-label", fakeRunner.Args[3])
	}

}

func TestSecretAdd_EmptyContent(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(""),
		Err:    nil,
	}

	_, err := commands.SecretAdd(fakeRunner, map[string]string{}, "desc", "label")
	if err == nil {
		t.Error("Expected error when content is empty, but got nil")
	}
}
