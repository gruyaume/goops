package environment

const JujuCharmFTPProxyEnvVar = "JUJU_CHARM_FTP_PROXY"

func (env Environment) JujuCharmFTPProxy() string {
	return env.Getter.Get(JujuCharmFTPProxyEnvVar)
}
