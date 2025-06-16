package charm

import (
	"fmt"
	"time"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/internal/integrations/tls_certificates"
)

const (
	CaCertificateSecretLabel   = "active-ca-certificates" // #nosec G101
	TLSCertificatesIntegration = "certificates"
)

func isConfigValid() (bool, error) {
	caCommonNameConfig, err := goops.GetConfig("ca-common-name")
	if err != nil {
		return false, fmt.Errorf("could not get config: %w", err)
	}

	if caCommonNameConfig == "" {
		return false, fmt.Errorf("ca-common-name config is empty")
	}

	return true, nil
}

func generateAndStoreRootCertificate() error {
	caCommonName, err := goops.GetConfigString("ca-common-name")
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	_, err = goops.GetSecretByLabel(CaCertificateSecretLabel, false, true)
	if err != nil {
		goops.LogInfof("could not get secret: %s", err.Error())

		caCertPEM, caKeyPEM, err := GenerateRootCertificate(caCommonName)
		if err != nil {
			return fmt.Errorf("could not generate root certificate: %w", err)
		}

		goops.LogInfof("Generated new root certificate with common name: %s", caCommonName)

		secretContent := map[string]string{
			"private-key":    caKeyPEM,
			"ca-certificate": caCertPEM,
		}

		expiry := time.Now().AddDate(1, 0, 0)

		output, err := goops.AddSecret(&goops.AddSecretOptions{
			Content:     secretContent,
			Description: "ca certificate and private key for the certificates charm",
			Expire:      expiry,
			Label:       CaCertificateSecretLabel,
			Rotate:      goops.RotateNever,
		})
		if err != nil {
			return fmt.Errorf("could not add secret: %w", err)
		}

		goops.LogInfof("Created new secret with ID: %s", output)

		return nil
	}

	secretInfo, err := goops.GetSecretInfoByLabel(CaCertificateSecretLabel)
	if err != nil {
		return fmt.Errorf("could not get secret info: %w", err)
	}

	if secretInfo == nil {
		return fmt.Errorf("secret info is nil")
	}

	return nil
}

func processOutstandingCertificateRequests() error {
	outstandingCertificateRequests, err := tls_certificates.GetOutstandingCertificateRequests(TLSCertificatesIntegration)
	if err != nil {
		return fmt.Errorf("could not get outstanding certificate requests: %w", err)
	}

	for _, request := range outstandingCertificateRequests {
		goops.LogInfof("Received a certificate signing request from: %s with common name: %s", request.RelationID, request.CertificateSigningRequest.CommonName)

		caCertificateSecret, err := goops.GetSecretByLabel(CaCertificateSecretLabel, false, true)
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

		err = tls_certificates.SetRelationCertificate(request.RelationID, providerCertificatte)
		if err != nil {
			goops.LogWarningf("Could not set relation certificate: %s", err.Error())
			continue
		}

		goops.LogInfof("Provided certificate to: %s", request.RelationID)
	}

	return nil
}

func setPorts() error {
	err := goops.SetPorts([]*goops.Port{
		{
			Port:     443,
			Protocol: "tcp",
		},
	})
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	return nil
}

func validateNetworkGet() error {
	bindAddress, err := goops.GetNetworkBindAddress("certificates")
	if err != nil {
		return fmt.Errorf("could not get network config: %w", err)
	}

	if bindAddress == "" {
		return fmt.Errorf("network config bind addresses is empty")
	}

	goops.LogInfof("Bind address: %s", bindAddress)

	ingressAddress, err := goops.GetNetworkIngressAddress("certificates")
	if err != nil {
		return fmt.Errorf("could not get network ingress addresses: %w", err)
	}

	if ingressAddress == "" {
		return fmt.Errorf("network config ingress address is empty")
	}

	goops.LogInfof("Ingress address: %s", ingressAddress)

	egressSubnets, err := goops.GetNetworkEgressSubnets("certificates")
	if err != nil {
		return fmt.Errorf("could not get network egress subnets: %w", err)
	}

	if len(egressSubnets) == 0 {
		return fmt.Errorf("network config egress subnets is empty")
	}

	if egressSubnets[0] == "" {
		return fmt.Errorf("network config egress subnet is empty")
	}

	return nil
}

