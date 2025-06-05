package goops

import "github.com/gruyaume/goops/environment"

type Environment struct {
	ActionName         string
	AgentSocketAddress string
	AgentSocketNetwork string
	APIAddresses       string
	AvailabilityZone   string
	CharmDir           string
	CharmFTPProxy      string
	CharmHTTPProxy     string
	CharmHTTPSProxy    string
	CharmNoProxy       string
	ContextID          string
	HookName           string
	MachineID          string
	ModelName          string
	ModelUUID          string
	PrincipalUnit      string
	UnitName           string
	Version            string
}

func ReadEnv() Environment {
	env := environment.NewEnvironment()

	return Environment{
		ActionName:         env.JujuActionName(),
		AgentSocketAddress: env.JujuAgentSocketAddress(),
		AgentSocketNetwork: env.JujuAgentSocketNetwork(),
		APIAddresses:       env.JujuAPIAddresses(),
		AvailabilityZone:   env.JujuAvailabilityZone(),
		CharmDir:           env.JujuCharmDir(),
		CharmFTPProxy:      env.JujuCharmFTPProxy(),
		CharmHTTPProxy:     env.JujuCharmHTTPProxy(),
		CharmHTTPSProxy:    env.JujuCharmHTTPSProxy(),
		CharmNoProxy:       env.JujuCharmNoProxy(),
		ContextID:          env.JujuContextID(),
		HookName:           env.JujuHookName(),
		MachineID:          env.JujuMachineID(),
		ModelName:          env.JujuModelName(),
		ModelUUID:          env.JujuModelUUID(),
		PrincipalUnit:      env.JujuPrincipalUnit(),
		UnitName:           env.JujuUnitName(),
		Version:            env.JujuVersion(),
	}
}
