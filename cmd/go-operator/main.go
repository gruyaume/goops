package main

import (
	"log"
	"os"

	"github.com/gruyaume/go-operator/internal/commands"
	"github.com/gruyaume/go-operator/internal/events"
)

func main() {
	log.Println("Starting go-operator")
	eventType, err := events.GetEventType()
	if err != nil {
		log.Println("could not get event type:", err)
		os.Exit(0)
	}
	log.Println("Event type:", eventType)
	commandRunner := &commands.DefaultRunner{}
	err = commands.StatusSet(commandRunner, commands.StatusActive)
	if err != nil {
		log.Println("could not set status:", err)
		os.Exit(1)
	}
	log.Println("Status set to active")

	mySecret, err := commands.SecretGet(commandRunner, "", "my-label", false, true)
	if err != nil {
		log.Println("could not get secret:", err)
		myNewSecretContent := map[string]string{
			"username": "admin",
			"password": "password",
		}
		_, err := commands.SecretAdd(commandRunner, myNewSecretContent, "my secret", "my-label")
		if err != nil {
			log.Println("could not add secret:", err)
			os.Exit(0)
		}
		log.Println("Created new secret")
	} else {
		log.Println("Secret content:", mySecret)
	}

	log.Println("go-operator finished")
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
