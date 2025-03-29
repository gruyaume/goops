package environment

const JujuActionNameEnvVar = "JUJU_ACTION_NAME"

func JujuActionName(getter EnvironmentGetter) string {
	return getter.Get(JujuActionNameEnvVar)
}
