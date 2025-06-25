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
	UnitStatus         Status
	AppStatus          Status
	Leader             bool
	Config             map[string]any
	Secrets            []*Secret
	ActionResults      map[string]string
	ActionParameters   map[string]any
	ActionError        error
	ApplicationVersion string
	Relations          []*Relation
	PeerRelations      []*PeerRelation
	Ports              []*Port
	StoredState        StoredState
	AppName            string
	UnitID             string
	JujuLog            []JujuLogLine
	Model              *Model
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
		"relation-model-get":      f.handleRelationModelGet,
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
		"status-get":              f.handleStatusGet,
		"status-set":              f.handleStatusSet,
		"juju-log":                f.handleJujuLog,
	}

	if handler, exists := handlers[name]; exists {
		handler(args)
		return f.Output, f.Err
	}

	return nil, fmt.Errorf("unknown command: %s", name)
}

type AppStatus struct {
	Name    StatusName `json:"status"`
	Message string     `json:"message"`
}

type appStatusReturn struct {
	AppStatus AppStatus `json:"application-status"`
}

func (f *fakeCommandRunner) handleStatusGet(args []string) {
	if args[0] == "--application" {
		if !f.Leader {
			f.Err = fmt.Errorf("command status-get failed: ERROR finding application status: this unit is not the leader")
			return
		}

		appStatus := appStatusReturn{
			AppStatus: AppStatus{
				Name:    f.AppStatus.Name,
				Message: f.AppStatus.Message,
			},
		}

		f.Output, f.Err = json.Marshal(appStatus)
	} else {
		unitStatus := f.UnitStatus

		f.Output, f.Err = json.Marshal(unitStatus)
	}
}

func (f *fakeCommandRunner) handleStatusSet(args []string) {
	if args[0] == "--application" {
		if !f.Leader {
			f.Err = fmt.Errorf("command status-set failed: ERROR setting application status: this unit is not the leader")
			return
		}

		if len(args) < 2 {
			f.Err = fmt.Errorf("status-set command requires an application status after --application")
			return
		}

		f.AppStatus = Status{
			Name:    StatusName(args[1]),
			Message: strings.Join(args[2:], " "),
		}
	} else {
		f.UnitStatus = Status{
			Name:    StatusName(args[0]),
			Message: strings.Join(args[1:], " "),
		}
	}
}

func (f *fakeCommandRunner) handleJujuLog(args []string) {
	var logLevel LogLevel

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

	for _, p := range f.Ports {
		if p.Port == port && p.Protocol == protocol {
			return
		}
	}

	f.Ports = append(f.Ports, &Port{
		Port:     port,
		Protocol: protocol,
	})
}

