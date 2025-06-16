package goops

import "fmt"

const (
	applicationVersionSetCommand = "application-version-set"
)

func SetAppVersion(version string) error {
	commandRunner := GetCommandRunner()

	args := []string{}
	if version != "" {
		args = append(args, version)
	}

	_, err := commandRunner.Run(applicationVersionSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set application version: %w", err)
	}

	return nil
}
