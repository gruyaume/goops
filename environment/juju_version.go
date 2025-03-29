package environment

const JujuVersionEnvVar = "JUJU_VERSION"

func JujuVersion(getter EnvironmentGetter) string {
	return getter.Get(JujuVersionEnvVar)
}
