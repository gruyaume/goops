package commands

import "fmt"

const (
	resourceGetCommand = "resource-get"
)

func (command Command) ResourceGet(resourceName string) (string, error) {
	if resourceName == "" {
		return "", fmt.Errorf("resource name cannot be empty")
	}

	args := []string{resourceName}

	output, err := command.Runner.Run(resourceGetCommand, args...)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
