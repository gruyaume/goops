package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestStorageAdd_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	storageAddOpts := &commands.StorageAddOptions{
		Name: "database-storage",
	}

	err := command.StorageAdd(storageAddOpts)
	if err != nil {
		t.Fatalf("StorageAdd returned an error: %v", err)
	}

	if fakeRunner.Command != "storage-add" {
		t.Errorf("Expected command %q, got %q", "storage-add", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "database-storage" {
		t.Errorf("Expected argument %q, got %q", "database-storage", fakeRunner.Args[0])
	}
}

func TestStorageAddWithCount_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	storageAddOpts := &commands.StorageAddOptions{
		Name:  "database-storage",
		Count: 2,
	}

	err := command.StorageAdd(storageAddOpts)
	if err != nil {
		t.Fatalf("StorageAdd returned an error: %v", err)
	}

	if fakeRunner.Command != "storage-add" {
		t.Errorf("Expected command %q, got %q", "storage-add", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "database-storage=2" {
		t.Errorf("Expected argument %q, got %q", "database-storage=2", fakeRunner.Args[0])
	}
}

func TestStorageGetByName_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"database-storage"`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	storageGetOpts := &commands.StorageGetOptions{
		Name: "database-storage",
	}

	storage, err := command.StorageGet(storageGetOpts)
	if err != nil {
		t.Fatalf("StorageGet returned an error: %v", err)
	}

	if fakeRunner.Command != "storage-get" {
		t.Errorf("Expected command %q, got %q", "storage-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 3 {
		t.Fatalf("Expected 3 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "-s" {
		t.Errorf("Expected argument %q, got %q", "-s", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "database-storage" {
		t.Errorf("Expected argument %q, got %q", "database-storage", fakeRunner.Args[1])
	}

	if fakeRunner.Args[2] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[2])
	}

	if storage != "database-storage" {
		t.Errorf("Expected storage %q, got %q", "database-storage", storage)
	}
}

func TestStorageGetByID_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"database-storage"`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	storageGetOpts := &commands.StorageGetOptions{
		ID: "21127934-8986-11e5-af63-feff819cdc9f",
	}

	storage, err := command.StorageGet(storageGetOpts)
	if err != nil {
		t.Fatalf("StorageGet returned an error: %v", err)
	}

	if fakeRunner.Command != "storage-get" {
		t.Errorf("Expected command %q, got %q", "storage-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 3 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "21127934-8986-11e5-af63-feff819cdc9f" {
		t.Errorf("Expected argument %q, got %q", "21127934-8986-11e5-af63-feff819cdc9f", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}

	if storage != "database-storage" {
		t.Errorf("Expected storage %q, got %q", "database-storage", storage)
	}
}

func TestStorageList_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["database-storage/0","database-storage/1"]`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	storageListOpts := &commands.StorageListOptions{
		Name: "database-storage",
	}

	storage, err := command.StorageList(storageListOpts)
	if err != nil {
		t.Fatalf("StorageList returned an error: %v", err)
	}

	if fakeRunner.Command != "storage-list" {
		t.Errorf("Expected command %q, got %q", "storage-list", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "database-storage" {
		t.Errorf("Expected argument %q, got %q", "database-storage", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}

	if len(storage) != 2 {
		t.Errorf("Expected 2 storage items, got %d", len(storage))
	}

	if storage[0] != "database-storage/0" {
		t.Errorf("Expected storage item %q, got %q", "database-storage/0", storage[0])
	}

	if storage[1] != "database-storage/1" {
		t.Errorf("Expected storage item %q, got %q", "database-storage/1", storage[1])
	}
}
