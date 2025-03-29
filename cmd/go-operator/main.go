package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/gruyaume/go-operator/internal/commands"
	"github.com/gruyaume/go-operator/internal/events"
)

const (
	CaCertificateSecretLabel   = "active-ca-certificates"
	CaCertificateValidityYears = 10
	TLSCertificatesIntegration = "certificates"
)

func generateRootCertificate1(commonName string) (string, string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("could not generate private key: %w", err)
	}
	csrTemplate := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: commonName,
		},
		DNSNames: []string{},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, priv)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate request: %w", err)
	}
	csrPEM := new(bytes.Buffer)
	err = pem.Encode(csrPEM, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to encode certificate request: %w", err)
	}
	certTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(CaCertificateValidityYears, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, &priv.PublicKey, priv)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate: %w", err)
	}
	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to encode certificate: %w", err)
	}
	caCert := certPEM.String()
	caKey := new(bytes.Buffer)
	err = pem.Encode(caKey, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to encode private key: %w", err)
	}
	caKeyPEM := caKey.String()
	return caCert, caKeyPEM, nil
}

func generateRootCertificate(commandRunner *commands.DefaultRunner, logger *commands.Logger) error {
	caCommonName, err := commands.ConfigGet(commandRunner, "ca-common-name")
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	_, err = commands.SecretGet(commandRunner, "", CaCertificateSecretLabel, false, true)
	if err != nil {
		logger.Info("could not get secret:", err.Error())
		caCert, caKeyPEM, err := generateRootCertificate1(caCommonName)
		if err != nil {
			return fmt.Errorf("could not generate root certificate: %w", err)
		}
		myNewSecretContent := map[string]string{
			"private-key":    caKeyPEM,
			"ca-certificate": caCert,
		}
		_, err = commands.SecretAdd(commandRunner, myNewSecretContent, "", CaCertificateSecretLabel)
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

	eventType, err := events.GetEventType()
	if err != nil {
		logger.Info("Could not get event type: ", err.Error())
		os.Exit(0)
	}
	logger.Info("Event type:", string(eventType))

	err = commands.StatusSet(commandRunner, commands.StatusActive)
	if err != nil {
		logger.Error("Could not set status:", err.Error())
		os.Exit(0)
	}
	logger.Info("Status set to active")

	err = generateRootCertificate(commandRunner, logger)
	if err != nil {
		logger.Error("Could not generate CA certificate:", err.Error())
		os.Exit(0)
	}
	relationIDs, err := commands.RelationIDs(commandRunner, TLSCertificatesIntegration)
	if err != nil {
		logger.Error("Could not get relation IDs:", err.Error())
		os.Exit(0)
	}
	for relationID := range relationIDs {
		logger.Info("Relation ID:", fmt.Sprintf("%d", relationID))
	}
	logger.Info("Finished go-operator")
	os.Exit(0)
}
