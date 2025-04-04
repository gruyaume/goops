package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/commands"
)

func HandleGetCACertificateAction(hookContext *goops.HookContext) error {
	secretGetOptions := &commands.SecretGetOptions{
		Label:   CaCertificateSecretLabel,
		Refresh: true,
	}

	caCertificateSecret, err := hookContext.Commands.SecretGet(secretGetOptions)
	if err != nil {
		actionFailOptions := &commands.ActionFailOptions{
			Message: "could not get CA certificate secret",
		}

		err := hookContext.Commands.ActionFail(actionFailOptions)
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not get CA certificate secret: %w", err, err)
		}

		return fmt.Errorf("could not get CA certificate secret: %w", err)
	}

	caCertPEM, ok := caCertificateSecret["ca-certificate"]
	if !ok {
		actionFailOptions := &commands.ActionFailOptions{
			Message: "could not find CA certificate in secret",
		}

		err := hookContext.Commands.ActionFail(actionFailOptions)
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not find CA certificate in secret: %w", err, err)
		}

		return fmt.Errorf("could not find CA certificate in secret")
	}

	actionSetOptions := &commands.ActionSetOptions{
		Content: map[string]string{
			"ca-certificate": caCertPEM,
		},
	}

	err = hookContext.Commands.ActionSet(actionSetOptions)
	if err != nil {
		actionFailOptions := &commands.ActionFailOptions{
			Message: "could not set action result",
		}

		err := hookContext.Commands.ActionFail(actionFailOptions)
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not set action result: %w", err, err)
		}

		return fmt.Errorf("could not set action result: %w", err)
	}

	return nil
}
