package environment

const JujuHookNameEnvVar = "JUJU_HOOK_NAME"

func (env Environment) JujuHookName() string {
	return env.Getter.Get(JujuHookNameEnvVar)
}
