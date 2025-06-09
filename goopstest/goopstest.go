package goopstest

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	UnitStatus         string
	AppStatus          string
	Leader             bool
	Config             map[string]string
	Secrets            []*Secret
	ActionResults      map[string]string
	ActionParameters   map[string]string
	ActionError        error
	ApplicationVersion string
	Relations          []*Relation
	Ports              []*Port
}

func (f *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.Output = []byte(``)
	f.Err = nil
	f.Command = name
	f.Args = args

	handlers := map[string]func([]string){
		"action-fail":             f.handleActionFail,
		"action-get":              f.handleActionGet,
		"action-set":              f.handleActionSet,
		"application-version-set": f.handleApplicationVersionSet,
		"config-get":              f.handleConfigGet,
		"is-leader":               f.handleIsLeader,
		"opened-ports":            f.handleOpenedPorts,
		"open-port":               f.handleOpenPort,
		"close-port":              f.handleClosePort,
		"relation-ids":            f.handleRelationIDs,
		"relation-get":            f.handleRelationGet,
		"relation-list":           f.handleRelationList,
		"relation-set":            f.handleRelationSet,
		"secret-add":              f.handleSecretAdd,
		"secret-get":              f.handleSecretGet,
		"secret-remove":           f.handleSecretRemove,
		"status-set":              f.handleStatusSet,
	}

	if handler, exists := handlers[name]; exists {
		handler(args)
		return f.Output, f.Err
	}

	return nil, fmt.Errorf("unknown command: %s", name)
}

func (f *fakeRunner) handleStatusSet(args []string) {
	if len(args) == 0 {
		f.Err = fmt.Errorf("status-set command requires at least one argument")
		return
	}

	if args[0] == "--application" {
		if len(args) < 2 {
			f.Err = fmt.Errorf("status-set command requires an application status after --application")
			return
		}

		f.AppStatus = args[1]
	} else {
		f.UnitStatus = args[0]
	}
}

func (f *fakeRunner) handleIsLeader(_ []string) {
	if f.Leader {
		f.Output = []byte(`true`)
	} else {
		f.Output = []byte(`false`)
	}
}

func (f *fakeRunner) handleOpenedPorts(_ []string) {
	portList := make([]string, len(f.Ports))
	for i, port := range f.Ports {
		portList[i] = fmt.Sprintf("%d/%s", port.Port, port.Protocol)
	}

	output, err := json.Marshal(portList)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal opened ports: %w", err)
		return
	}

	f.Output = output
}

func (f *fakeRunner) handleOpenPort(args []string) {
	if len(args) != 1 {
		f.Err = fmt.Errorf("open-port command requires exactly one argument")
		return
	}

	portInfo := strings.Split(args[0], "/")

	if len(portInfo) != 2 {
		f.Err = fmt.Errorf("invalid port format, expected <port>/<protocol>")
		return
	}

	port, err := strconv.Atoi(portInfo[0])
	if err != nil || port < 0 || port > 65535 {
		f.Err = fmt.Errorf("invalid port number: %s", portInfo[0])
		return
	}

	protocol := portInfo[1]

	if protocol != "tcp" && protocol != "udp" {
		f.Err = fmt.Errorf("invalid protocol: %s, must be 'tcp' or 'udp'", protocol)
		return
	}

	f.Ports = append(f.Ports, &Port{
		Port:     port,
		Protocol: protocol,
	})
}

func (f *fakeRunner) handleClosePort(args []string) {
	if len(args) != 1 {
		f.Err = fmt.Errorf("close-port command requires exactly one argument")
		return
	}

	portInfo := strings.Split(args[0], "/")

	if len(portInfo) != 2 {
		f.Err = fmt.Errorf("invalid port format, expected <port>/<protocol>")
		return
	}

	port, err := strconv.Atoi(portInfo[0])
	if err != nil || port < 0 || port > 65535 {
		f.Err = fmt.Errorf("invalid port number: %s", portInfo[0])
		return
	}

	protocol := portInfo[1]

	if protocol != "tcp" && protocol != "udp" {
		f.Err = fmt.Errorf("invalid protocol: %s, must be 'tcp' or 'udp'", protocol)
		return
	}

	for i, p := range f.Ports {
		if p.Port == port && p.Protocol == protocol {
			f.Ports = append(f.Ports[:i], f.Ports[i+1:]...)
			return
		}
	}

	f.Err = fmt.Errorf("port %d/%s not found", port, protocol)
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
			if relation.ID != "" {
				f.Output = []byte(fmt.Sprintf(`["%s"]`, relation.ID))
			} else {
				f.Output = []byte(`[]`)
			}
		}
	}
}

