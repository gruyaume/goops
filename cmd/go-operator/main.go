package main

import (
	"fmt"
	"os"

	"github.com/gruyaume/go-operator/internal/commands"
	"github.com/gruyaume/go-operator/internal/events"
)

const (
	TLSCertificatesIntegration = "certificates"
)

func GenerateCACertificate(commandRunner *commands.DefaultRunner, logger *commands.Logger) error {
	_, err := commands.SecretGet(commandRunner, "", "my-label", false, true)
	if err != nil {
		logger.Info("could not get secret:", err.Error())
		myNewSecretContent := map[string]string{
			"username": "admin",
			"password": "password",
		}
		_, err := commands.SecretAdd(commandRunner, myNewSecretContent, "my secret", "my-label")
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

	err = GenerateCACertificate(commandRunner, logger)
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

// type CreateCertificateAuthorityParams struct {
// 	CommonName          string `json:"common_name"`
// 	SANsDNS             string `json:"sans_dns"`
// 	CountryName         string `json:"country_name"`
// 	StateOrProvinceName string `json:"state_or_province_name"`
// 	LocalityName        string `json:"locality_name"`
// 	OrganizationName    string `json:"organization_name"`
// 	OrganizationalUnit  string `json:"organizational_unit_name"`
// 	NotValidAfter       string `json:"not_valid_after"`
// }

// func createCertificateAuthority(fields CreateCertificateAuthorityParams) (string, string, string, error) {
// 	// Create the private key for the CA
// 	priv, err := rsa.GenerateKey(rand.Reader, 4096)
// 	if err != nil {
// 		return "", "", ""
// 	}

// 	privPEM := new(bytes.Buffer)
// 	err = pem.Encode(privPEM, &pem.Block{
// 		Type:  "RSA PRIVATE KEY",
// 		Bytes: x509.MarshalPKCS1PrivateKey(priv),
// 	})
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to encode private key: %w", err)
// 	}

// 	// Create the certificate request for the CA
// 	csrTemplate := &x509.CertificateRequest{
// 		Subject: pkix.Name{
// 			CommonName:         fields.CommonName,
// 			Country:            []string{fields.CountryName},
// 			Province:           []string{fields.StateOrProvinceName},
// 			Locality:           []string{fields.LocalityName},
// 			Organization:       []string{fields.OrganizationName},
// 			OrganizationalUnit: []string{fields.OrganizationalUnit},
// 		},
// 		DNSNames: []string{fields.SANsDNS},
// 	}

// 	if fields.SANsDNS != "" {
// 		csrTemplate.DNSNames = []string{fields.SANsDNS}
// 	}

// 	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, priv)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to create certificate request: %w", err)
// 	}

// 	csrPEM := new(bytes.Buffer)
// 	err = pem.Encode(csrPEM, &pem.Block{
// 		Type:  "CERTIFICATE REQUEST",
// 		Bytes: csrBytes,
// 	})
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to encode certificate request: %w", err)
// 	}

// 	template := &x509.Certificate{
// 		SerialNumber: big.NewInt(time.Now().UnixNano()),
// 		Subject: pkix.Name{
// 			CommonName:         fields.CommonName,
// 			Country:            []string{fields.CountryName},
// 			Province:           []string{fields.StateOrProvinceName},
// 			Locality:           []string{fields.LocalityName},
// 			Organization:       []string{fields.OrganizationName},
// 			OrganizationalUnit: []string{fields.OrganizationalUnit},
// 		},
// 		NotBefore:             time.Now(),
// 		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
// 		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
// 		BasicConstraintsValid: true,
// 		IsCA:                  true,
// 	}

// 	if fields.NotValidAfter != "" {
// 		notAfter, err := time.Parse(time.RFC3339, fields.NotValidAfter)
// 		if err != nil {
// 			return "", "", "", fmt.Errorf("failed to parse NotValidAfter: %w", err)
// 		}
// 		template.NotAfter = notAfter
// 	} else {
// 		template.NotAfter = time.Now().AddDate(10, 0, 0) // Default 10 years
// 	}

// 	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to create certificate: %w", err)
// 	}

// 	certPEM := new(bytes.Buffer)
// 	err = pem.Encode(certPEM, &pem.Block{
// 		Type:  "CERTIFICATE",
// 		Bytes: derBytes,
// 	})
// 	if err != nil {
// 		return "", "", "", fmt.Errorf("failed to encode certificate: %w", err)
// 	}

// 	return csrPEM.String(), privPEM.String(), certPEM.String(), nil
// }
