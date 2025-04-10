package environment

const JujuCharmHTTPSProxyEnvVar = "JUJU_CHARM_HTTPS_PROXY"

func (env Environment) JujuCharmHTTPSProxy() string {
	return env.Getter.Get(JujuCharmHTTPSProxyEnvVar)
}
