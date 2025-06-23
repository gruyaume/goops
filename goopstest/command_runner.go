package goopstest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gruyaume/goops"
)

type fakeCommandRunner struct {
	Command            string
	Args               []string
	Output             []byte
	Err                error
	UnitStatus         string
	AppStatus          string
	Leader             bool
	Config             map[string]any
	Secrets            []*Secret
	ActionResults      map[string]string
	ActionParameters   map[string]any
	ActionError        error
	ApplicationVersion string
	Relations          []*Relation
	Ports              []*Port
	StoredState        StoredState
	AppName            string
	UnitID             int
	JujuLog            []JujuLogLine
}

func (f *fakeCommandRunner) Run(name string, args ...string) ([]byte, error) {
	f.Output = []byte(``)
	f.Err = nil
	f.Command = name
	f.Args = args

	handlers := map[string]func([]string){
		"action-fail":             f.handleActionFail,
		"action-get":              f.handleActionGet,
		"action-set":              f.handleActionSet,
		"action-log":              f.handleActionLog,
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
		"secret-info-get":         f.handleSecretInfoGet,
		"secret-ids":              f.handleSecretIDs,
		"secret-grant":            f.handleSecretGrant,
		"secret-set":              f.handleSecretSet,
		"secret-revoke":           f.handleSecretRevoke,
		"state-get":               f.handleStateGet,
		"state-set":               f.handleStateSet,
		"state-delete":            f.handleStateDelete,
		"status-set":              f.handleStatusSet,
		"juju-log":                f.handleJujuLog,
	}

	if handler, exists := handlers[name]; exists {
		handler(args)
		return f.Output, f.Err
	}

	return nil, fmt.Errorf("unknown command: %s", name)
}

