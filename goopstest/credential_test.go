package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func GetCredential() error {
	value, err := goops.GetCredential()
	if err != nil {
		return err
	}

	if len(value) != 0 {
		return fmt.Errorf("expected empty credential, got %v", value)
	}

	return nil
}

func TestGetCredential(t *testing.T) {
	ctx := goopstest.NewContext(GetCredential)

	stateIn := goopstest.State{}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no error, got: %v", ctx.CharmErr)
	}
}
