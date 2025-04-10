package environment

const JujuAgentSocketNetworkEnvVar = "JUJU_AGENT_SOCKET_NETWORK"

func (env Environment) JujuAgentSocketNetwork() string {
	return env.Getter.Get(JujuAgentSocketNetworkEnvVar)
}
