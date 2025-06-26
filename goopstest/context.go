package goopstest

import (
	"fmt"

	"github.com/gruyaume/goops"
)

type LogLevel string

const (
	DefaultJujuVersion = "3.6.0"
)

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

type MountMeta struct {
	Location string
	Storage  string
}

type ContainerMeta struct {
	Mounts   []MountMeta
	Resource string
}

type IntegrationMeta struct {
	Interface string
}

type ResourceMeta struct {
	Description string
	Type        string
	Filename    string
}

type StorageMeta struct {
	MinimumSize string
	Type        string
}

type Metadata struct {
	Containers  map[string]ContainerMeta
	Description string
	Name        string
	Provides    map[string]IntegrationMeta
	Resources   map[string]ResourceMeta
	Storage     map[string]StorageMeta
	Summary     string
}

type Context struct {
	CharmFunc     func() error
	Metadata      Metadata
	AppName       string
	UnitID        string
	JujuVersion   string
	ActionResults map[string]string
	ActionError   error
	JujuLog       []JujuLogLine
	CharmErr      error
}

func (c *Context) Run(hookName string, state State) (State, error) {
	state.Relations = setRelationIDs(state.Relations)
	state.Relations = setUnitIDs(state.Relations)
	state.PeerRelations = setPeerRelationIDs(state.PeerRelations)

	if state.Model.Name == "" {
		state.Model.Name = "test-model"
	}

	if state.Model.UUID == "" {
		state.Model.UUID = "12345678-1234-5678-1234-567812345678"
	}

	nilStatus := Status{}
	if state.UnitStatus == nilStatus {
		state.UnitStatus = Status{
			Name: StatusUnknown,
		}
	}

	fakeCommand := &fakeCommandRunner{
		Output:        []byte(``),
		Err:           nil,
		Leader:        state.Leader,
		Config:        state.Config,
		Secrets:       state.Secrets,
		Relations:     state.Relations,
		PeerRelations: state.PeerRelations,
		Ports:         state.Ports,
		StoredState:   state.StoredState,
		AppName:       c.AppName,
		UnitID:        c.UnitID,
		Model:         state.Model,
		UnitStatus:    state.UnitStatus,
		AppStatus:     state.AppStatus,
		Metadata:      c.Metadata,
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

	err := c.CharmFunc()
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

func (c *Context) RunAction(actionName string, state State, params map[string]any) (State, error) {
	fakeCommandRunner := &fakeCommandRunner{
		Output:           []byte(``),
		Err:              nil,
		Leader:           state.Leader,
		Config:           state.Config,
		Secrets:          state.Secrets,
		Relations:        state.Relations,
		PeerRelations:    state.PeerRelations,
		ActionParameters: params,
		Ports:            state.Ports,
		StoredState:      state.StoredState,
		AppName:          c.AppName,
		UnitID:           c.UnitID,
		Model:            state.Model,
		UnitStatus:       state.UnitStatus,
		AppStatus:        state.AppStatus,
	}

	if state.Model.Name == "" {
		state.Model.Name = "test-model"
	}

	if state.Model.UUID == "" {
		state.Model.UUID = "12345678-1234-5678-1234-567812345678"
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

	err := c.CharmFunc()
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
func setUnitIDs(relations []Relation) []Relation {
	for _, relation := range relations {
		if relation.RemoteUnitsData == nil {
			relation.RemoteUnitsData = make(map[UnitID]DataBag)
		}

		if len(relation.RemoteUnitsData) == 0 {
			relation.RemoteUnitsData[UnitID(relation.RemoteAppName+"/0")] = DataBag{}
		}
	}

	return relations
}

// For each relation, we set the ID to: <name>:<number>
func setRelationIDs(relations []Relation) []Relation {
	for i := range relations {
		if relations[i].ID == "" {
			relations[i].ID = fmt.Sprintf("%s:%d", relations[i].Endpoint, i)
		}
	}

	return relations
}

// For each peer relation, we set the ID to: <name>:<number>
func setPeerRelationIDs(peerRelations []PeerRelation) []PeerRelation {
	for i, relation := range peerRelations {
		if relation.ID == "" {
			relation.ID = fmt.Sprintf("%s:%d", relation.Endpoint, i)
		}
	}

	return peerRelations
}
