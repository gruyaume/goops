package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func GetRelationIDs() error {
	relationIDs, err := goops.GetRelationIDs("certificates")
	if err != nil {
		return err
	}

	if len(relationIDs) != 1 {
		return fmt.Errorf("expected 1 relation ID, got %d", len(relationIDs))
	}

	if relationIDs[0] != "certificates:0" {
		return fmt.Errorf("expected relation ID 'certificates:0', got '%s'", relationIDs[0])
	}

	return nil
}

func TestCharmGetRelationIDs(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRelationIDs,
	}

	certRelation := &goopstest.Relation{
		Endpoint: "certificates",
	}
	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			certRelation,
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
}