func safeCopy(data DataBag) DataBag {
	if data == nil {
		return make(DataBag)
	}

	return data
}

func (f *fakeRunner) findRelationByID(id string) *Relation {
	for i := range f.Relations {
		if f.Relations[i].ID == id {
			return f.Relations[i]
		}
	}

	return nil
}

func parseRelationGetArgs(args []string) (isApp bool, relationID string, unitID string, err error) {
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--app":
			isApp = true
		case strings.HasPrefix(args[i], "-r="):
			relationID = strings.TrimPrefix(args[i], "-r=")
		case args[i] == "-" && i+1 < len(args):
			unitID = args[i+1]
		}
	}

	if relationID == "" || unitID == "" {
		return false, "", "", fmt.Errorf("relation ID or unit ID not provided")
	}

	return isApp, relationID, unitID, nil
}

func (f *fakeRunner) handleRelationGet(args []string) {
	isApp, relationID, unitID, err := parseRelationGetArgs(args)
	if err != nil {
		f.Err = err
		return
	}

	relation := f.findRelationByID(relationID)
	if relation == nil {
		f.Err = fmt.Errorf("relation %s not found", relationID)
		return
	}

	var data any

	if isApp {
		data = safeCopy(relation.RemoteAppData)
	} else {
		unitData, ok := relation.RemoteUnitsData[UnitID(unitID)]
		if !ok {
			f.Err = fmt.Errorf("unit ID %s not found in relation %s", unitID, relationID)
			return
		}

		data = unitData
	}

	output, err := json.Marshal(data)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal relation data: %w", err)
		return
	}

	f.Output = output
}

func (f *fakeRunner) handleRelationList(args []string) {
	relationID := strings.TrimPrefix(args[0], "-r=")

	for _, relation := range f.Relations {
		if relation.ID == relationID {
			unitIDs := make([]string, 0, len(relation.RemoteUnitsData))
			for unitID := range relation.RemoteUnitsData {
				unitIDs = append(unitIDs, string(unitID))
			}

			output, err := json.Marshal(unitIDs)
			if err != nil {
				f.Err = fmt.Errorf("failed to marshal relation units: %w", err)
				return
			}

			f.Output = output

			return
		}
	}
}

func parseRelationSetArgs(args []string) (isApp bool, relationID string, data map[string]string, err error) {
	filteredArgs := make([]string, 0, len(args))

	for _, arg := range args {
		switch {
		case arg == "--app":
			isApp = true
		case strings.HasPrefix(arg, "-r="):
			relationID = strings.TrimPrefix(arg, "-r=")
		default:
			filteredArgs = append(filteredArgs, arg)
		}
	}

	if relationID == "" {
		return false, "", nil, fmt.Errorf("relation ID not provided")
	}

	data = parseKeyValueArgs(filteredArgs)

	return isApp, relationID, data, nil
}

func (f *fakeRunner) handleRelationSet(args []string) {
	isApp, relationID, data, err := parseRelationSetArgs(args)
	if err != nil {
		f.Err = err
		return
	}

	for _, relation := range f.Relations {
		if relation.ID != relationID {
			continue
		}

		target := &relation.LocalUnitData
		if isApp {
			target = &relation.LocalAppData
		}

		if *target == nil {
			*target = make(DataBag)
		}

		for k, v := range data {
			(*target)[k] = v
		}
	}
}

func filterOutLabelArgs(args []string) []string {
	filtered := make([]string, 0, len(args))

	for _, arg := range args {
		if !strings.HasPrefix(arg, "--label=") {
			filtered = append(filtered, arg)
		}
	}

	return filtered
}

func (f *fakeRunner) handleSecretAdd(args []string) {
	label := extractLabelFromArgs(args)
	filtered := filterOutLabelArgs(args)

	content := parseKeyValueArgs(filtered)

	f.Secrets = append(f.Secrets, &Secret{
		Label:   label,
		Content: content,
	})
}

