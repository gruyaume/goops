package environment

const JujuUnitNameEnvVar = "JUJU_UNIT_NAME"

func (env Environment) JujuUnitName() string {
	return env.Getter.Get(JujuUnitNameEnvVar)
}
