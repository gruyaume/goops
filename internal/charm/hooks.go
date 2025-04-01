package charm

import (
	"fmt"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/goops/internal/integrations/tls_certificates"
)

const (
	CaCertificateSecretLabel   = "active-ca-certificates" // #nosec G101
	TLSCertificatesIntegration = "certificates"
)

func isConfigValid(hookContext *goops.HookContext) (bool, error) {
	caCommonNameConfig, err := hookContext.Commands.ConfigGet("ca-common-name")
	if err != nil {
		return false, fmt.Errorf("could not get config: %w", err)
	}

	if caCommonNameConfig == "" {
		return false, fmt.Errorf("ca-common-name config is empty")
	}

	return true, nil
}

func generateAndStoreRootCertificate(hookContext *goops.HookContext) error {
	caCommonName, err := hookContext.Commands.ConfigGetString("ca-common-name")
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	_, err = hookContext.Commands.SecretGet("", CaCertificateSecretLabel, false, true)
	if err != nil {
		hookContext.Commands.JujuLog(commands.Info, "could not get secret:", err.Error())

		caCertPEM, caKeyPEM, err := GenerateRootCertificate(caCommonName)
		if err != nil {
			return fmt.Errorf("could not generate root certificate: %w", err)
		}

		hookContext.Commands.JujuLog(commands.Info, "Generated new root certificate")

		secretContent := map[string]string{
			"private-key":    caKeyPEM,
			"ca-certificate": caCertPEM,
		}

		_, err = hookContext.Commands.SecretAdd(secretContent, "", CaCertificateSecretLabel)
		if err != nil {
			return fmt.Errorf("could not add secret: %w", err)
		}

		hookContext.Commands.JujuLog(commands.Info, "Created new secret")

		return nil
	}

	hookContext.Commands.JujuLog(commands.Info, "Secret found")

	return nil
}

func processOutstandingCertificateRequests(hookContext *goops.HookContext) error {
	outstandingCertificateRequests, err := tls_certificates.GetOutstandingCertificateRequests(hookContext, TLSCertificatesIntegration)
	if err != nil {
		return fmt.Errorf("could not get outstanding certificate requests: %w", err)
	}

	for _, request := range outstandingCertificateRequests {
		hookContext.Commands.JujuLog(commands.Info, "Received a certificate signing request from:", request.RelationID, "with common name:", request.CertificateSigningRequest.CommonName)

		caCertificateSecret, err := hookContext.Commands.SecretGet("", CaCertificateSecretLabel, false, true)
		if err != nil {
			return fmt.Errorf("could not get CA certificate secret: %w", err)
		}

		caKeyPEM, ok := caCertificateSecret["private-key"]
		if !ok {
			return fmt.Errorf("could not find CA private key in secret")
		}

		caCertPEM, ok := caCertificateSecret["ca-certificate"]
		if !ok {
			return fmt.Errorf("could not find CA certificate in secret")
		}

		certPEM, err := GenerateCertificate(caKeyPEM, caCertPEM, request.CertificateSigningRequest.Raw)
		if err != nil {
			return fmt.Errorf("could not generate certificate: %w", err)
		}

		providerCertificatte := tls_certificates.ProviderCertificate{
			RelationID:                request.RelationID,
			Certificate:               tls_certificates.Certificate{Raw: certPEM},
			CertificateSigningRequest: request.CertificateSigningRequest,
			CA:                        tls_certificates.Certificate{Raw: caCertPEM},
			Chain: []tls_certificates.Certificate{
				{Raw: caCertPEM},
			},
			Revoked: false,
		}

		err = tls_certificates.SetRelationCertificate(hookContext, request.RelationID, providerCertificatte)
		if err != nil {
			hookContext.Commands.JujuLog(commands.Warning, "Could not set relation certificate:", err.Error())
			continue
		}

		hookContext.Commands.JujuLog(commands.Info, "Provided certificate to:", request.RelationID)
	}

	return nil
}

func setPorts(hookContext *goops.HookContext) error {
	ports := []commands.Port{
		{
			Port:     443,
			Protocol: "tcp",
		},
	}

	err := hookContext.Commands.SetPorts(ports)
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	return nil
}

func HandleDefaultHook(hookContext *goops.HookContext) error {
	isLeader, err := hookContext.Commands.IsLeader()
	if err != nil {
		return fmt.Errorf("could not check if unit is leader: %w", err)
	}

	if !isLeader {
		return fmt.Errorf("unit is not leader")
	}

	err = setPorts(hookContext)
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	hookContext.Commands.JujuLog(commands.Info, "Set unit ports")

	valid, err := isConfigValid(hookContext)
	if err != nil {
		return fmt.Errorf("could not check config: %w", err)
	}

	if !valid {
		return fmt.Errorf("config is not valid")
	}

	err = generateAndStoreRootCertificate(hookContext)
	if err != nil {
		return fmt.Errorf("could not generate CA certificate: %w", err)
	}

	err = processOutstandingCertificateRequests(hookContext)
	if err != nil {
		return fmt.Errorf("could not process outstanding certificate requests: %w", err)
	}

	err = hookContext.Commands.StatusSet(commands.StatusActive, "")
	if err != nil {
		hookContext.Commands.JujuLog(commands.Error, "Could not set status:", err.Error())
		return fmt.Errorf("could not set status: %w", err)
	}

	hookContext.Commands.JujuLog(commands.Info, "Status set to active")

	return nil
}
