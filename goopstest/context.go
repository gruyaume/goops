package goopstest

import (
	"fmt"

	"github.com/gruyaume/goops"
)

type LogLevel string

const (
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
	LogLevelDebug   LogLevel = "DEBUG"
)

type JujuLogLine struct {
	Level   LogLevel
	Message string
}

type Context struct {
	Charm         func() error
	Metadata      goops.Metadata
	AppName       string
	UnitID        int
	JujuVersion   string
	ActionResults map[string]string
	ActionError   error
	JujuLog       []JujuLogLine
	CharmErr      error
}

func (c *Context) Run(hookName string, state *State) (*State, error) {
	if c.Charm == nil {
		return nil, fmt.Errorf("charm function is not set in the context")
	}

	setRelationIDs(state.Relations)
	setUnitIDs(state.Relations)

	if state.Model == nil {
		state.Model = &Model{
			Name: "test-model",
			UUID: "12345678-1234-5678-1234-567812345678",
		}
	}

	fakeCommand := &fakeCommandRunner{
		Output:      []byte(``),
		Err:         nil,
		Leader:      state.Leader,
		Config:      state.Config,
		Secrets:     state.Secrets,
		Relations:   state.Relations,
		Ports:       state.Ports,
		StoredState: state.StoredState,
		AppName:     c.AppName,
		UnitID:      c.UnitID,
		Model:       state.Model,
	}

	fakeEnv := &fakeEnvGetter{
		HookName:    hookName,
		Model:       state.Model,
		AppName:     c.AppName,
		UnitID:      c.UnitID,
		JujuVersion: c.JujuVersion,
		Metadata:    c.Metadata,
	}

	fakePebble := &fakePebbleGetter{
		Containers: state.Containers,
	}

	goops.SetPebbleGetter(fakePebble)
	goops.SetCommandRunner(fakeCommand)
	goops.SetEnvGetter(fakeEnv)

	err := c.Charm()
	if err != nil {
		c.CharmErr = err
	}

	state.UnitStatus = fakeCommand.UnitStatus
	state.AppStatus = fakeCommand.AppStatus
	state.Secrets = fakeCommand.Secrets
	state.ApplicationVersion = fakeCommand.ApplicationVersion
	state.Ports = fakeCommand.Ports
	state.StoredState = fakeCommand.StoredState
	state.Containers = fakePebble.Containers

	c.JujuLog = fakeCommand.JujuLog

	return state, nil
}

func (c *Context) RunAction(actionName string, state *State, params map[string]any) (*State, error) {
	fakeCommandRunner := &fakeCommandRunner{
		Output:           []byte(``),
		Err:              nil,
		Leader:           state.Leader,
		Config:           state.Config,
		Secrets:          state.Secrets,
		ActionParameters: params,
		StoredState:      state.StoredState,
	}

	if state.Model == nil {
		state.Model = &Model{
			Name: "test-model",
			UUID: "12345678-1234-5678-1234-567812345678",
		}
	}

	fakeEnvGetter := &fakeEnvGetter{
		ActionName:  actionName,
		Model:       state.Model,
		AppName:     c.AppName,
		UnitID:      c.UnitID,
		JujuVersion: c.JujuVersion,
	}

	goops.SetCommandRunner(fakeCommandRunner)
	goops.SetEnvGetter(fakeEnvGetter)

	err := c.Charm()
	if err != nil {
		c.CharmErr = err
	}

	state.UnitStatus = fakeCommandRunner.UnitStatus
	state.AppStatus = fakeCommandRunner.AppStatus
	state.Secrets = fakeCommandRunner.Secrets
	state.ApplicationVersion = fakeCommandRunner.ApplicationVersion
	c.ActionResults = fakeCommandRunner.ActionResults
	c.ActionError = fakeCommandRunner.ActionError

	return state, nil
}

// For each relation, we set the remoteUnitsData so that it contains at leader 1 unit
func setUnitIDs(relations []*Relation) {
	for _, relation := range relations {
		if relation.RemoteUnitsData == nil {
			relation.RemoteUnitsData = make(map[UnitID]DataBag)
		}

		if len(relation.RemoteUnitsData) == 0 {
			relation.RemoteUnitsData[UnitID(relation.RemoteAppName+"/0")] = DataBag{}
		}
	}
}

// For each relation, we set the ID to: <name>:<number>
func setRelationIDs(relations []*Relation) {
	for i, relation := range relations {
		if relation.ID == "" {
			relation.ID = fmt.Sprintf("%s:%d", relation.Endpoint, i)
		}
	}
}
