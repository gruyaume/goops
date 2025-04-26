package metadata_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruyaume/goops/environment"
	"github.com/gruyaume/goops/metadata"
)

const yamlContent = `
description: "An example charm"
name: example
summary: "Just a test"
provides:
  certificates:
    interface: tls-certificates
containers:
  web:
    resource: docker-image
    mounts:
      - location: /data
        storage: fast
resources:
  db:
    description: "database relation"
    type: relation
    upstream-source: mysql
storage:
  logs:
    minimum-size: 1G
    type: filesystem
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
	if err := os.WriteFile(metadataPath, []byte(yamlContent), 0o644); err != nil {
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

	if meta.Name != "example" {
		t.Errorf("Name = %q; want %q", meta.Name, "example")
	}

	if meta.Description != "An example charm" {
		t.Errorf("Description = %q; want %q", meta.Description, "An example charm")
	}

	if meta.Summary != "Just a test" {
		t.Errorf("Summary = %q; want %q", meta.Summary, "Just a test")
	}

	if prov, ok := meta.Provides["certificates"]; !ok {
		t.Errorf("Provides[\"certificates\"] missing")
	} else if prov.Interface != "tls-certificates" {
		t.Errorf("Provides[\"certificates\"].Interface = %q; want %q", prov.Interface, "tls-certificates")
	}

	cont, ok := meta.Containers["web"]
	if !ok {
		t.Fatalf("Containers[\"web\"] missing")
	}

	if cont.Resource != "docker-image" {
		t.Errorf("Containers[\"web\"].Resource = %q; want %q", cont.Resource, "docker-image")
	}

	if len(cont.Mounts) != 1 {
		t.Fatalf("expected 1 mount, got %d", len(cont.Mounts))
	}

	if m := cont.Mounts[0]; m.Location != "/data" || m.Storage != "fast" {
		t.Errorf("Mount = %+v; want Location=%q, Storage=%q", m, "/data", "fast")
	}

	res, ok := meta.Resources["db"]
	if !ok {
		t.Errorf("Resources[\"db\"] missing")
	} else {
		if res.Description != "database relation" {
			t.Errorf("Resources[\"db\"].Description = %q; want %q", res.Description, "database relation")
		}

		if res.Type != "relation" {
			t.Errorf("Resources[\"db\"].Type = %q; want %q", res.Type, "relation")
		}

		if res.UpstreamSource != "mysql" {
			t.Errorf("Resources[\"db\"].UpstreamSource = %q; want %q", res.UpstreamSource, "mysql")
		}
	}

	st, ok := meta.Storage["logs"]
	if !ok {
		t.Errorf("Storage[\"logs\"] missing")
	} else {
		if st.MinimumSize != "1G" {
			t.Errorf("Storage[\"logs\"].MinimumSize = %q; want %q", st.MinimumSize, "1G")
		}

		if st.Type != "filesystem" {
			t.Errorf("Storage[\"logs\"].Type = %q; want %q", st.Type, "filesystem")
		}
	}
}
