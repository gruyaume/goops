package goopstest

import (
	"fmt"

	"github.com/gruyaume/goops"
)

type Context struct {
	Charm func() error
}

type fakeRunner struct {
	Command string
	Args    []string
	Output  []byte
	Err     error
	Status  string
	Leader  bool
}

func (f *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.Command = name
	f.Args = args

	if name == "status-set" {
		f.Status = args[0]
	}

	if name == "is-leader" {
		if f.Leader {
			f.Output = []byte(`true`)
		} else {
			f.Output = []byte(`false`)
		}
	}

	return f.Output, f.Err
}

type fakeGetter struct {
	HookName string
}

func (f *fakeGetter) Get(key string) string {
	if key == "JUJU_HOOK_NAME" {
		return f.HookName
	}

	return ""
}

func (c *Context) Run(hookName string, state *State) (*State, error) {
	fakeRunner := &fakeRunner{
		Output: []byte(``),
		Err:    nil,
		Leader: state.Leader,
	}

	fakeGetter := &fakeGetter{
		HookName: hookName,
	}

	goops.SetRunner(fakeRunner)
	goops.SetEnvironment(fakeGetter)

	err := c.Charm()
	if err != nil {
		return nil, fmt.Errorf("failed to run charm: %w", err)
	}

	state.UnitStatus = fakeRunner.Status

	return state, nil
}

type State struct {
	Leader     bool
	UnitStatus string
}
