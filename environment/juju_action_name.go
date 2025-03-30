package environment

const JujuActionNameEnvVar = "JUJU_ACTION_NAME"

func (env Environment) JujuActionName() string {
	return env.Getter.Get(JujuActionNameEnvVar)
}
