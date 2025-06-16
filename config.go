package goops

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	configGetCommand = "config-get"
)

var ErrConfigNotSet = errors.New("config option not set")

func GetConfig(config any) error {
	commandRunner := GetRunner()

	args := []string{"--format=json"}

	output, err := commandRunner.Run(configGetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if len(output) == 0 {
		return ErrConfigNotSet
	}

	err = json.Unmarshal(output, config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}
