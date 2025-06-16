package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

type MyConfig struct {
	Color    string `json:"color"`
	Quantity int    `json:"quantity"`
	ForSale  bool   `json:"for_sale"`
}

func TestGetConfigSuccess(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"color": "red", "quantity": 42, "for_sale": true}`),
		Err:    nil,
	}

	goops.SetRunner(fakeRunner)

	var dat MyConfig

	err := goops.GetConfig(&dat)
	if err != nil {
		t.Fatalf("Couldn't get config options: %v", err)
	}

	if dat.Color != "red" {
		t.Fatalf("Expected color 'red', got '%s'", dat.Color)
	}

	if dat.Quantity != 42 {
		t.Fatalf("Expected quantity 42, got %d", dat.Quantity)
	}

	if !dat.ForSale {
		t.Fatalf("Expected for_sale to be true, got %v", dat.ForSale)
	}
}

func TestGetConfigFailure(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"config not found"`),
		Err:    nil,
	}

	goops.SetRunner(fakeRunner)

	var dat MyConfig

	err := goops.GetConfig(&dat)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}
