package main

import (
	"fmt"
	"os"

	"github.com/gruyaume/go-operator/internal/charm"
	"github.com/gruyaume/go-operator/internal/commands"
)

const (
	CaCertificateSecretLabel   = "active-ca-certificates" // #nosec G101
	TLSCertificatesIntegration = "certificates"
)

func generateAndStoreRootCertificate(commandRunner *commands.DefaultRunner, logger *commands.Logger) error {
	caCommonName, err := commands.ConfigGet(commandRunner, "ca-common-name")
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	_, err = commands.SecretGet(commandRunner, "", CaCertificateSecretLabel, false, true)
	if err != nil {
		logger.Info("could not get secret:", err.Error())
		caCert, caKeyPEM, err := charm.GenerateRootCertificate(caCommonName)
		if err != nil {
			return fmt.Errorf("could not generate root certificate: %w", err)
		}
		logger.Info("Generated new root certificate")
		secretContent := map[string]string{
			"private-key":    caKeyPEM,
			"ca-certificate": caCert,
		}
		_, err = commands.SecretAdd(commandRunner, secretContent, "", CaCertificateSecretLabel)
		if err != nil {
			return fmt.Errorf("could not add secret: %w", err)
		}
		logger.Info("Created new secret")
		return nil
	}
	logger.Info("Secret found")
	return nil
}

func isConfigValid(commandRunner *commands.DefaultRunner) (bool, error) {
	caCommonNameConfig, err := commands.ConfigGet(commandRunner, "ca-common-name")
	if err != nil {
		return false, fmt.Errorf("could not get config: %w", err)
	}
	if caCommonNameConfig == "" {
		return false, fmt.Errorf("ca-common-name config is empty")
	}
	return true, nil
}

func processOutstandingCertificateRequests(commandRunner *commands.DefaultRunner, logger *commands.Logger) error {
	outstandingCertificateRequests, err := charm.GetOutstandingCertificateRequests(commandRunner, TLSCertificatesIntegration)
	if err != nil {
		return fmt.Errorf("could not get outstanding certificate requests: %w", err)
	}
	for _, request := range outstandingCertificateRequests {
		logger.Info("Received a certificate signing request from:", request.RelationID, "with common name:", request.CertificateSigningRequest.CommonName)
		caCertificateSecret, err := commands.SecretGet(commandRunner, "", CaCertificateSecretLabel, false, true)
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
		certPEM, err := charm.GenerateCertificate(caKeyPEM, caCertPEM, request.CertificateSigningRequest.Raw)
		if err != nil {
			return fmt.Errorf("could not generate certificate: %w", err)
		}
		providerCertificatte := charm.ProviderCertificate{
			RelationID:                request.RelationID,
			Certificate:               charm.Certificate{Raw: certPEM},
			CertificateSigningRequest: request.CertificateSigningRequest,
			CA:                        charm.Certificate{Raw: caCertPEM},
			Chain: []charm.Certificate{
				{Raw: caCertPEM},
			},
			Revoked: false,
		}
		err = charm.SetRelationCertificate(commandRunner, request.RelationID, providerCertificatte)
		if err != nil {
			logger.Warning("Could not set relation certificate:", err.Error())
			continue
		}
		logger.Info("Provided certificate to:", request.RelationID)
	}
	return nil
}

func main() {
	commandRunner := &commands.DefaultRunner{}
	logger := commands.NewLogger(commandRunner)
	logger.Info("Started go-operator")

	isLeader, err := commands.IsLeader(commandRunner)
	if err != nil {
		logger.Info("Could not check if leader:", err.Error())
		os.Exit(0)
	}
	if !isLeader {
		logger.Info("not leader, exiting")
		os.Exit(0)
	}
	logger.Info("Unit is leader")

	valid, err := isConfigValid(commandRunner)
	if err != nil {
		logger.Info("Could not check config:", err.Error())
		os.Exit(0)
	}
	if !valid {
		logger.Info("Config is not valid, exiting")
		os.Exit(0)
	}
	logger.Info("Config is valid")

	err = commands.StatusSet(commandRunner, commands.StatusActive)
	if err != nil {
		logger.Error("Could not set status:", err.Error())
		os.Exit(0)
	}
	logger.Info("Status set to active")

	err = generateAndStoreRootCertificate(commandRunner, logger)
	if err != nil {
		logger.Error("Could not generate CA certificate:", err.Error())
		os.Exit(0)
	}
	err = processOutstandingCertificateRequests(commandRunner, logger)
	if err != nil {
		logger.Error("Could not process outstanding certificate requests:", err.Error())
		os.Exit(0)
	}
	logger.Info("Finished go-operator")
	os.Exit(0)
}