func (f *fakeCommandRunner) handleClosePort(args []string) {
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

	if len(f.Relations) == 0 && len(f.PeerRelations) == 0 {
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

	for _, relation := range f.PeerRelations {
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

func (f *fakeCommandRunner) findPeerRelationByID(id string) *PeerRelation {
	for i := range f.PeerRelations {
		if f.PeerRelations[i].ID == id {
			return f.PeerRelations[i]
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

func getAppNameFromUnitID(unitID string) string {
	if strings.Contains(unitID, "/") {
		parts := strings.Split(unitID, "/")
		if len(parts) == 2 {
			return parts[0]
		}
	}

	return ""
}

func (f *fakeCommandRunner) handleRelationGet(args []string) {
	isApp, relationID, unitID, err := parseRelationGetArgs(args)
	if err != nil {
		f.Err = err
		return
	}

	relation := f.findRelationByID(relationID)
	peerRelation := f.findPeerRelationByID(relationID)

	if relation == nil && peerRelation == nil {
		f.Err = fmt.Errorf("command relation-get failed: ERROR invalid value %q for option -r: relation not found", relationID)
		return
	}

	if relation != nil {
		argAppName := getAppNameFromUnitID(unitID)
		ctxAppName := getAppNameFromUnitID(f.UnitID)

		isLocal := argAppName == ctxAppName

		if !isLocal && argAppName != relation.RemoteAppName {
			f.Err = fmt.Errorf("command relation-get failed: ERROR permission denied")
			return
		}

		if isApp && isLocal && !f.Leader {
			f.Err = fmt.Errorf("command relation-get failed: ERROR permission denied")
			return
		}

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

	if peerRelation != nil {
		data, err := f.selectPeerRelationData(peerRelation, isApp, unitID)
		if err != nil {
			f.Err = err
			return
		}

		f.Output, err = json.Marshal(data)
		if err != nil {
			f.Err = fmt.Errorf("failed to marshal relation data: %w", err)
		}
	}
}

func (f *fakeCommandRunner) selectRelationData(rel *Relation, isApp bool, isLocal bool, unitID string) (any, error) {
	if isApp {
		if isLocal {
			return safeCopy(rel.LocalAppData), nil
		}

		return safeCopy(rel.RemoteAppData), nil
	}

	if isLocal {
		if f.UnitID != unitID {
			return nil, nil
		}

		if rel.LocalUnitData == nil {
			return nil, fmt.Errorf("local unit data not found for relation %s", rel.ID)
		}

		return safeCopy(rel.LocalUnitData), nil
	}

	unitData, ok := rel.RemoteUnitsData[UnitID(unitID)]
	if !ok {
		return nil, fmt.Errorf("command relation-get failed: ERROR cannot read settings for unit %q in relation %q: unit %q: settings not found", unitID, rel.ID, unitID)
	}

	return unitData, nil
}

func (f *fakeCommandRunner) selectPeerRelationData(rel *PeerRelation, isApp bool, unitID string) (any, error) {
	if isApp {
		return safeCopy(rel.LocalAppData), nil
	}

	if f.UnitID == unitID {
		return safeCopy(rel.LocalUnitData), nil
	}

	unitData, ok := rel.PeersData[UnitID(unitID)]
	if !ok {
		return nil, fmt.Errorf("command relation-get failed: ERROR cannot read settings for unit %q in peer relation %q: unit %q: settings not found", unitID, rel.ID, unitID)
	}

	return unitData, nil
}

func (f *fakeCommandRunner) handleRelationList(args []string) {
	meta, _ := splitPrefixedArgs(args, "-")
	relationID := meta["r"]

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

	for _, peerRelation := range f.PeerRelations {
		if peerRelation.ID == relationID {
			unitIDs := make([]string, 0, len(peerRelation.PeersData))
			for unitID := range peerRelation.PeersData {
				unitIDs = append(unitIDs, string(unitID))
			}

			output, err := json.Marshal(unitIDs)
			if err != nil {
				f.Err = fmt.Errorf("failed to marshal peer relation units: %w", err)
				return
			}

			f.Output = output

			return
		}
	}

	f.Err = fmt.Errorf("command relation-list failed: ERROR invalid value %q for option -r: relation not found", relationID)
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

	relation := f.findRelationByID(relationID)
	peerRelation := f.findPeerRelationByID(relationID)

	if relation == nil && peerRelation == nil {
		f.Err = fmt.Errorf("command relation-set failed: ERROR invalid value %q for option -r: relation not found", relationID)
		return
	}

	if isApp && !f.Leader {
		f.Err = fmt.Errorf("command relation-set failed: ERROR cannot write relation settings")
		return
	}

	updateDataBag := func(target *DataBag) {
		if *target == nil {
			*target = make(DataBag)
		}

		for k, v := range data {
			(*target)[k] = v
		}
	}

	if relation != nil {
		target := &relation.LocalUnitData
		if isApp {
			target = &relation.LocalAppData
		}

		updateDataBag(target)
	}

	if peerRelation != nil {
		target := &peerRelation.LocalUnitData
		if isApp {
			target = &peerRelation.LocalAppData
		}

		updateDataBag(target)
	}
}

func (f *fakeCommandRunner) handleRelationModelGet(args []string) {
	meta, _ := splitPrefixedArgs(args, "-")
	relationID := meta["r"]

	if relationID == "" {
		f.Err = fmt.Errorf("command relation-model-get failed: ERROR no relation ID specified with -r")
		return
	}

	relation := f.findRelationByID(relationID)
	peerRelation := f.findPeerRelationByID(relationID)

	if relation == nil && peerRelation == nil {
		f.Err = fmt.Errorf("command relation-model-get failed: ERROR invalid value %q for option -r: relation not found", relationID)
		return
	}

	var uuid string
	if relation != nil {
		uuid = relation.RemoteModelUUID
	}

	if uuid == "" && f.Model != nil {
		uuid = f.Model.UUID
	}

	outputBytes, err := json.Marshal(map[string]string{"uuid": uuid})
	if err != nil {
		f.Err = fmt.Errorf("failed to marshal relation model get output: %w", err)
		return
	}

	f.Output = outputBytes
}

func (f *fakeCommandRunner) handleSecretAdd(args []string) {
	meta, remaining := splitPrefixedArgs(args, "--")
	label := meta["label"]
	owner := meta["owner"]
	description := meta["description"]
	rotation := meta["rotate"]
	expiry := meta["expire"]

	if !f.Leader && owner != "unit" {
		f.Err = fmt.Errorf("command secret-add failed: ERROR this unit is not the leader")
		return
	}

	expiryTime, err := parseRFC3339(expiry)
	if err != nil {
		f.Err = fmt.Errorf("invalid expiry format: %w", err)
		return
	}

	content := parseKeyValueArgs(remaining)

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
	meta, remaining := splitPrefixedArgs(args, "--")
	label := meta["label"]

	if label != "" {
		secret := findSecretByLabel(f.Secrets, label)
		if secret == nil {
			f.Err = fmt.Errorf("secret with label %q not found", label)
			return
		}

		f.setSecretOutput(secret)

		return
	}

	var id string

	for _, arg := range remaining {
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

func splitPrefixedArgs(args []string, prefix string) (map[string]string, []string) {
	meta := make(map[string]string)
	remaining := make([]string, 0, len(args))

	for _, arg := range args {
		if strings.HasPrefix(arg, prefix) {
			parts := strings.SplitN(strings.TrimPrefix(arg, prefix), "=", 2)
			if len(parts) == 2 {
				meta[parts[0]] = parts[1]
				continue
			}
		}

		remaining = append(remaining, arg)
	}

	return meta, remaining
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
	meta, remaining := splitPrefixedArgs(args, "--")
	label := meta["label"]

	var id string

	for _, arg := range remaining {
		if !strings.HasPrefix(arg, "--") {
			id = arg
			break
		}
	}

	var secret *Secret

	switch {
	case label != "":
		secret = findSecretByLabel(f.Secrets, label)
		if secret == nil || (!f.Leader && secret.Owner != "unit") {
			f.Err = fmt.Errorf(`ERROR secret %q not found`, label)
			return
		}
	case id != "":
		secret = findSecretByID(f.Secrets, id)
		if secret == nil || (!f.Leader && secret.Owner != "unit") {
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
		f.Output = json.RawMessage("null")
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
	secretID := args[0]

	if !f.Leader {
		f.Err = fmt.Errorf(`ERROR secret "%s" not found`, secretID)
		return
	}
}

func (f *fakeCommandRunner) handleSecretSet(args []string) {
	if !f.Leader {
		f.Output = []byte(`null`)
		return
	}

	id := args[0]
	meta, remaining := splitPrefixedArgs(args[1:], "--")

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

		expiryTime, err := parseRFC3339(meta["expiry"])
		if err != nil {
			f.Err = fmt.Errorf("invalid expiry format: %w", err)
			return
		}

		secret.Expire = expiryTime

		return
	}

	f.Err = fmt.Errorf("secret with ID %q not found", id)
}

func parseRFC3339(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339, s)
}

func (f *fakeCommandRunner) handleSecretRevoke(args []string) {}

func (f *fakeCommandRunner) handleStateGet(args []string) {
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
