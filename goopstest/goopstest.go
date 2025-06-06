package goopstest

import (
	"fmt"

	"github.com/gruyaume/goops"
)

type Context struct {
	Charm func() error
}

type FakeRunner struct {
	Command string
	Args    []string
	Output  []byte
	Err     error
	Status  string
}

func (f *FakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.Command = name
	f.Args = args

	if name == "status-set" {
		f.Status = args[0]
	}

	return f.Output, f.Err
}

type FakeGetter struct {
	HookName string
}

func (f *FakeGetter) Get(key string) string {
	if key == "JUJU_HOOK_NAME" {
		return f.HookName
	}

	return ""
}

func (c *Context) Run(event string, state State) (*State, error) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	fakeGetter := &FakeGetter{
		HookName: event,
	}

	goops.SetRunner(fakeRunner)
	goops.SetEnvironment(fakeGetter)

	err := c.Charm()
	if err != nil {
		return nil, fmt.Errorf("failed to run charm: %w", err)
	}

	return &State{
		UnitStatus: fakeRunner.Status,
	}, nil
}

type State struct {
	UnitStatus string
}
