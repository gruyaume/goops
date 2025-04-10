package environment

const JujuAgentSocketAddressEnvVar = "JUJU_AGENT_SOCKET_ADDRESS"

func (env Environment) JujuAgentSocketAddress() string {
	return env.Getter.Get(JujuAgentSocketAddressEnvVar)
}
