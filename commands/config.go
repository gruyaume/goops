package commands

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	configGetCommand = "config-get"
)

var ErrConfigNotSet = errors.New("config option not set")

type ConfigGetOptions struct {
	Key string
}

func (command Command) ConfigGet(opts *ConfigGetOptions) (any, error) {
	args := []string{opts.Key, "--format=json"}

	output, err := command.Runner.Run(configGetCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get config: %w", err)
	}

	var configValue any

	err = json.Unmarshal(output, &configValue)
	if err != nil {
		return "", fmt.Errorf("failed to parse config value: %w", err)
	}

	return configValue, nil
}

func (command Command) ConfigGetString(opts *ConfigGetOptions) (string, error) {
	value, err := command.ConfigGet(opts)
	if err != nil {
		return "", err
	}

	if value == nil {
		return "", ErrConfigNotSet
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("config value is not a string: %v", value)
	}

	return strValue, nil
}

func (command Command) ConfigGetInt(opts *ConfigGetOptions) (int, error) {
	value, err := command.ConfigGet(opts)
	if err != nil {
		return 0, err
	}

	if value == nil {
		return 0, ErrConfigNotSet
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("config value is not a number: %v", value)
	}

	return int(floatValue), nil
}

func (command Command) ConfigGetBool(opts *ConfigGetOptions) (bool, error) {
	value, err := command.ConfigGet(opts)
	if err != nil {
		return false, err
	}

	if value == nil {
		return false, ErrConfigNotSet
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("config value is not a bool: %v", value)
	}

	return boolValue, nil
}
