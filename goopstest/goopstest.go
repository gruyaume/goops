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
	Command            string
	Args               []string
	Output             []byte
	Err                error
	Status             string
	Leader             bool
	Config             map[string]string
	Secrets            []*Secret
	ActionResults      map[string]string
	ActionParameters   map[string]string
	ActionError        error
	ApplicationVersion string
	Relations          []*Relation
}

func (f *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.Command = name
	f.Args = args

	switch name {
	case "action-fail":
		f.handleActionFail(args)
	case "action-get":
		f.handleActionGet(args)
	case "action-log":
		// Not yet implemented
	case "action-set":
		f.handleActionSet(args)
	case "application-version-set":
		f.handleApplicationVersionSet(args)
	case "config-get":
		f.handleConfigGet(args)
	case "credential-get":
		// Not yet implemented
	case "goal-state":
		// Not yet implemented
	case "is-leader":
		f.handleIsLeader()
	case "juju-log":
		// Not yet implemented
	case "network-get":
		// Not yet implemented
	case "open-port":
		// Not yet implemented
	case "close-port":
		// Not yet implemented
	case "opened-ports":
		// Not yet implemented
	case "juju-reboot":
		// Not yet implemented
	case "relation-ids":
		f.handleRelationIDs(args)
	case "relation-get":
		// Not yet implemented
	case "relation-list":
		// Not yet implemented
	case "relation-set":
		// Not yet implemented
	case "relation-model-get":
		// Not yet implemented
	case "resource-get":
		// Not yet implemented
	case "secret-add":
		f.handleSecretAdd(args)
	case "secret-get":
		f.handleSecretGet(args)
	case "secret-grant":
		// Not yet implemented
	case "secret-ids":
		// Not yet implemented
	case "secret-info-get":
		// Not yet implemented
	case "secret-remove":
		f.handleSecretRemove(args)
	case "secret-revoke":
		// Not yet implemented
	case "secret-set":
		// Not yet implemented
	case "state-delete":
		// Not yet implemented
	case "state-get":
		// Not yet implemented
	case "state-set":
		// Not yet implemented
	case "status-get":
		// Not yet implemented
	case "status-set":
		f.handleStatusSet(args)
	case "storage-add":
		// Not yet implemented
	case "storage-get":
		// Not yet implemented
	case "storage-list":
		// Not yet implemented
	case "unit-get":
		// Not yet implemented
	default:
		return nil, fmt.Errorf("unknown command: %s", name)
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

func (f *fakeRunner) handleRelationIDs(args []string) {
	for _, relation := range f.Relations {
		if len(args) > 0 && args[0] == relation.Endpoint {
			// If the endpoint matches, return the relation ID
			if relation.ID != "" {
				f.Output = []byte(fmt.Sprintf(`["%s"]`, relation.ID))
			} else {
				f.Output = []byte(`[]`)
			}
		}
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

	f.Secrets = append(f.Secrets, &Secret{
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

func (f *fakeRunner) handleActionGet(args []string) {
	key := args[0]

	if value, ok := f.ActionParameters[key]; ok {
		output, err := json.Marshal(value)
		if err != nil {
			f.Err = fmt.Errorf("failed to marshal action parameter: %w", err)
			return
		}

		f.Output = output
	} else {
		f.Err = fmt.Errorf("action parameter %s not found", key)
		f.Output = []byte(`""`)
	}
}

func (f *fakeRunner) handleApplicationVersionSet(args []string) {
	f.ApplicationVersion = args[0]
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

// For each relation, we set the ID to: <name>:<number>
func setRelationIDs(relations []*Relation) {
	for i, relation := range relations {
		if relation.ID == "" {
			relation.ID = fmt.Sprintf("%s:%d", relation.Endpoint, i)
		}
	}
}

func (c *Context) Run(hookName string, state *State) (*State, error) {
	setRelationIDs(state.Relations)
	fakeRunner := &fakeRunner{
		Output:    []byte(``),
		Err:       nil,
		Leader:    state.Leader,
		Config:    state.Config,
		Secrets:   state.Secrets,
		Relations: state.Relations,
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
	state.ApplicationVersion = fakeRunner.ApplicationVersion

	return state, nil
}

func (c *Context) RunAction(actionName string, state *State, params map[string]string) (*State, error) {
	fakeRunner := &fakeRunner{
		Output:           []byte(``),
		Err:              nil,
		Leader:           state.Leader,
		Config:           state.Config,
		Secrets:          state.Secrets,
		ActionParameters: params,
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

type Relation struct {
	Endpoint      string
	Interface     string
	ID            string
	LocalAppData  map[string]string
	LocalUnitData map[string]string
}

type State struct {
	Leader             bool
	UnitStatus         string
	Config             map[string]string
	Secrets            []*Secret
	ApplicationVersion string
	Relations          []*Relation
}
