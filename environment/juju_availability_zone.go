package environment

const JujuAvailabilityZoneEnvVar = "JUJU_AVAILABILITY_ZONE"

func (env Environment) JujuAvailabilityZone() string {
	return env.Getter.Get(JujuAvailabilityZoneEnvVar)
}
