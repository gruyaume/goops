package charm

import (
	"fmt"
	"time"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/commands"
	"github.com/gruyaume/goops/internal/integrations/tls_certificates"
	"github.com/gruyaume/goops/metadata"
)

const (
	CaCertificateSecretLabel   = "active-ca-certificates" // #nosec G101
	TLSCertificatesIntegration = "certificates"
)

func isConfigValid(hookContext *goops.HookContext) (bool, error) {
	configGetOptions := &commands.ConfigGetOptions{
		Key: "ca-common-name",
	}

	caCommonNameConfig, err := hookContext.Commands.ConfigGet(configGetOptions)
	if err != nil {
		return false, fmt.Errorf("could not get config: %w", err)
	}

	if caCommonNameConfig == "" {
		return false, fmt.Errorf("ca-common-name config is empty")
	}

	return true, nil
}

func generateAndStoreRootCertificate(hookContext *goops.HookContext) error {
	configGetOptions := &commands.ConfigGetOptions{
		Key: "ca-common-name",
	}

	caCommonName, err := hookContext.Commands.ConfigGetString(configGetOptions)
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	secretGetOptions := &commands.SecretGetOptions{
		Label:   CaCertificateSecretLabel,
		Refresh: true,
	}

	_, err = hookContext.Commands.SecretGet(secretGetOptions)
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

		expiry := time.Now().AddDate(1, 0, 0)

		secretAddOptions := &commands.SecretAddOptions{
			Content:     secretContent,
			Description: "ca certificate and private key for the certificates charm",
			Expire:      expiry,
			Label:       CaCertificateSecretLabel,
		}

		output, err := hookContext.Commands.SecretAdd(secretAddOptions)
		if err != nil {
			return fmt.Errorf("could not add secret: %w", err)
		}

		hookContext.Commands.JujuLog(commands.Info, "Created new secret")
		hookContext.Commands.JujuLog(commands.Info, "Secret ID:", output)

		return nil
	}

	secretInfoGetOpts := &commands.SecretInfoGetOptions{
		Label: CaCertificateSecretLabel,
	}

	secretInfo, err := hookContext.Commands.SecretInfoGet(secretInfoGetOpts)
	if err != nil {
		return fmt.Errorf("could not get secret info: %w", err)
	}

	if secretInfo == nil {
		return fmt.Errorf("secret info is nil")
	}

	return nil
}

