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
	ctx := goopstest.NewContext(GetRelationIDsForPeers, goopstest.WithUnitID("example/0"))

	peersData := map[goopstest.UnitID]goopstest.DataBag{
		goopstest.UnitID("example/1"): {},
	}

	peerRelation := goopstest.PeerRelation{
		Endpoint:  "example-peer",
		Interface: "example-peer",
		ID:        "example-peer:0",
		PeersData: peersData,
	}

	stateIn := goopstest.State{
		PeerRelations: []goopstest.PeerRelation{
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
	ctx := goopstest.NewContext(GetRelationIDsForPeers, goopstest.WithUnitID("example/0"))

	stateIn := goopstest.State{
		PeerRelations: []goopstest.PeerRelation{},
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
	ctx := goopstest.NewContext(ListPeerRelationUnits, goopstest.WithUnitID("example/0"))

	stateIn := goopstest.State{
		PeerRelations: []goopstest.PeerRelation{
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

func GetPeerRelationModelUUID() error {
	modelUUID, err := goops.GetRelationModelUUID("example-peer:0")
	if err != nil {
		return fmt.Errorf("could not get peer data: %w", err)
	}

	if modelUUID != "example-model-uuid" {
		return fmt.Errorf("expected model UUID 'example-model-uuid', got '%s'", modelUUID)
	}

	return nil
}

func TestGetPeerRelationModelUUID(t *testing.T) {
	ctx := goopstest.NewContext(GetPeerRelationModelUUID, goopstest.WithUnitID("example/0"))

	stateIn := goopstest.State{
		PeerRelations: []goopstest.PeerRelation{
			{
				ID: "example-peer:0",
			},
		},
		Model: goopstest.Model{
			UUID: "example-model-uuid",
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

func GetUnitPeerRelationData() error {
	data, err := goops.GetUnitRelationData("example-peer:0", "example/0")
	if err != nil {
		return fmt.Errorf("could not get unit relation data: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("expected unit relation data, got none")
	}

	if data["key"] != "value" {
		return fmt.Errorf("expected unit relation data 'key=value', got '%s=%s'", "key", data["key"])
	}

	return nil
}

// For peer relations, each unit can read its own databag
func TestGetSelfUnitPeerRelationData(t *testing.T) {
	tests := []struct {
		leader bool
	}{
		{leader: true},
		{leader: false},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("Leader=%v", tc.leader), func(t *testing.T) {
			ctx := goopstest.NewContext(GetUnitPeerRelationData, goopstest.WithUnitID("example/0"))

			stateIn := goopstest.State{
				Leader: tc.leader,
				PeerRelations: []goopstest.PeerRelation{
					{
						ID: "example-peer:0",
						LocalUnitData: goopstest.DataBag{
							"key": "value",
						},
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
		})
	}
}

// For peer relations, each unit can read the other unit's databag
func TestGetOtherUnitPeerRelationData(t *testing.T) {
	tests := []struct {
		leader bool
	}{
		{leader: true},
		{leader: false},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("Leader=%v", tc.leader), func(t *testing.T) {
			ctx := goopstest.NewContext(GetUnitPeerRelationData, goopstest.WithUnitID("example/1"))

			stateIn := goopstest.State{
				Leader: tc.leader,
				PeerRelations: []goopstest.PeerRelation{
					{
						ID:            "example-peer:0",
						LocalUnitData: goopstest.DataBag{},
						PeersData: map[goopstest.UnitID]goopstest.DataBag{
							"example/0": {
								"key": "value",
							},
						},
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
		})
	}
}

func GetAppPeerRelationData() error {
	data, err := goops.GetAppRelationData("example-peer:0", "example-peer/0")
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("expected app relation data, got none")
	}

	if data["app_key"] != "app_value" {
		return fmt.Errorf("expected app relation data 'app_key=app_value', got '%s=%s'", "app_key", data["app_key"])
	}

	return nil
}

// Each unit can read the application's databag
func TestGetAppPeerRelationData(t *testing.T) {
	tests := []struct {
		leader bool
	}{
		{leader: true},
		{leader: false},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("Leader=%v", tc.leader), func(t *testing.T) {
			ctx := goopstest.NewContext(GetAppPeerRelationData, goopstest.WithUnitID("example-peer/0"))

			stateIn := goopstest.State{
				Leader: tc.leader,
				PeerRelations: []goopstest.PeerRelation{
					{
						ID: "example-peer:0",
						LocalAppData: goopstest.DataBag{
							"app_key": "app_value",
						},
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
		})
	}
}

func TestGetAppPeerRelationDataNoRelation(t *testing.T) {
	ctx := goopstest.NewContext(GetAppPeerRelationData, goopstest.WithUnitID("example-peer/0"))

	stateIn := goopstest.State{
		PeerRelations: []goopstest.PeerRelation{},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Errorf("expected charm error to be set, got nil")
	}

	expectedErr := "failed to get relation data: command relation-get failed: ERROR invalid value \"example-peer:0\" for option -r: relation not found"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want '%s'", ctx.CharmErr.Error(), expectedErr)
	}
}

func SetPeerUnitRelationData() error {
	relationData := map[string]string{
		"whatever_key": "whatever_value",
	}

	err := goops.SetUnitRelationData("example-peer:0", relationData)
	if err != nil {
		return err
	}

	return nil
}

func TestSetPeerUnitRelationData(t *testing.T) {
	ctx := goopstest.NewContext(SetPeerUnitRelationData, goopstest.WithUnitID("example/0"))

	stateIn := goopstest.State{
		PeerRelations: []goopstest.PeerRelation{
			{
				ID: "example-peer:0",
			},
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no error, got %v", ctx.CharmErr)
	}

	if len(stateOut.PeerRelations) != 1 {
		t.Errorf("expected 1 peer relation, got %d", len(stateOut.PeerRelations))
	}

	if stateOut.PeerRelations[0].LocalUnitData["whatever_key"] != "whatever_value" {
		t.Errorf("expected 'whatever_key=whatever_value', got '%s=%s'", "whatever_key", stateOut.PeerRelations[0].LocalUnitData["whatever_key"])
	}
}

func SetPeerAppRelationData() error {
	relationData := map[string]string{
		"app_key": "app_value",
	}

	err := goops.SetAppRelationData("example-peer:0", relationData)
	if err != nil {
		return err
	}

	return nil
}

func TestSetPeerAppRelationDataLeader(t *testing.T) {
	ctx := goopstest.NewContext(SetPeerAppRelationData, goopstest.WithUnitID("example-peer/0"))

	stateIn := goopstest.State{
		Leader: true,
		PeerRelations: []goopstest.PeerRelation{
			{
				ID: "example-peer:0",
			},
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no error, got %v", ctx.CharmErr)
	}

	if len(stateOut.PeerRelations) != 1 {
		t.Errorf("expected 1 peer relation, got %d", len(stateOut.PeerRelations))
	}

	if stateOut.PeerRelations[0].LocalAppData["app_key"] != "app_value" {
		t.Errorf("expected 'app_key=app_value', got '%s=%s'", "app_key", stateOut.PeerRelations[0].LocalAppData["app_key"])
	}
}

func TestSetPeerAppRelationDataNonLeader(t *testing.T) {
	ctx := goopstest.NewContext(SetPeerAppRelationData, goopstest.WithUnitID("example-peer/0"))

	stateIn := goopstest.State{
		Leader: false,
		PeerRelations: []goopstest.PeerRelation{
			{
				ID: "example-peer:0",
			},
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("expected charm error to be set, got nil")
	}

	expectedErr := "failed to set relation data: command relation-set failed: ERROR cannot write relation settings"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want '%s'", ctx.CharmErr.Error(), expectedErr)
	}

	if len(stateOut.PeerRelations) != 1 {
		t.Errorf("expected 1 peer relation, got %d", len(stateOut.PeerRelations))
	}

	if stateOut.PeerRelations[0].LocalAppData != nil {
		t.Errorf("expected no app relation data, got %v", stateOut.PeerRelations[0].LocalAppData)
	}
}
