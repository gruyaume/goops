package environment

const CloudAPIVersionEnvVar = "CLOUD_API_VERSION"

func (env Environment) CloudAPIVersion() string {
	return env.Getter.Get(CloudAPIVersionEnvVar)
}
