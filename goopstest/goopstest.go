package goopstest

import (
	"encoding/json"
	"fmt"
	"strings"

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
	Config  map[string]string
	Secrets []Secret
}

func (f *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.Command = name
	f.Args = args

	switch name {
	case "status-set":
		f.Status = args[0]
	case "is-leader":
		if f.Leader {
			f.Output = []byte(`true`)
		} else {
			f.Output = []byte(`false`)
		}
	case "config-get":
		if value, ok := f.Config[args[0]]; ok {
			f.Output = []byte(fmt.Sprintf(`"%s"`, value))
		} else {
			f.Output = []byte(`""`)
			f.Err = fmt.Errorf("config key %s not found", args[0])
		}
	case "secret-get":
		for _, secret := range f.Secrets {
			if strings.Contains(args[0], "--label") && strings.Contains(args[0], "--label"+"="+secret.Label) {
				output, err := json.Marshal(secret.Content)
				if err != nil {
					f.Err = err
					break
				}

				f.Output = output

				break
			}
		}
	}

	return f.Output, f.Err
}

type fakeGetter struct {
	HookName string
}

func (f *fakeGetter) Get(key string) string {
	switch key {
	case "JUJU_HOOK_NAME":
		return f.HookName
	}

	return ""
}

func (c *Context) Run(hookName string, state *State) (*State, error) {
	fakeRunner := &fakeRunner{
		Output:  []byte(``),
		Err:     nil,
		Leader:  state.Leader,
		Config:  state.Config,
		Secrets: state.Secrets,
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

type Secret struct {
	Label   string
	Content map[string]string
}

type State struct {
	Leader     bool
	UnitStatus string
	Config     map[string]string
	Secrets    []Secret
}
