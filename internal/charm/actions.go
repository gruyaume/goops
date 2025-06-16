package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
)

type GetCACertificateActionParams struct {
	AcceptTOS bool `json:"accept-tos"`
}

func HandleGetCACertificateAction() error {
	getCAActionParams := GetCACertificateActionParams{}

	err := goops.GetActionParams(&getCAActionParams)
	if err != nil {
		return fmt.Errorf("could not get action parameter 'accept-tos': %w", err)
	}

	if !getCAActionParams.AcceptTOS {
		err := goops.FailActionf("You must accept the terms of service to get the CA certificate")
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not get action parameter 'accept-tos': %w", err, err)
		}

		return fmt.Errorf("you must accept the terms of service to get the CA certificate")
	}

	caCertificateSecret, err := goops.GetSecretByLabel(CaCertificateSecretLabel, false, true)
	if err != nil {
		err := goops.FailActionf("could not get CA certificate secret")
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not get CA certificate secret: %w", err, err)
		}

		return fmt.Errorf("could not get CA certificate secret: %w", err)
	}

	caCertPEM, ok := caCertificateSecret["ca-certificate"]
	if !ok {
		err := goops.FailActionf("could not find CA certificate in secret")
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not find CA certificate in secret: %w", err, err)
		}

		return fmt.Errorf("could not find CA certificate in secret")
	}

	err = goops.SetActionResults(map[string]string{
		"ca-certificate": caCertPEM,
	})
	if err != nil {
		err := goops.FailActionf("could not set action result")
		if err != nil {
			return fmt.Errorf("could not fail action: %w and could not set action result: %w", err, err)
		}

		return fmt.Errorf("could not set action result: %w", err)
	}

	return nil
}
