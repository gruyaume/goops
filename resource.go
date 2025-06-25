package goops

const (
	resourceGetCommand = "resource-get"
)

// GetResource retrieves the local path to a resource file for the given resource name.
func GetResource(name string) (string, error) {
	commandRunner := GetCommandRunner()

	args := []string{name}

	output, err := commandRunner.Run(resourceGetCommand, args...)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
