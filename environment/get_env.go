package environment

import "os"

type EnvironmentGetter interface {
	Get(name string) string
}

type ExecutionEnvironment struct{}

func (r *ExecutionEnvironment) Get(name string) string {
	return os.Getenv(name)
}

type Environment struct {
	Getter EnvironmentGetter
}

func NewEnvironment() *Environment {
	return &Environment{
		Getter: &ExecutionEnvironment{},
	}
}
