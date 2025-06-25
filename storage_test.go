package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestStorageAdd_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.AddStorage("database-storage", 1)
	if err != nil {
		t.Fatalf("StorageAdd returned an error: %v", err)
	}

	if fakeRunner.Command != "storage-add" {
		t.Errorf("Expected command %q, got %q", "storage-add", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "database-storage=1" {
		t.Errorf("Expected argument %q, got %q", "database-storage=1", fakeRunner.Args[0])
	}
}

func TestStorageAddWithCount_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.AddStorage("database-storage", 2)
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

func TestStorageGetByID_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"kind":"filesystem","location":"/var/lib/juju/storage/config/0"}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	storage, err := goops.GetStorageByID("database-storage")
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

	if storage.Kind != "filesystem" {
		t.Errorf("Expected storage %q, got %q", "filesystem", storage.Kind)
	}

	if storage.Location != "/var/lib/juju/storage/config/0" {
		t.Errorf("Expected storage location %q, got %q", "/var/lib/juju/storage/config/0", storage.Location)
	}
}

func TestStorageList_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["database-storage/0","database-storage/1"]`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	storage, err := goops.ListStorage("database-storage")
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
