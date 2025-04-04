package commands

import (
	"encoding/json"
	"fmt"
)

const (
	unitGetCommand = "unit-get"
)

type UnitGetOptions struct {
	PrivateAddress bool
	PublicAddress  bool
}

func (command Command) UnitGet(opts *UnitGetOptions) (string, error) {
	if !opts.PrivateAddress && !opts.PublicAddress {
		return "", fmt.Errorf("must specify either private or public address")
	}

	if opts.PrivateAddress && opts.PublicAddress {
		return "", fmt.Errorf("cannot specify both private and public address")
	}

	args := []string{}

	if opts.PrivateAddress {
		args = append(args, "private-address")
	}

	if opts.PublicAddress {
		args = append(args, "public-address")
	}

	args = append(args, "--format=json")

	output, err := command.Runner.Run(unitGetCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get unit: %w", err)
	}

	var result string
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("failed to parse unit get output: %w", err)
	}

	return result, nil
}
