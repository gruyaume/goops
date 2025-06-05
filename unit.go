package goops

import (
	"encoding/json"
	"fmt"
)

const (
	unitGetCommand = "unit-get"
)

func getUnit(key string) (string, error) {
	commandRunner := GetRunner()

	args := []string{}

	args = append(args, key)

	args = append(args, "--format=json")

	output, err := commandRunner.Run(unitGetCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get unit: %w", err)
	}

	var result string
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("failed to parse unit get output: %w", err)
	}

	return result, nil
}

func GetUnitPublicAddress() (string, error) {
	return getUnit("public-address")
}

func GetUnitPrivateAddress() (string, error) {
	return getUnit("private-address")
}
