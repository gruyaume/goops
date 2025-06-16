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

func TestBla(t *testing.T) {
	var dat MyConfig

	if err := goops.GetConfig(&dat); err != nil {
		panic(err)
	}

	if dat.Color != "red" {
		t.Fatalf("Expected color 'red', got '%s'", dat.Color)
	}
}
