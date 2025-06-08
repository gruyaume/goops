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

func ListRelations1Result() error {
	relationUnits, err := goops.ListRelations("certificates:0")
	if err != nil {
		return err
	}

	if len(relationUnits) != 1 {
		return fmt.Errorf("expected 1 relation unit, got %d", len(relationUnits))
	}

	if relationUnits[0] != "provider/0" {
		return fmt.Errorf("expected relation unit 'provider/0', got '%s'", relationUnits[0])
	}

	return nil
}

func ListRelations3Result() error {
	relationUnits, err := goops.ListRelations("certificates:0")
	if err != nil {
		return err
	}

	expected := map[string]bool{
		"provider/0": false,
		"provider/1": false,
		"provider/2": false,
	}

	if len(relationUnits) != len(expected) {
		return fmt.Errorf("expected %d relation units, got %d", len(expected), len(relationUnits))
	}

	for _, unit := range relationUnits {
		if _, ok := expected[unit]; !ok {
			return fmt.Errorf("unexpected relation unit: %s", unit)
		}

		expected[unit] = true
	}

	for unit, found := range expected {
		if !found {
			return fmt.Errorf("missing expected relation unit: %s", unit)
		}
	}

	return nil
}

func TestCharmListRelations(t *testing.T) {
	tests := []struct {
		name        string
		handler     func() error
		remoteUnits int
	}{
		{
			name:        "ListRelations1Result",
			handler:     ListRelations1Result,
			remoteUnits: 1,
		},
		{
			name:        "ListRelations3Result",
			handler:     ListRelations3Result,
			remoteUnits: 3,
		},
	}

	for _, tc := range tests {
		ctx := goopstest.Context{
			Charm: tc.handler,
		}

		remoteUnitsData := map[goopstest.UnitID]goopstest.DataBag{}

		for i := 0; i < tc.remoteUnits; i++ {
			unitID := fmt.Sprintf("provider/%d", i)
			remoteUnitsData[goopstest.UnitID(unitID)] = goopstest.DataBag{}
		}

		certRelation := &goopstest.Relation{
			Endpoint:        "certificates",
			RemoteAppName:   "provider",
			RemoteUnitsData: remoteUnitsData,
		}
		stateIn := &goopstest.State{
			Relations: []*goopstest.Relation{
				certRelation,
			},
		}

		stateOut, err := ctx.Run("start", stateIn)
		if err != nil {
			t.Fatalf("Run returned an error: %v", err)
		}

		if len(stateOut.Relations) != 1 {
			t.Fatalf("expected 1 relation, got %d", len(stateOut.Relations))
		}
	}
}

func GetUnitRelationData() error {
	relationData, err := goops.GetUnitRelationData("certificates:0", "provider/0")
	if err != nil {
		return err
	}

	if len(relationData) == 0 {
		return fmt.Errorf("expected relation data, got empty map")
	}

	csr, ok := relationData["certificate_signing_requests"]
	if !ok {
		return fmt.Errorf("expected 'certificate_signing_requests' key in relation data")
	}

	if csr != "csr-data" {
		return fmt.Errorf("expected 'csr-data', got '%s'", csr)
	}

	return nil
}

func TestCharmGetUnitRelationData(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetUnitRelationData,
	}

	remoteUnitsData := map[goopstest.UnitID]goopstest.DataBag{
		goopstest.UnitID("provider/0"): {
			"certificate_signing_requests": "csr-data",
		},
	}

	certRelation := &goopstest.Relation{
		Endpoint:        "certificates",
		RemoteAppName:   "provider",
		RemoteUnitsData: remoteUnitsData,
	}
	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			certRelation,
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(stateOut.Relations))
	}
}

func GetAppRelationData() error {
	relationData, err := goops.GetAppRelationData("certificates:0", "provider/0")
	if err != nil {
		return err
	}

	if len(relationData) == 0 {
		return fmt.Errorf("expected relation data, got empty map")
	}

	csr, ok := relationData["certificate_signing_requests"]
	if !ok {
		return fmt.Errorf("expected 'certificate_signing_requests' key in relation data")
	}

	if csr != "csr-data" {
		return fmt.Errorf("expected 'csr-data', got '%s'", csr)
	}

	return nil
}