func processOutstandingCertificateRequests(hookContext *goops.HookContext) error {
	outstandingCertificateRequests, err := tls_certificates.GetOutstandingCertificateRequests(hookContext, TLSCertificatesIntegration)
	if err != nil {
		return fmt.Errorf("could not get outstanding certificate requests: %w", err)
	}

	for _, request := range outstandingCertificateRequests {
		hookContext.Commands.JujuLog(commands.Info, "Received a certificate signing request from:", request.RelationID, "with common name:", request.CertificateSigningRequest.CommonName)

		secretGetOptions := &commands.SecretGetOptions{
			Label:   CaCertificateSecretLabel,
			Refresh: true,
		}

		caCertificateSecret, err := hookContext.Commands.SecretGet(secretGetOptions)
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
	ports := &commands.SetPortsOptions{
		Ports: []*commands.Port{
			{
				Port:     443,
				Protocol: "tcp",
			},
		},
	}

	err := hookContext.Commands.SetPorts(ports)
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	return nil
}

func validateNetworkGet(hookContext *goops.HookContext) error {
	networkGetOpts := &commands.NetworkGetOptions{
		BindingName: "certificates",
	}

	networkConfig, err := hookContext.Commands.NetworkGet(networkGetOpts)
	if err != nil {
		return fmt.Errorf("could not get network config: %w", err)
	}

	if networkConfig == nil {
		return fmt.Errorf("network config is nil")
	}

	if len(networkConfig.BindAddresses) == 0 {
		return fmt.Errorf("network config bind addresses is empty")
	}

	if len(networkConfig.BindAddresses[0].Addresses) == 0 {
		return fmt.Errorf("network config bind address addresses is empty")
	}

	if networkConfig.BindAddresses[0].Addresses[0].Value == "" {
		return fmt.Errorf("network config bind address address value is empty- This can happen in the first stage of the deployment")
	}

	if len(networkConfig.IngressAddresses) == 0 {
		return fmt.Errorf("network config ingress addresses is empty")
	}

	if networkConfig.IngressAddresses[0] == "" {
		return fmt.Errorf("network config ingress address is empty")
	}

	if len(networkConfig.EgressSubnets) == 0 {
		return fmt.Errorf("network config egress subnets is empty")
	}

	if networkConfig.EgressSubnets[0] == "" {
		return fmt.Errorf("network config egress subnet is empty")
	}

	return nil
}

func validateState(hookContext *goops.HookContext) error {
	stateKey := "my-key"
	stateValue := "my-value"

	stateGetOptions := &commands.StateGetOptions{
		Key: stateKey,
	}

	_, err := hookContext.Commands.StateGet(stateGetOptions)
	if err != nil {
		hookContext.Commands.JujuLog(commands.Info, "could not get state:", err.Error())

		stateSetOptions := &commands.StateSetOptions{
			Key:   stateKey,
			Value: stateValue,
		}

		err := hookContext.Commands.StateSet(stateSetOptions)
		if err != nil {
			return fmt.Errorf("could not set state: %w", err)
		}

		hookContext.Commands.JujuLog(commands.Info, "set state:", stateKey, "=", stateValue)

		return nil
	} else {
		hookContext.Commands.JujuLog(commands.Info, "state already set:", stateKey, "=", stateValue)

		stateDeleteOptions := &commands.StateDeleteOptions{
			Key: stateKey,
		}

		err := hookContext.Commands.StateDelete(stateDeleteOptions)
		if err != nil {
			return fmt.Errorf("could not delete state: %w", err)
		}

		hookContext.Commands.JujuLog(commands.Info, "deleted state:", stateKey)
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

	meta, err := metadata.GetCharmMetadata(hookContext.Environment)
	if err != nil {
		return fmt.Errorf("could not get charm metadata: %w", err)
	}

	hookContext.Commands.JujuLog(commands.Info, "Charm Name:", meta.Name)

	err = setPorts(hookContext)
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	unitGetOpts := &commands.UnitGetOptions{
		PrivateAddress: true,
	}

	privateAddress, err := hookContext.Commands.UnitGet(unitGetOpts)
	if err != nil {
		return fmt.Errorf("could not get unit private address: %w", err)
	}

	if privateAddress == "" {
		return fmt.Errorf("unit private address is empty")
	}

	hookContext.Commands.JujuLog(commands.Info, "Unit private address:", privateAddress)

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

	goalState, err := hookContext.Commands.GoalState()
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

	_, err = hookContext.Commands.CredentialGet()
	if err == nil {
		return fmt.Errorf("expected not to get container on caas model: %w", err)
	}

	err = validateNetworkGet(hookContext)
	if err != nil {
		return fmt.Errorf("could not validate network get: %w", err)
	}

	relationIDsOptions := &commands.RelationIDsOptions{
		Name: TLSCertificatesIntegration,
	}

	certificatesRelationID, err := hookContext.Commands.RelationIDs(relationIDsOptions)
	if err != nil {
		return fmt.Errorf("could not get relation ID: %w", err)
	}

	if len(certificatesRelationID) > 0 {
		relationModelGetOptions := &commands.RelationModelGetOptions{
			ID: certificatesRelationID[0],
		}

		relationModel, err := hookContext.Commands.RelationModelGet(relationModelGetOptions)
		if err != nil {
			return fmt.Errorf("could not get relation model: %w", err)
		}

		if relationModel == nil {
			return fmt.Errorf("relation model is nil")
		}

		if relationModel.UUID == "" {
			return fmt.Errorf("relation model UUID is empty")
		}
	}

	err = validateState(hookContext)
	if err != nil {
		return fmt.Errorf("could not validate state: %w", err)
	}

	applicationVersionSetOptions := &commands.ApplicationVersionSetOptions{
		Version: "1.0.0",
	}

	err = hookContext.Commands.ApplicationVersionSet(applicationVersionSetOptions)
	if err != nil {
		return fmt.Errorf("could not set application version: %w", err)
	}

	existingStatus, err := hookContext.Commands.StatusGet()
	if err != nil {
		return fmt.Errorf("could not get status: %w", err)
	}

	hookContext.Commands.JujuLog(commands.Info, "Current status:", string(existingStatus.Name), existingStatus.Message)

	statusOpts := &commands.StatusSetOptions{
		Name:    commands.StatusActive,
		Message: "A happy charm",
	}

	err = hookContext.Commands.StatusSet(statusOpts)
	if err != nil {
		return fmt.Errorf("could not set status: %w", err)
	}

	hookContext.Commands.JujuLog(commands.Info, "Status set to active")

	return nil
}
