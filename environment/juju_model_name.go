package environment

const JujuModelNameEnvVar = "JUJU_MODEL_NAME"

func (env Environment) JujuModelName() string {
	return env.Getter.Get(JujuModelNameEnvVar)
}
