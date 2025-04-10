package environment

const JujuCharmNoProxyEnvVar = "JUJU_CHARM_NO_PROXY"

func (env Environment) JujuCharmNoProxy() string {
	return env.Getter.Get(JujuCharmNoProxyEnvVar)
}
