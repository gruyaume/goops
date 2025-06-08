package goopstest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gruyaume/goops"
)

type Context struct {
	Charm         func() error
	ActionResults map[string]string
	ActionError   error
}

type fakeRunner struct {
	Command       string
	Args          []string
	Output        []byte
	Err           error
	Status        string
	Leader        bool
	Config        map[string]string
	Secrets       []Secret
	ActionResults map[string]string
	ActionError   error
}

func (f *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.Command = name
	f.Args = args

	switch name {
	case "status-set":
		f.handleStatusSet(args)
	case "is-leader":
		f.handleIsLeader()
	case "config-get":
		f.handleConfigGet(args)
	case "secret-get":
		f.handleSecretGet(args)
	case "secret-add":
		f.handleSecretAdd(args)
	case "secret-remove":
		f.handleSecretRemove(args)
	case "action-set":
		f.handleActionSet(args)
	case "action-fail":
		f.handleActionFail(args)
	}

	return f.Output, f.Err
}

func (f *fakeRunner) handleStatusSet(args []string) {
	f.Status = args[0]
}

func (f *fakeRunner) handleIsLeader() {
	if f.Leader {
		f.Output = []byte(`true`)
	} else {
		f.Output = []byte(`false`)
	}
}

func (f *fakeRunner) handleConfigGet(args []string) {
	if value, ok := f.Config[args[0]]; ok {
		f.Output = []byte(fmt.Sprintf(`"%s"`, value))
	} else {
		f.Output = []byte(`""`)
		f.Err = fmt.Errorf("config key %s not found", args[0])
	}
}

func (f *fakeRunner) handleSecretAdd(args []string) {
	content := make(map[string]string)

	var label string

	for _, arg := range args {
		if strings.HasPrefix(arg, "--label=") {
			label = strings.TrimPrefix(arg, "--label=")
		} else if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				content[parts[0]] = parts[1]
			}
		}
	}

	f.Secrets = append(f.Secrets, Secret{
		Label:   label,
		Content: content,
	})
}

func (f *fakeRunner) handleSecretGet(args []string) {
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

func (f *fakeRunner) handleSecretRemove(args []string) {
	for i, secret := range f.Secrets {
		if strings.Contains(args[0], secret.ID) || strings.Contains(args[0], "--label="+secret.Label) {
			f.Secrets = append(f.Secrets[:i], f.Secrets[i+1:]...)
			break
		}
	}
}

func (f *fakeRunner) handleActionSet(args []string) {
	f.ActionResults = make(map[string]string)

	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				f.ActionResults[parts[0]] = parts[1]
			}
		}
	}
}

func (f *fakeRunner) handleActionFail(args []string) {
	f.ActionError = fmt.Errorf("%s", strings.Join(args, " "))
}

type fakeGetter struct {
	HookName   string
	ActionName string
}

func (f *fakeGetter) Get(key string) string {
	switch key {
	case "JUJU_HOOK_NAME":
		return f.HookName
	case "JUJU_ACTION_NAME":
		return f.ActionName
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
	state.Secrets = fakeRunner.Secrets

	return state, nil
}

func (c *Context) RunAction(actionName string, state *State) (*State, error) {
	fakeRunner := &fakeRunner{
		Output:  []byte(``),
		Err:     nil,
		Leader:  state.Leader,
		Config:  state.Config,
		Secrets: state.Secrets,
	}

	fakeGetter := &fakeGetter{
		ActionName: actionName,
	}

	goops.SetRunner(fakeRunner)
	goops.SetEnvironment(fakeGetter)

	err := c.Charm()
	if err != nil {
		return nil, err
	}

	state.UnitStatus = fakeRunner.Status
	state.Secrets = fakeRunner.Secrets
	c.ActionResults = fakeRunner.ActionResults
	c.ActionError = fakeRunner.ActionError

	return state, nil
}

type Secret struct {
	ID      string
	Label   string
	Content map[string]string
}

type State struct {
	Leader     bool
	UnitStatus string
	Config     map[string]string
	Secrets    []Secret
}
