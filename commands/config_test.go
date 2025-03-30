package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestConfigGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"banana"`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	result, err := command.ConfigGet("fruit")
	if err != nil {
		t.Fatalf("ConfigGet returned an error: %v", err)
	}

	if _, ok := result.(string); !ok {
		t.Fatalf("Expected result to be a string, got %T", result)
	}

	if result != "banana" {
		t.Fatalf("Expected %q, got %q", "banana", result)
	}

	if fakeRunner.Command != commands.ConfigGetCommand {
		t.Errorf("Expected command %q, got %q", commands.ConfigGetCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != "fruit" {
		t.Errorf("Expected argument %q, got %q", "fruit", fakeRunner.Args[0])
	}
	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}
}

func TestConfigGetString_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"banana"`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	result, err := command.ConfigGetString("fruit")
	if err != nil {
		t.Fatalf("ConfigGetString returned an error: %v", err)
	}
	if result != "banana" {
		t.Fatalf("Expected %q, got %q", "banana", result)
	}

}

func TestConfigGetString_BadType(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`123`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	_, err := command.ConfigGetString("fruit")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if err.Error() != "config value is not a string: 123" {
		t.Fatalf("Expected error %q, got %q", "config value is not a string: 123", err.Error())
	}
}

func TestConfigGetInt_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`123`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	result, err := command.ConfigGetInt("fruit")
	if err != nil {
		t.Fatalf("ConfigGetInt returned an error: %v", err)
	}
	if result != 123 {
		t.Fatalf("Expected %d, got %d", 123, result)
	}
}

func TestConfigGetInt_BadType(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"banana"`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	_, err := command.ConfigGetInt("fruit")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if err.Error() != "config value is not a number: banana" {
		t.Fatalf("Expected error %q, got %q", "config value is not a number: banana", err.Error())
	}
}

func TestConfigGetBool_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`true`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	result, err := command.ConfigGetBool("fruit")
	if err != nil {
		t.Fatalf("ConfigGetBool returned an error: %v", err)
	}
	if result != true {
		t.Fatalf("Expected %t, got %t", true, result)
	}
}

func TestConfigGetBool_BadType(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`123`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	_, err := command.ConfigGetBool("fruit")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if err.Error() != "config value is not a bool: 123" {
		t.Fatalf("Expected error %q, got %q", "config value is not a bool: 123", err.Error())
	}
}
