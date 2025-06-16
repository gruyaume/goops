package goops

const (
	resourceGetCommand = "resource-get"
)

func GetResource(name string) (string, error) {
	commandRunner := GetCommandRunner()

	args := []string{name}

	output, err := commandRunner.Run(resourceGetCommand, args...)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
