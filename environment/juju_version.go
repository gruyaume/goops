package environment

const JujuVersionEnvVar = "JUJU_VERSION"

func (env Environment) JujuVersion() string {
	return env.Getter.Get(JujuVersionEnvVar)
}
