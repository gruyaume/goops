package environment

const JujuCharmHTTPProxyEnvVar = "JUJU_CHARM_HTTP_PROXY"

func (env Environment) JujuCharmHTTPProxy() string {
	return env.Getter.Get(JujuCharmHTTPProxyEnvVar)
}
