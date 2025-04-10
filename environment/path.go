package environment

const PathEnvVar = "PATH"

func (env Environment) Path() string {
	return env.Getter.Get(PathEnvVar)
}
