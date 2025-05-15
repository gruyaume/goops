package metadata

import (
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

func GetCharmMetadata(path string) *Metadata {
	yamlFile, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return nil
	}

	var c Metadata

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil
	}

	return &c
}
