package goopstest_test

import (
	"encoding/json"
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

func GetRelationIDsNoRelation() error {
	relationIDs, err := goops.GetRelationIDs("certificates")
	if err != nil {
		return err
	}

	if len(relationIDs) != 0 {
		return fmt.Errorf("expected no relation IDs, got %d", len(relationIDs))
	}

	return nil
}

func TestCharmGetRelationIDsNoRelation(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRelationIDsNoRelation,
	}

	stateIn := &goopstest.State{}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Relations) != 0 {
		t.Fatalf("expected no relations, got %d", len(stateOut.Relations))
	}
}

func GetRelationIDsNoName() error {
	_, err := goops.GetRelationIDs("")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmGetRelationIDsNoName(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRelationIDsNoName,
	}

	stateIn := &goopstest.State{}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	expectedErr := "failed to get relation IDs: command relation-ids failed: ERROR no endpoint name specified"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
	}
}

func GetRelationIDsNoResult() error {
	relationIDs, err := goops.GetRelationIDs("nonexistent")
	if err != nil {
		return err
	}

	if len(relationIDs) != 0 {
		return fmt.Errorf("expected no relation IDs, got %d", len(relationIDs))
	}

	return nil
}

func TestCharmGetRelationIDsNoResult(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRelationIDsNoResult,
	}

	stateIn := &goopstest.State{}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatal("Expected no charm error, got one")
	}
}

func ListRelationUnits1Result() error {
	relationUnits, err := goops.ListRelationUnits("certificates:0")
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

func ListRelationUnits3Result() error {
	relationUnits, err := goops.ListRelationUnits("certificates:0")
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

func TestCharmListRelationUnits(t *testing.T) {
	tests := []struct {
		name        string
		handler     func() error
		remoteUnits int
	}{
		{
			name:        "ListRelationUnits1Result",
			handler:     ListRelationUnits1Result,
			remoteUnits: 1,
		},
		{
			name:        "ListRelationUnits3Result",
			handler:     ListRelationUnits3Result,
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

func TestListRelationUnitsResultNotFound(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ListRelationUnits1Result,
	}

	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	expectedErr := "failed to list relation data: command relation-list failed: ERROR invalid value \"certificates:0\" for option -r: relation not found"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
	}
}

func TestListRelationUnitsInActionHook(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ListRelationUnits1Result,
	}

	remoteUnitsData := map[goopstest.UnitID]goopstest.DataBag{
		goopstest.UnitID("provider/0"): {},
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

	stateOut, err := ctx.RunAction("run-action", stateIn, nil)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(stateOut.Relations))
	}
}

func GetRemoteUnitRelationData() error {
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

func TestCharmGetRemoteUnitRelationData(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRemoteUnitRelationData,
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

func TestCharmGetRemoteUnitRelationDataNoRelation(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRemoteUnitRelationData,
	}

	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	expectedErr := "failed to get relation data: command relation-get failed: ERROR invalid value \"certificates:0\" for option -r: relation not found"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
	}
}

func TestCharmGetRemoteUnitRelationDataNoRemoteUnit(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRemoteUnitRelationData,
	}

	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			{
				Endpoint:      "certificates",
				RemoteAppName: "provider",
				RemoteUnitsData: map[goopstest.UnitID]goopstest.DataBag{
					goopstest.UnitID("certificates/22"): {
						"certificate_signing_requests": "csr-data",
					},
				},
			},
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	expectedErr := "failed to get relation data: command relation-get failed: ERROR cannot read settings for unit \"provider/0\" in relation \"certificates:0\": unit \"provider/0\": settings not found"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
	}
}

func TestCharmGetRemoteUnitRelationDataNoData(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRemoteUnitRelationData,
	}

	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			{
				Endpoint:      "certificates",
				RemoteAppName: "provider",
				RemoteUnitsData: map[goopstest.UnitID]goopstest.DataBag{
					goopstest.UnitID("provider/0"): {},
				},
			},
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("Expected CharmErr to be set, got nil")
	}

	expectedErr := "expected relation data, got empty map"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
	}
}

