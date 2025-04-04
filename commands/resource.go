package commands

import "fmt"

const (
	resourceGetCommand = "resource-get"
)

type ResourceGetOptions struct {
	Name string
}

func (command Command) ResourceGet(opts *ResourceGetOptions) (string, error) {
	if opts.Name == "" {
		return "", fmt.Errorf("resource name cannot be empty")
	}

	args := []string{opts.Name}

	output, err := command.Runner.Run(resourceGetCommand, args...)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
