package goops

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Mount struct {
	Location string `yaml:"location"`
	Storage  string `yaml:"storage"`
}

type Container struct {
	Mounts   []Mount `yaml:"mounts"`
	Resource string  `yaml:"resource"`
}

type Integration struct {
	Interface string `yaml:"interface"`
}

type Resource struct {
	Description    string `yaml:"description"`
	Type           string `yaml:"type"`
	UpstreamSource string `yaml:"upstream-source"`
}

type Storage struct {
	MinimumSize string `yaml:"minimum-size"`
	Type        string `yaml:"type"`
}

type Metadata struct {
	Containers  map[string]Container   `yaml:"containers"`
	Description string                 `yaml:"description"`
	Name        string                 `yaml:"name"`
	Provides    map[string]Integration `yaml:"provides"`
	Resources   map[string]Resource    `yaml:"resources"`
	Storage     map[string]Storage     `yaml:"storage"`
	Summary     string                 `yaml:"summary"`
}

func ReadMetadata() (*Metadata, error) {
	env := ReadEnv()

	path := env.CharmDir + "/metadata.yaml"

	yamlFile, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var c Metadata

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &c, nil
}
