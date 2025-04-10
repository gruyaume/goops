package environment

const JujuContextIDEnvVar = "JUJU_CONTEXT_ID"

func (env Environment) JujuContextID() string {
	return env.Getter.Get(JujuContextIDEnvVar)
}
