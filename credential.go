package goops

import "encoding/json"

const (
	credentialGetCommand = "credential-get" // #nosec G101
)

func GetCredential() (map[string]string, error) {
	commandRunner := GetCommandRunner()

	args := []string{"--format=json"}

	output, err := commandRunner.Run(credentialGetCommand, args...)
	if err != nil {
		return nil, err
	}

	var credential map[string]string

	err = json.Unmarshal(output, &credential)
	if err != nil {
		return nil, err
	}

	return credential, nil
}