func validateState() error {
	stateKey := "my-key"
	stateValue := "my-value"

	_, err := goops.GetState(stateKey)
	if err != nil {
		goops.LogInfof("could not get state: %s", err.Error())

		err := goops.SetState(stateKey, stateValue)
		if err != nil {
			return fmt.Errorf("could not set state: %w", err)
		}

		goops.LogInfof("set state: %s = %s", stateKey, stateValue)

		return nil
	} else {
		goops.LogInfof("state already set: %s = %s", stateKey, stateValue)

		err := goops.DeleteState(stateKey)
		if err != nil {
			return fmt.Errorf("could not delete state: %w", err)
		}

		goops.LogInfof("deleted state: %s", stateKey)
	}

	return nil
}

func Configure() error {
	meta, err := goops.ReadMetadata()
	if err != nil {
		return fmt.Errorf("could not read metadata: %w", err)
	}

	isLeader, err := goops.IsLeader()
	if err != nil {
		return fmt.Errorf("could not check if unit is leader: %w", err)
	}

	if !isLeader {
		return fmt.Errorf("unit is not leader")
	}

	goops.LogInfof("Charm Name: %s", meta.Name)

	err = setPorts()
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	privateAddress, err := goops.GetUnitPrivateAddress()
	if err != nil {
		return fmt.Errorf("could not get unit private address: %w", err)
	}

	if privateAddress == "" {
		return fmt.Errorf("unit private address is empty")
	}

	goops.LogInfof("Unit private address: %s", privateAddress)

	goops.LogInfof("Set unit ports")

	valid, err := isConfigValid()
	if err != nil {
		return fmt.Errorf("could not check config: %w", err)
	}

	if !valid {
		return fmt.Errorf("config is not valid")
	}

	err = generateAndStoreRootCertificate()
	if err != nil {
		return fmt.Errorf("could not generate CA certificate: %w", err)
	}

	err = processOutstandingCertificateRequests()
	if err != nil {
		return fmt.Errorf("could not process outstanding certificate requests: %w", err)
	}

	goalState, err := goops.GetGoalState()
	if err != nil {
		return fmt.Errorf("could not get goal state: %w", err)
	}

	if goalState == nil {
		return fmt.Errorf("goal state is nil")
	}

	if goalState.Units == nil {
		return fmt.Errorf("goal state units is nil")
	}

	if goalState.Units["example/0"] == nil {
		return fmt.Errorf("goal state unit is nil")
	}

	_, err = goops.GetCredential()
	if err == nil {
		return fmt.Errorf("expected not to get container on caas model: %w", err)
	}

	err = validateNetworkGet()
	if err != nil {
		return fmt.Errorf("could not validate network get: %w", err)
	}

	certificatesRelationID, err := goops.GetRelationIDs(TLSCertificatesIntegration)
	if err != nil {
		return fmt.Errorf("could not get relation ID: %w", err)
	}

	if len(certificatesRelationID) > 0 {
		uuid, err := goops.GetRelationModel(certificatesRelationID[0])
		if err != nil {
			return fmt.Errorf("could not get relation model: %w", err)
		}

		if uuid == "" {
			return fmt.Errorf("relation model UUID is empty")
		}
	}

	err = validateState()
	if err != nil {
		return fmt.Errorf("could not validate state: %w", err)
	}

	err = goops.SetApplicationVersion("1.0.0")
	if err != nil {
		return fmt.Errorf("could not set application version using goops: %w", err)
	}

	existingStatus, err := goops.GetStatus()
	if err != nil {
		return fmt.Errorf("could not get status: %w", err)
	}

	goops.LogInfof("Current status: %s %s", existingStatus.Code, existingStatus.Message)

	err = goops.SetUnitStatus(goops.StatusActive, "A happy charm")
	if err != nil {
		return fmt.Errorf("could not set unit status: %w", err)
	}

	goops.LogInfof("Status set to active")

	return nil
}
