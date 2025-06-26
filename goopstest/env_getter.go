package goopstest

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type fakeEnvGetter struct {
	HookName    string
	ActionName  string
	Model       Model
	AppName     string
	UnitID      string
	JujuVersion string
	Metadata    Metadata
}

func (f *fakeEnvGetter) Get(key string) string {
	switch key {
	case "JUJU_HOOK_NAME":
		return f.HookName
	case "JUJU_ACTION_NAME":
		return f.ActionName
	case "JUJU_MODEL_NAME":
		return f.Model.Name
	case "JUJU_MODEL_UUID":
		return f.Model.UUID
	case "JUJU_UNIT_NAME":
		return f.UnitID
	case "JUJU_VERSION":
		return f.JujuVersion
	}

	return ""
}

func (f *fakeEnvGetter) ReadFile(name string) ([]byte, error) {
	if strings.HasSuffix(name, "metadata.yaml") {
		data, err := yaml.Marshal(f.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}

		return data, nil
	}

	return nil, fmt.Errorf("file %s not found", name)
}
