package environment

import "os"

type EnvironmentGetter interface {
	Get(name string) string
}

type ExecutionEnvironment struct{}

func (r *ExecutionEnvironment) Get(name string) string {
	return os.Getenv(name)
}
