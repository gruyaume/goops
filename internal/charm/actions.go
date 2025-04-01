package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
)

func HandleGetCACertificateAction(hookContext *goops.HookContext) error {
	caCertificateSecret, err := hookContext.Commands.SecretGet("", CaCertificateSecretLabel, false, true)
	if err != nil {
		err := hookContext.Commands.ActionFail("could not get CA certificate secret")
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not get CA certificate secret: %w", err, err)
		}

		return fmt.Errorf("could not get CA certificate secret: %w", err)
	}

	caCertPEM, ok := caCertificateSecret["ca-certificate"]
	if !ok {
		err := hookContext.Commands.ActionFail("could not find CA certificate in secret")
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not find CA certificate in secret: %w", err, err)
		}

		return fmt.Errorf("could not find CA certificate in secret")
	}

	err = hookContext.Commands.ActionSet(map[string]string{"ca-certificate": caCertPEM})
	if err != nil {
		err := hookContext.Commands.ActionFail("could not set action result")
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not set action result: %w", err, err)
		}

		return fmt.Errorf("could not set action result: %w", err)
	}

	return nil
}
