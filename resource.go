package goops

const (
	resourceGetCommand = "resource-get"
)

type ResourceGetOptions struct {
	Name string
}

func GetResource(name string) (string, error) {
	commandRunner := GetRunner()

	args := []string{name}

	output, err := commandRunner.Run(resourceGetCommand, args...)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
