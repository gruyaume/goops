package metadata_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruyaume/goops/environment"
	"github.com/gruyaume/goops/metadata"
)

const exampleMetadata = `
containers:
  notary:
    mounts:
    - location: /etc/notary/config
      storage: config
    - location: /var/lib/notary/database
      storage: database
    resource: notary-image
description: An example charm
name: notary-k8s
provides:
  certificates:
    interface: tls-certificates
resources:
  notary-image:
    description: OCI image for the Notary application
    type: oci-image
    upstream-source: ghcr.io/canonical/notary:0.0.3
storage:
  config:
    minimum-size: 5M
    type: filesystem
  database:
    minimum-size: 1G
    type: filesystem
summary: Certificate management made easy
`

type TestExecutionEnvironment struct {
	CharmDir string
}

func (r *TestExecutionEnvironment) Get(name string) string {
	return r.CharmDir
}

func NewTestExecutionEnvironment(charmDir string) *TestExecutionEnvironment {
	return &TestExecutionEnvironment{
		CharmDir: charmDir,
	}
}

func TestGetCharmMetadata_Success(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("JUJU_CHARM_DIR", tmpDir)

	metadataPath := filepath.Join(tmpDir, "metadata.yaml")
	if err := os.WriteFile(metadataPath, []byte(exampleMetadata), 0o644); err != nil {
		t.Fatalf("failed to write metadata.yaml: %v", err)
	}

	env := &environment.Environment{
		Getter: &TestExecutionEnvironment{
			CharmDir: tmpDir,
		},
	}

	meta, err := metadata.GetCharmMetadata(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Name != "notary-k8s" {
		t.Errorf("Name = %q; want %q", meta.Name, "notary-k8s")
	}

	if meta.Description != "An example charm" {
		t.Errorf("Description = %q; want %q", meta.Description, "An example charm")
	}

	if meta.Summary != "Certificate management made easy" {
		t.Errorf("Summary = %q; want %q", meta.Summary, "Certificate management made easy")
	}

	if prov, ok := meta.Provides["certificates"]; !ok {
		t.Errorf("Provides[\"certificates\"] missing")
	} else if prov.Interface != "tls-certificates" {
		t.Errorf("Provides[\"certificates\"].Interface = %q; want %q", prov.Interface, "tls-certificates")
	}

	cont, ok := meta.Containers["notary"]
	if !ok {
		t.Fatalf("Containers[\"notary\"] missing")
	}

	if cont.Resource != "notary-image" {
		t.Errorf("Containers[\"notary\"].Resource = %q; want %q", cont.Resource, "notary-image")
	}

	if len(cont.Mounts) != 2 {
		t.Fatalf("expected 2 mounts, got %d", len(cont.Mounts))
	}

	if m := cont.Mounts[0]; m.Location != "/etc/notary/config" || m.Storage != "config" {
		t.Errorf("Mount = %+v; want Location=%q, Storage=%q", m, "/etc/notary/config", "config")
	}

	if m := cont.Mounts[1]; m.Location != "/var/lib/notary/database" || m.Storage != "database" {
		t.Errorf("Mount = %+v; want Location=%q, Storage=%q", m, "/var/lib/notary/database", "database")
	}

	res, ok := meta.Resources["notary-image"]
	if !ok {
		t.Errorf("Resources[\"notary-image\"] missing")
	} else {
		if res.Description != "OCI image for the Notary application" {
			t.Errorf("Resources[\"notary-image\"].Description = %q; want %q", res.Description, "OCI image for the Notary application")
		}

		if res.Type != "oci-image" {
			t.Errorf("Resources[\"notary-image\"].Type = %q; want %q", res.Type, "oci-image")
		}

		if res.UpstreamSource != "ghcr.io/canonical/notary:0.0.3" {
			t.Errorf("Resources[\"notary-image\"].UpstreamSource = %q; want %q", res.UpstreamSource, "ghcr.io/canonical/notary:0.0.3")
		}
	}

	st, ok := meta.Storage["config"]
	if !ok {
		t.Errorf("Storage[\"config\"] missing")
	} else {
		if st.MinimumSize != "5M" {
			t.Errorf("Storage[\"config\"].MinimumSize = %q; want %q", st.MinimumSize, "5M")
		}

		if st.Type != "filesystem" {
			t.Errorf("Storage[\"config\"].Type = %q; want %q", st.Type, "filesystem")
		}
	}

	st, ok = meta.Storage["database"]
	if !ok {
		t.Errorf("Storage[\"database\"] missing")
	} else {
		if st.MinimumSize != "1G" {
			t.Errorf("Storage[\"database\"].MinimumSize = %q; want %q", st.MinimumSize, "1G")
		}

		if st.Type != "filesystem" {
			t.Errorf("Storage[\"database\"].Type = %q; want %q", st.Type, "filesystem")
		}
	}
}
