package commands

import "encoding/json"

const (
	CredentialGetCommand = "credential-get"
)

func (command Command) CredentialGet() (map[string]string, error) {
	args := []string{"--format=json"}

	output, err := command.Runner.Run(CredentialGetCommand, args...)
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
