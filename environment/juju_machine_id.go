package environment

const JujuMachineIDEnvVar = "JUJU_MACHINE_ID"

func (env Environment) JujuMachineID() string {
	return env.Getter.Get(JujuMachineIDEnvVar)
}