func (f *fakeCommandRunner) handleStatusSet(args []string) {
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

func (f *fakeCommandRunner) handleJujuLog(args []string) {
	var logLevel LogLevel

	if len(args) > 0 {
		switch args[0] {
		case "--log-level=INFO":
			logLevel = LogLevelInfo
		case "--log-level=WARNING":
			logLevel = LogLevelWarning
		case "--log-level=ERROR":
			logLevel = LogLevelError
		case "--log-level=DEBUG":
			logLevel = LogLevelDebug
		default:
			logLevel = LogLevelInfo
		}

		args = args[1:]
	}

	message := strings.Join(args, " ")
	newLogEntry := JujuLogLine{
		Level:   logLevel,
		Message: message,
	}
	f.JujuLog = append(f.JujuLog, newLogEntry)
}

func (f *fakeCommandRunner) handleIsLeader(_ []string) {
	if f.Leader {
		f.Output = []byte(`true`)
	} else {
		f.Output = []byte(`false`)
	}
}

func (f *fakeCommandRunner) handleOpenedPorts(_ []string) {
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

func (f *fakeCommandRunner) handleOpenPort(args []string) {
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

func (f *fakeCommandRunner) handleClosePort(args []string) {
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

func (f *fakeCommandRunner) handleConfigGet(_ []string) {
	if len(f.Config) == 0 {
		f.Output = []byte(`{}`)
		return
	}

	output, err := json.Marshal(f.Config)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal config: %w", err)
		return
	}

	f.Output = output
}

func (f *fakeCommandRunner) handleRelationIDs(args []string) {
	if args[0] == "" {
		f.Err = fmt.Errorf("command relation-ids failed: ERROR no endpoint name specified")

		return
	}

	if len(f.Relations) == 0 {
		f.Output = []byte(`[]`)
		return
	}

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

func (f *fakeCommandRunner) findRelationByID(id string) *Relation {
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

func (f *fakeCommandRunner) handleRelationGet(args []string) {
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

	isLocal := unitID == f.AppName+"/"+strconv.Itoa(f.UnitID)

	data, err := f.selectRelationData(relation, isApp, isLocal, unitID)
	if err != nil {
		f.Err = err
		return
	}

	f.Output, err = json.Marshal(data)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal relation data: %w", err)
	}
}

func (f *fakeCommandRunner) selectRelationData(rel *Relation, isApp, isLocal bool, unitID string) (any, error) {
	if isApp {
		if isLocal {
			return safeCopy(rel.LocalAppData), nil
		}

		return safeCopy(rel.RemoteAppData), nil
	}

	if isLocal {
		if rel.LocalUnitData == nil {
			return nil, fmt.Errorf("local unit data not found for relation %s", rel.ID)
		}

		return safeCopy(rel.LocalUnitData), nil
	}

	unitData, ok := rel.RemoteUnitsData[UnitID(unitID)]
	if !ok {
		return nil, fmt.Errorf("unit ID %s not found in relation %s", unitID, rel.ID)
	}

	return unitData, nil
}

func (f *fakeCommandRunner) handleRelationList(args []string) {
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

func (f *fakeCommandRunner) handleRelationSet(args []string) {
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

func (f *fakeCommandRunner) handleSecretAdd(args []string) {
	label, owner, description, rotation, expiry := extractSecretArgs(args)
	filtered := filterOutLabelArgs(args)

	if !f.Leader && owner != "unit" {
		f.Err = fmt.Errorf("command secret-add failed: ERROR this unit is not the leader")
		return
	}

	var expiryTime time.Time

	var err error

	if expiry != "" {
		expiryTime, err = time.Parse(time.RFC3339, expiry)
		if err != nil {
			f.Err = fmt.Errorf("invalid expiry format: %w", err)
			return
		}
	}

	content := parseKeyValueArgs(filtered)

	f.Secrets = append(f.Secrets, &Secret{
		Label:       label,
		Content:     content,
		Owner:       owner,
		Description: description,
		Rotate:      rotation,
		Expire:      expiryTime,
	})
}

func (f *fakeCommandRunner) handleSecretGet(args []string) {
	var label, id string

	for _, arg := range args {
		if strings.HasPrefix(arg, "--label=") {
			label = strings.TrimPrefix(arg, "--label=")
			break
		}
	}

	if label != "" {
		secret := findSecretByLabel(f.Secrets, label)
		if secret == nil {
			f.Err = fmt.Errorf("secret with label %q not found", label)
			return
		}

		f.setSecretOutput(secret)

		return
	}

	// No label; try extracting ID from positional args
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			id = arg
			break
		}
	}

	if id == "" {
		f.Err = fmt.Errorf("no --label or ID specified")
		return
	}

	secret := findSecretByID(f.Secrets, id)
	if secret == nil {
		f.Err = fmt.Errorf("secret with ID %q not found", id)
		return
	}

	f.setSecretOutput(secret)
}

func (f *fakeCommandRunner) setSecretOutput(secret *Secret) {
	output, err := json.Marshal(secret.Content)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal secret content: %w", err)
		return
	}

	f.Output = output
}

func findSecretByID(secrets []*Secret, id string) *Secret {
	for _, s := range secrets {
		if s.ID == id {
			return s
		}
	}

	return nil
}

func extractSecretArgs(args []string) (label, owner, desc, rotate, expire string) {
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--label="):
			label = strings.TrimPrefix(arg, "--label=")
		case strings.HasPrefix(arg, "--owner="):
			owner = strings.TrimPrefix(arg, "--owner=")
		case strings.HasPrefix(arg, "--description="):
			desc = strings.TrimPrefix(arg, "--description=")
		case strings.HasPrefix(arg, "--rotate="):
			rotate = strings.TrimPrefix(arg, "--rotate=")
		case strings.HasPrefix(arg, "--expire="):
			expire = strings.TrimPrefix(arg, "--expire=")
		}
	}

	return
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

func (f *fakeCommandRunner) handleSecretRemove(args []string) {
	if !f.Leader {
		return
	}

	for i, secret := range f.Secrets {
		if strings.Contains(args[0], secret.ID) || strings.Contains(args[0], "--label="+secret.Label) {
			f.Secrets = append(f.Secrets[:i], f.Secrets[i+1:]...)
			break
		}
	}
}

type SecretInfo struct {
	Revision    int    `json:"revision"`
	Label       string `json:"label"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	Rotation    string `json:"rotation"`
	Expiry      string `json:"expiry"`
}

func (f *fakeCommandRunner) handleSecretInfoGet(args []string) {
	var id, label string

	for _, arg := range args {
		if strings.HasPrefix(arg, "--label=") {
			label = strings.TrimPrefix(arg, "--label=")
		} else if !strings.HasPrefix(arg, "--") {
			id = arg
		}
	}

	var secret *Secret

	switch {
	case label != "":
		secret = findSecretByLabel(f.Secrets, label)

		if secret == nil {
			f.Err = fmt.Errorf(`ERROR secret %q not found`, label)

			return
		}

		if !f.Leader && secret.Owner != "unit" {
			f.Err = fmt.Errorf(`ERROR secret %q not found`, label)
			return
		}
	case id != "":
		secret = findSecretByID(f.Secrets, id)
		if secret == nil {
			f.Err = fmt.Errorf(`ERROR secret %q not found`, id)
			return
		}

		if !f.Leader && secret.Owner != "unit" {
			f.Err = fmt.Errorf(`ERROR secret %q not found`, id)
			return
		}
	default:
		f.Err = fmt.Errorf("no --label or ID specified")
		return
	}

	secretInfo := map[string]SecretInfo{
		secret.ID: {
			Revision:    1,
			Label:       secret.Label,
			Owner:       secret.Owner,
			Description: secret.Description,
			Rotation:    secret.Rotate,
			Expiry:      secret.Expire.Format(time.RFC3339),
		},
	}

	output, err := json.Marshal(secretInfo)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal secret info: %w", err)
		return
	}

	f.Output = output
}

func (f *fakeCommandRunner) handleSecretIDs(_ []string) {
	if len(f.Secrets) == 0 {
		f.Output = []byte(`null`)
		return
	}

	ids := []string{}

	for _, secret := range f.Secrets {
		if !f.Leader && secret.Owner != "unit" {
			continue
		}

		ids = append(ids, secret.ID)
	}

	output, err := json.Marshal(ids)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal secret IDs: %w", err)
		return
	}

	f.Output = output
}

func (f *fakeCommandRunner) handleSecretGrant(args []string) {
	if len(args) == 0 {
		f.Err = fmt.Errorf("secret-grant command requires at least one argument")
		return
	}

	secretID := args[0]

	if !f.Leader {
		f.Err = fmt.Errorf(`ERROR secret "%s" not found`, secretID)
		return
	}
}

func (f *fakeCommandRunner) handleSecretSet(args []string) {
	if len(args) == 0 {
		f.Err = fmt.Errorf("secret-set command requires at least one argument")
		return
	}

	if !f.Leader {
		f.Output = []byte(`null`)
		return
	}

	id := args[0]
	args = args[1:]

	meta, remaining := parseSecretMetadata(args)

	content := make(map[string]string)

	for _, arg := range remaining {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			f.Err = fmt.Errorf("invalid secret-set argument: %s", arg)
			return
		}

		content[parts[0]] = parts[1]
	}

	for _, secret := range f.Secrets {
		if secret.ID != id {
			continue
		}

		secret.Content = content
		if meta["label"] != "" {
			secret.Label = meta["label"]
		}

		if meta["owner"] != "" {
			secret.Owner = meta["owner"]
		}

		if meta["description"] != "" {
			secret.Description = meta["description"]
		}

		if meta["rotation"] != "" {
			secret.Rotate = meta["rotation"]
		}

		if meta["expiry"] != "" {
			expiryTime, err := time.Parse(time.RFC3339, meta["expiry"])
			if err != nil {
				f.Err = fmt.Errorf("invalid expiry format: %w", err)
				return
			}

			secret.Expire = expiryTime
		}

		return
	}

	f.Err = fmt.Errorf("secret with ID %q not found", id)
}

func (f *fakeCommandRunner) handleSecretRevoke(args []string) {}

func parseSecretMetadata(args []string) (map[string]string, []string) {
	meta := map[string]string{
		"label":       "",
		"owner":       "",
		"description": "",
		"rotation":    "",
		"expiry":      "",
	}

	remaining := make([]string, 0, len(args))

	for _, arg := range args {
		matched := false

		for key := range meta {
			prefix := "--" + key + "="
			if strings.HasPrefix(arg, prefix) {
				meta[key] = strings.TrimPrefix(arg, prefix)
				matched = true

				break
			}
		}

		if !matched {
			remaining = append(remaining, arg)
		}
	}

	return meta, remaining
}

func (f *fakeCommandRunner) handleStateGet(args []string) {
	if len(args) == 0 {
		f.Err = fmt.Errorf("state-get command requires at least one argument")
		return
	}

	key := args[0]

	if f.StoredState == nil {
		f.Output = []byte(`""`)
		f.Err = fmt.Errorf("stored state is nil")

		return
	}

	value, exists := f.StoredState[key]
	if !exists {
		f.Output = []byte(`""`)
		f.Err = fmt.Errorf("state key %s not found", key)

		return
	}

	output, err := json.Marshal(value)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal state value: %w", err)

		return
	}

	f.Output = output
}

func (f *fakeCommandRunner) handleStateSet(args []string) {
	if len(args) == 0 {
		f.Err = fmt.Errorf("state-set command requires at least one argument")
		return
	}

	if f.StoredState == nil {
		f.StoredState = make(StoredState)
	}

	for _, arg := range args {
		if parts := strings.SplitN(arg, "=", 2); len(parts) == 2 {
			key := parts[0]
			value := parts[1]

			if key == "" {
				f.Err = fmt.Errorf("state key cannot be empty")
				return
			}

			f.StoredState[key] = value
		} else {
			f.Err = fmt.Errorf("invalid state-set argument: %s", arg)
			return
		}
	}
}

func (f *fakeCommandRunner) handleStateDelete(args []string) {
	if len(args) == 0 {
		f.Err = fmt.Errorf("state-delete command requires at least one argument")
		return
	}

	key := args[0]

	if f.StoredState == nil {
		f.Err = fmt.Errorf("stored state is nil")
		return
	}

	if _, exists := f.StoredState[key]; !exists {
		f.Err = fmt.Errorf("state key %s not found", key)
		return
	}

	delete(f.StoredState, key)
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

func (f *fakeCommandRunner) handleActionSet(args []string) {
	if goops.ReadEnv().ActionName == "" {
		f.Err = fmt.Errorf("command action-set failed: ERROR not running an action")
		return
	}

	f.ActionResults = parseKeyValueArgs(args)
}

func (f *fakeCommandRunner) handleActionLog(_ []string) {
	if goops.ReadEnv().ActionName == "" {
		f.Err = fmt.Errorf("command action-log failed: ERROR not running an action")
		return
	}
}

func (f *fakeCommandRunner) handleActionFail(args []string) {
	if goops.ReadEnv().ActionName == "" {
		f.Err = fmt.Errorf("command action-fail failed: ERROR not running an action")
		return
	}

	f.ActionError = fmt.Errorf("%s", strings.Join(args, " "))
}

func (f *fakeCommandRunner) handleActionGet(_ []string) {
	env := goops.ReadEnv()
	if env.ActionName == "" {
		f.Err = fmt.Errorf("command action-get failed: ERROR not running an action")
		return
	}

	if f.ActionParameters == nil {
		f.ActionParameters = make(map[string]any)
	}

	output, err := json.Marshal(f.ActionParameters)
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal action parameters: %w", err)
		return
	}

	f.Output = output
}

func (f *fakeCommandRunner) handleApplicationVersionSet(args []string) {
	f.ApplicationVersion = args[0]
}
