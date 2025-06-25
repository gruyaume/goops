package goops

import "fmt"

const (
	appVersionSetCommand = "application-version-set"
)

// SetAppVersion sets the application version.
// The version set will be displayed in “juju status” output for the application.
func SetAppVersion(version string) error {
	commandRunner := GetCommandRunner()

	args := []string{}
	if version != "" {
		args = append(args, version)
	}

	_, err := commandRunner.Run(appVersionSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set application version: %w", err)
	}

	return nil
}
