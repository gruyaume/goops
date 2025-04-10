package environment

const JujuAPIAddressesEnvVar = "JUJU_API_ADDRESSES"

func (env Environment) JujuAPIAddresses() string {
	return env.Getter.Get(JujuAPIAddressesEnvVar)
}
