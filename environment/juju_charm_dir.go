package environment

const JujuCharmDirEnvVar = "JUJU_CHARM_DIR"

func (env Environment) JujuCharmDir() string {
	return env.Getter.Get(JujuCharmDirEnvVar)
}
