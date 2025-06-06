package goops_test

type FakeRunner struct {
	Command string
	Args    []string
	Output  []byte
	Err     error
}

func (f *FakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.Command = name
	f.Args = args

	return f.Output, f.Err
}
