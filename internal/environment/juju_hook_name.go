package environment

const JujuHookNameEnvVar = "JUJU_HOOK_NAME"

func JujuHookName(getter EnvironmentGetter) string {
	return getter.Get(JujuHookNameEnvVar)
}
