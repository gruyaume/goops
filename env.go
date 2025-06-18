package goops

import (
	"os"
)

var defaultGetter EnvironmentGetter = &realExecutionEnvironment{}

type realExecutionEnvironment struct{}

type Environment struct {
	ActionName         string
	AgentSocketAddress string
	AgentSocketNetwork string
	APIAddresses       string
	AvailabilityZone   string
	CharmDir           string
	CharmFTPProxy      string
	CharmHTTPProxy     string
	CharmHTTPSProxy    string
	CharmNoProxy       string
	CloudAPIVersion    string
	ContextID          string
	HookName           string
	MachineID          string
	ModelName          string
	ModelUUID          string
	Path               string
	PrincipalUnit      string
	UnitName           string
	Version            string
}

type EnvironmentGetter interface {
	Get(name string) string
	ReadFile(name string) ([]byte, error)
}

func (r *realExecutionEnvironment) Get(name string) string {
	return os.Getenv(name)
}

func (r *realExecutionEnvironment) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name) // #nosec G304
}

type JujuEnvironment struct {
	Getter EnvironmentGetter
}

func GetEnvGetter() EnvironmentGetter {
	return defaultGetter
}

func SetEnvGetter(envGetter EnvironmentGetter) {
	defaultGetter = envGetter
}

func ReadEnv() Environment {
	envGetter := GetEnvGetter()

	return Environment{
		ActionName:         envGetter.Get("JUJU_ACTION_NAME"),
		AgentSocketAddress: envGetter.Get("JUJU_AGENT_SOCKET_ADDRESS"),
		AgentSocketNetwork: envGetter.Get("JUJU_AGENT_SOCKET_NETWORK"),
		APIAddresses:       envGetter.Get("JUJU_API_ADDRESSES"),
		AvailabilityZone:   envGetter.Get("JUJU_AVAILABILITY_ZONE"),
		CloudAPIVersion:    envGetter.Get("CLOUD_API_VERSION"),
		CharmDir:           envGetter.Get("JUJU_CHARM_DIR"),
		CharmFTPProxy:      envGetter.Get("JUJU_CHARM_FTP_PROXY"),
		CharmHTTPProxy:     envGetter.Get("JUJU_CHARM_HTTP_PROXY"),
		CharmHTTPSProxy:    envGetter.Get("JUJU_CHARM_HTTPS_PROXY"),
		CharmNoProxy:       envGetter.Get("JUJU_CHARM_NO_PROXY"),
		ContextID:          envGetter.Get("JUJU_CONTEXT_ID"),
		HookName:           envGetter.Get("JUJU_HOOK_NAME"),
		MachineID:          envGetter.Get("JUJU_MACHINE_ID"),
		ModelName:          envGetter.Get("JUJU_MODEL_NAME"),
		ModelUUID:          envGetter.Get("JUJU_MODEL_UUID"),
		PrincipalUnit:      envGetter.Get("JUJU_PRINCIPAL_UNIT"),
		UnitName:           envGetter.Get("JUJU_UNIT_NAME"),
		Version:            envGetter.Get("JUJU_VERSION"),
		Path:               envGetter.Get("PATH"),
	}
}