func TestCharmGetAppRelationData(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetAppRelationData,
	}

	remoteAppData := goopstest.DataBag{
		"certificate_signing_requests": "csr-data",
	}

	certRelation := &goopstest.Relation{
		Endpoint:      "certificates",
		RemoteAppName: "provider",
		RemoteAppData: remoteAppData,
	}
	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			certRelation,
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(stateOut.Relations))
	}
}

func SetUnitRelationData() error {
	relationData := map[string]string{
		"certificate_signing_requests": "csr-data",
	}

	err := goops.SetUnitRelationData("certificates:0", relationData)
	if err != nil {
		return err
	}

	return nil
}

func TestCharmSetUnitRelationData(t *testing.T) {
	ctx := goopstest.Context{
		Charm: SetUnitRelationData,
	}

	certRelation := &goopstest.Relation{
		Endpoint: "certificates",
	}
	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			certRelation,
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(stateOut.Relations))
	}

	localUnitData := stateOut.Relations[0].LocalUnitData
	if len(localUnitData) != 1 {
		t.Fatalf("expected 1 local unit data, got %d", len(localUnitData))
	}

	if localUnitData["certificate_signing_requests"] != "csr-data" {
		t.Fatalf("expected 'csr-data', got '%s'", localUnitData["certificate_signing_requests"])
	}
}

func SetAppRelationData() error {
	relationData := map[string]string{
		"certificate_signing_requests": "csr-data",
	}

	err := goops.SetAppRelationData("certificates:0", relationData)
	if err != nil {
		return err
	}

	return nil
}

func TestCharmSetAppRelationData(t *testing.T) {
	ctx := goopstest.Context{
		Charm: SetAppRelationData,
	}

	certRelation := &goopstest.Relation{
		Endpoint: "certificates",
	}
	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			certRelation,
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(stateOut.Relations))
	}

	appData := stateOut.Relations[0].LocalAppData
	if len(appData) != 1 {
		t.Fatalf("expected 1 app data, got %d", len(appData))
	}

	if appData["certificate_signing_requests"] != "csr-data" {
		t.Fatalf("expected 'csr-data', got '%s'", appData["certificate_signing_requests"])
	}
}

func RelationEndToEnd() error {
	relationName := "certificates"

	relationIDs, err := goops.GetRelationIDs(relationName)
	if err != nil {
		return fmt.Errorf("could not get relation IDs: %w", err)
	}

	requirerCertificateRequests := make([]string, 0)

	for _, relationID := range relationIDs {
		relationUnits, err := goops.ListRelations(relationID)
		if err != nil {
			return fmt.Errorf("could not list relation data: %w", err)
		}

		for _, unitID := range relationUnits {
			relationData, err := goops.GetUnitRelationData(relationID, unitID)
			if err != nil {
				return fmt.Errorf("could not get relation data: %w", err)
			}

			csr, ok := relationData["certificate_signing_requests"]
			if !ok {
				continue
			}

			requirerCertificateRequests = append(requirerCertificateRequests, csr)
		}
	}

	if len(requirerCertificateRequests) == 0 {
		return fmt.Errorf("no certificate signing requests found in relation data")
	}

	return nil
}

func TestCharmRelationEndToEnd(t *testing.T) {
	ctx := goopstest.Context{
		Charm: RelationEndToEnd,
	}

	remoteUnitsData := map[goopstest.UnitID]goopstest.DataBag{
		goopstest.UnitID("provider/0"): {
			"certificate_signing_requests": "csr-data-0",
		},
		goopstest.UnitID("provider/1"): {
			"certificate_signing_requests": "csr-data-1",
		},
	}

	certRelation := &goopstest.Relation{
		Endpoint:        "certificates",
		RemoteAppName:   "provider",
		RemoteUnitsData: remoteUnitsData,
	}
	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			certRelation,
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(stateOut.Relations))
	}
}