func GetLocalUnitRelationData() error {
	relationData, err := goops.GetUnitRelationData("certificates:0", "requirer/0")
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

func TestCharmGetLocalUnitRelationData(t *testing.T) {
	ctx := goopstest.Context{
		Charm:   GetLocalUnitRelationData,
		AppName: "requirer",
		UnitID:  "requirer/0",
	}

	certRelation := &goopstest.Relation{
		Endpoint: "certificates",
		LocalUnitData: map[string]string{
			"certificate_signing_requests": "csr-data",
		},
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

// Each unit can only read its own databag. Reading another local unit's data should return nothing.
func TestGetOtherLocalUnitRelationData(t *testing.T) {
	ctx := goopstest.Context{
		Charm:   GetLocalUnitRelationData,
		AppName: "requirer",
		UnitID:  "requirer/1", // This unit should not be able to read data from unit 0
	}

	certRelation := &goopstest.Relation{
		Endpoint: "certificates",
		LocalUnitData: map[string]string{
			"certificate_signing_requests": "csr-data",
		},
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

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	expectedErr := "expected relation data, got empty map"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
	}
}

func GetRemoteAppRelationData() error {
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

func TestCharmGetRemoteAppRelationData(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRemoteAppRelationData,
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

func TestCharmGetRemoteAppRelationDataNoRelation(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetRemoteAppRelationData,
	}

	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	expectedErr := "failed to get relation data: command relation-get failed: ERROR invalid value \"certificates:0\" for option -r: relation not found"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
	}
}

func GetLocalAppRelationData() error {
	relationData, err := goops.GetAppRelationData("certificates:0", "requirer/0")
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

func TestCharmGetLocalAppRelationData(t *testing.T) {
	ctx := goopstest.Context{
		Charm:   GetLocalAppRelationData,
		AppName: "requirer",
		UnitID:  "requirer/0",
	}

	certRelation := &goopstest.Relation{
		Endpoint: "certificates",
		LocalAppData: map[string]string{
			"certificate_signing_requests": "csr-data",
		},
	}
	stateIn := &goopstest.State{
		Leader: true,
		Relations: []*goopstest.Relation{
			certRelation,
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no CharmErr, got %v", ctx.CharmErr)
	}

	if len(stateOut.Relations) != 1 {
		t.Fatalf("expected 1 relation, got %d", len(stateOut.Relations))
	}
}

// Only leader units can read to the local application databag;
func TestCharmGetLocalAppRelationDataNonLeader(t *testing.T) {
	ctx := goopstest.Context{
		Charm:   GetLocalAppRelationData,
		AppName: "requirer",
		UnitID:  "requirer/0",
	}

	certRelation := &goopstest.Relation{
		Endpoint: "certificates",
		LocalAppData: map[string]string{
			"certificate_signing_requests": "csr-data",
		},
	}
	stateIn := &goopstest.State{
		Leader: false,
		Relations: []*goopstest.Relation{
			certRelation,
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("Expected CharmErr to be set, got nil")
	}

	expectedErr := "failed to get relation data: command relation-get failed: ERROR permission denied"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
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

type TLSConfig struct {
	InsecureSkipVerify bool `json:"insecure_skip_verify"`
}

type StaticConfig struct {
	Targets []string `json:"targets"`
}

type Job struct {
	Scheme        string         `json:"scheme"`
	TLSConfig     TLSConfig      `json:"tls_config"`
	MetricsPath   string         `json:"metrics_path"`
	StaticConfigs []StaticConfig `json:"static_configs"`
}

type ScrapeMetadata struct {
	Model       string `json:"model"`
	ModelUUID   string `json:"model_uuid"`
	Application string `json:"application"`
	Unit        string `json:"unit"`
	CharmName   string `json:"charm_name"`
}

func SetAppRelationData2() error {
	relationIDs, err := goops.GetRelationIDs("metrics")
	if err != nil {
		return fmt.Errorf("could not get relation IDs: %w", err)
	}

	jobs := []*Job{
		{
			Scheme:      "https",
			TLSConfig:   TLSConfig{InsecureSkipVerify: true},
			MetricsPath: "/metrics",
			StaticConfigs: []StaticConfig{
				{
					Targets: []string{"localhost:8080"},
				},
			},
		},
	}

	scrapeJobs, err := json.Marshal(jobs)
	if err != nil {
		return fmt.Errorf("could not marshal scrape jobs to JSON: %w", err)
	}

	scrapeMetadata := &ScrapeMetadata{
		Model:       "test-model",
		ModelUUID:   "12345678-1234-5678-1234-567812345678",
		Application: "test-application",
		Unit:        "test-unit/0",
		CharmName:   "test-charm",
	}

	scrapeMetadataBytes, err := json.Marshal(scrapeMetadata)
	if err != nil {
		return fmt.Errorf("could not marshal scrape metadata to JSON: %w", err)
	}

	relationData := map[string]string{
		"scrape_jobs":     string(scrapeJobs),
		"scrape_metadata": string(scrapeMetadataBytes),
	}

	err = goops.SetAppRelationData(relationIDs[0], relationData)
	if err != nil {
		return fmt.Errorf("could not set relation data: %w", err)
	}

	return nil
}

func TestCharmSetAppRelationData2(t *testing.T) {
	ctx := goopstest.Context{
		Charm: SetAppRelationData2,
	}

	certRelation := &goopstest.Relation{
		Endpoint: "metrics",
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
	if len(appData) != 2 {
		t.Fatalf("expected 2 app data, got %d", len(appData))
	}

	if _, ok := appData["scrape_jobs"]; !ok {
		t.Fatal("expected 'scrape_jobs' key in app data")
	}

	if _, ok := appData["scrape_metadata"]; !ok {
		t.Fatal("expected 'scrape_metadata' key in app data")
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
		relationUnits, err := goops.ListRelationUnits(relationID)
		if err != nil {
			return fmt.Errorf("could not list relation units: %w", err)
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

func RelationModelGetUUID() error {
	modelUUID, err := goops.GetRelationModelUUID("certificates:0")
	if err != nil {
		return err
	}

	if modelUUID != "a4e65ff5-2358-4595-8ace-cc820c120e24" {
		return fmt.Errorf("expected model UUID 'a4e65ff5-2358-4595-8ace-cc820c120e24', got '%s'", modelUUID)
	}

	return nil
}

func TestCharmRelationModelGetUUID(t *testing.T) {
	ctx := goopstest.Context{
		Charm: RelationModelGetUUID,
	}

	certRelation := &goopstest.Relation{
		Endpoint: "certificates",
	}
	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			certRelation,
		},
		Model: &goopstest.Model{
			UUID: "a4e65ff5-2358-4595-8ace-cc820c120e24",
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no CharmErr, got %v", ctx.CharmErr)
	}
}

func TestCharmRelationModelGetUUIDWithRemoteModelUUID(t *testing.T) {
	ctx := goopstest.Context{
		Charm: RelationModelGetUUID,
	}

	certRelation := &goopstest.Relation{
		Endpoint:        "certificates",
		RemoteModelUUID: "a4e65ff5-2358-4595-8ace-cc820c120e24",
	}
	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{
			certRelation,
		},
		Model: &goopstest.Model{
			UUID: "a-different-uuid",
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no CharmErr, got %v", ctx.CharmErr)
	}
}

func TestCharmRelationModelGetUUIDNoRelation(t *testing.T) {
	ctx := goopstest.Context{
		Charm: RelationModelGetUUID,
	}

	stateIn := &goopstest.State{
		Relations: []*goopstest.Relation{},
		Model: &goopstest.Model{
			UUID: "a4e65ff5-2358-4595-8ace-cc820c120e24",
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("expected CharmErr to be set, got nil")
	}

	expectedError := "failed to get relation model data: command relation-model-get failed: ERROR invalid value \"certificates:0\" for option -r: relation not found"
	if ctx.CharmErr.Error() != expectedError {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedError)
	}
}
