package environment

const JujuPrincipalUnitEnvVar = "JUJU_PRINCIPAL_UNIT"

func (env Environment) JujuPrincipalUnit() string {
	return env.Getter.Get(JujuPrincipalUnitEnvVar)
}
