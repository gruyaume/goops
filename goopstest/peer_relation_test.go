package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func GetRelationIDsForPeers() error {
	relationIDs, err := goops.GetRelationIDs("example-peer")
	if err != nil {
		return fmt.Errorf("could not get peer relation ID: %w", err)
	}

	if len(relationIDs) != 1 {
		return fmt.Errorf("expected 1 peer relation ID, got %d", len(relationIDs))
	}

	if relationIDs[0] != "example-peer:0" {
		return fmt.Errorf("expected peer relation ID 'example-peer:0', got '%s'", relationIDs[0])
	}

	return nil
}

func TestGetRelationIDsForPeers(t *testing.T) {
	ctx := goopstest.Context{
		Charm:   GetRelationIDsForPeers,
		AppName: "example",
		UnitID:  "example/0",
	}

	peersData := map[goopstest.UnitID]goopstest.DataBag{
		goopstest.UnitID("example/1"): {},
	}

	peerRelation := &goopstest.PeerRelation{
		Endpoint:  "example-peer",
		Interface: "example-peer",
		ID:        "example-peer:0",
		PeersData: peersData,
	}

	stateIn := &goopstest.State{
		PeerRelations: []*goopstest.PeerRelation{
			peerRelation,
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Errorf("expected no error, got %v", ctx.CharmErr)
	}

	if len(stateOut.PeerRelations) != 1 {
		t.Errorf("expected 1 peer relation, got %d", len(stateOut.PeerRelations))
	}
}

func TestGetRelationIDsNoPeers(t *testing.T) {
	ctx := goopstest.Context{
		Charm:   GetRelationIDsForPeers,
		AppName: "example",
		UnitID:  "example/0",
	}

	stateIn := &goopstest.State{
		PeerRelations: []*goopstest.PeerRelation{},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Errorf("expected charm error to be set, got nil")
	}

	if ctx.CharmErr.Error() != "expected 1 peer relation ID, got 0" {
		t.Errorf("got CharmErr=%q, want 'expected 1 peer relation ID, got 0'", ctx.CharmErr.Error())
	}
}

func ListPeerRelationUnits() error {
	peerUnits, err := goops.ListRelationUnits("example-peer:0")
	if err != nil {
		return fmt.Errorf("could not list peer relation units: %w", err)
	}

	if len(peerUnits) != 1 {
		return fmt.Errorf("expected at least one peer unit, got none")
	}

	if peerUnits[0] != "example/1" {
		return fmt.Errorf("expected peer unit ID 'example/1', got '%s'", peerUnits[0])
	}

	return nil
}

func TestListPeerRelationUnits(t *testing.T) {
	ctx := goopstest.Context{
		Charm:   ListPeerRelationUnits,
		AppName: "example",
		UnitID:  "example/0",
	}

	stateIn := &goopstest.State{
		PeerRelations: []*goopstest.PeerRelation{
			{
				ID:        "example-peer:0",
				PeersData: map[goopstest.UnitID]goopstest.DataBag{"example/1": {}},
			},
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Errorf("expected no error, got %v", ctx.CharmErr)
	}

	if len(stateOut.PeerRelations) != 1 {
		t.Errorf("expected 1 peer relation, got %d", len(stateOut.PeerRelations))
	}
}
