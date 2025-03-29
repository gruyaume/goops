package commands

import (
	"encoding/json"
	"fmt"
)

const (
	ConfigGetCommand = "config-get"
)

func ConfigGet(runner CommandRunner, key string) (string, error) {
	args := []string{key, "--format=json"}
	output, err := runner.Run(ConfigGetCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get config: %w", err)
	}
	var configValue string
	err = json.Unmarshal(output, &configValue)
	if err != nil {
		return "", fmt.Errorf("failed to parse config value: %w", err)
	}
	return configValue, nil
}
