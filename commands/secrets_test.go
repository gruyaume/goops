package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestSecretIDs_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["123", "456"]`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	result, err := command.SecretIDs()
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

	if fakeRunner.Command != "secret-ids" {
		t.Errorf("Expected command %q, got %q", "secret-ids", fakeRunner.Command)
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
	command := commands.Command{
		Runner: fakeRunner,
	}

	result, err := command.SecretGet("123", "my-label", false, true)
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

	if fakeRunner.Command != "secret-get" {
		t.Errorf("Expected command %q, got %q", "secret-get", fakeRunner.Command)
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
	command := commands.Command{
		Runner: fakeRunner,
	}

	result, err := command.SecretAdd(content, description, label)
	if err != nil {
		t.Fatalf("SecretAdd returned an error: %v", err)
	}

	expectedOutput := `{"result":"success"}`
	if result != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, result)
	}

	if fakeRunner.Command != "secret-add" {
		t.Errorf("Expected command %q, got %q", "secret-add", fakeRunner.Command)
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
	command := commands.Command{
		Runner: fakeRunner,
	}

	_, err := command.SecretAdd(map[string]string{}, "desc", "label")
	if err == nil {
		t.Error("Expected error when content is empty, but got nil")
	}
}

func TestSecretGrant_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"result":"success"}`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	err := command.SecretGrant("123", "certificates:0", "")
	if err != nil {
		t.Fatalf("SecretGrant returned an error: %v", err)
	}

	if fakeRunner.Command != "secret-grant" {
		t.Fatalf("Expected command %q, got %q", "secret-grant", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "123" {
		t.Errorf("Expected ID arg %q, got %q", "123", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--relation=certificates:0" {
		t.Errorf("Expected secret ID arg %q, got %q", "--relation=certificates:0", fakeRunner.Args[1])
	}
}
