package environment

const JujuModelUUIDEnvVar = "JUJU_MODEL_UUID"

func (env Environment) JujuModelUUID() string {
	return env.Getter.Get(JujuModelUUIDEnvVar)
}