func (f *fakeRunner) handleSecretGet(args []string) {
	label := extractLabelFromArgs(args)
	if label == "" {
		f.Err = fmt.Errorf("no --label specified")
		return
	}

	secret := findSecretByLabel(f.Secrets, label)
	if secret == nil {
		f.Err = fmt.Errorf("secret with label %q not found", label)
		return
	}

	output, err := json.Marshal(secret.Content)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal secret content: %w", err)
		return
	}

	f.Output = output
}

// extractLabelFromArgs returns the label from args if present.
func extractLabelFromArgs(args []string) string {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--label=") {
			return strings.TrimPrefix(arg, "--label=")
		}
	}

	return ""
}

// findSecretByLabel returns the pointer to the secret with the given label.
func findSecretByLabel(secrets []*Secret, label string) *Secret {
	for _, secret := range secrets {
		if secret.Label == label {
			return secret
		}
	}

	return nil
}

func (f *fakeRunner) handleSecretRemove(args []string) {
	for i, secret := range f.Secrets {
		if strings.Contains(args[0], secret.ID) || strings.Contains(args[0], "--label="+secret.Label) {
			f.Secrets = append(f.Secrets[:i], f.Secrets[i+1:]...)
			break
		}
	}
}

func parseKeyValueArgs(args []string) map[string]string {
	result := make(map[string]string)

	for _, arg := range args {
		if parts := strings.SplitN(arg, "=", 2); len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}

	return result
}

func (f *fakeRunner) handleActionSet(args []string) {
	f.ActionResults = parseKeyValueArgs(args)
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
	Model      *Model
}

func (f *fakeGetter) Get(key string) string {
	switch key {
	case "JUJU_HOOK_NAME":
		return f.HookName
	case "JUJU_ACTION_NAME":
		return f.ActionName
	case "JUJU_MODEL_NAME":
		return f.Model.Name
	case "JUJU_MODEL_UUID":
		return f.Model.UUID
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

func (c *Context) Run(hookName string, state *State) (*State, error) {
	setRelationIDs(state.Relations)
	setUnitIDs(state.Relations)

	if state.Model == nil {
		state.Model = &Model{
			Name: "test-model",
			UUID: "12345678-1234-5678-1234-567812345678",
		}
	}

	fakeRunner := &fakeRunner{
		Output:    []byte(``),
		Err:       nil,
		Leader:    state.Leader,
		Config:    state.Config,
		Secrets:   state.Secrets,
		Relations: state.Relations,
		Ports:     state.Ports,
	}

	fakeGetter := &fakeGetter{
		HookName: hookName,
		Model:    state.Model,
	}

	goops.SetRunner(fakeRunner)
	goops.SetEnvironment(fakeGetter)

	err := c.Charm()
	if err != nil {
		return nil, fmt.Errorf("failed to run charm: %w", err)
	}

	state.UnitStatus = fakeRunner.UnitStatus
	state.AppStatus = fakeRunner.AppStatus
	state.Secrets = fakeRunner.Secrets
	state.ApplicationVersion = fakeRunner.ApplicationVersion
	state.Ports = fakeRunner.Ports

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

	if state.Model == nil {
		state.Model = &Model{
			Name: "test-model",
			UUID: "12345678-1234-5678-1234-567812345678",
		}
	}

	fakeGetter := &fakeGetter{
		ActionName: actionName,
		Model:      state.Model,
	}

	goops.SetRunner(fakeRunner)
	goops.SetEnvironment(fakeGetter)

	err := c.Charm()
	if err != nil {
		return nil, err
	}

	state.UnitStatus = fakeRunner.UnitStatus
	state.AppStatus = fakeRunner.AppStatus
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

type UnitID string

type DataBag map[string]string

type Relation struct {
	Endpoint        string
	Interface       string
	ID              string
	RemoteAppName   string
	LocalAppData    DataBag
	LocalUnitData   DataBag
	RemoteAppData   DataBag
	RemoteUnitsData map[UnitID]DataBag
}

type Port struct {
	Port     int
	Protocol string
}

type Model struct {
	Name string
	UUID string
}

type State struct {
	Leader             bool
	UnitStatus         string
	AppStatus          string
	Config             map[string]string
	Secrets            []*Secret
	ApplicationVersion string
	Relations          []*Relation
	Ports              []*Port
	Model              *Model
}
