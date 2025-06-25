package goopstest

import (
	"time"

	"github.com/canonical/pebble/client"
)

type Secret struct {
	ID          string
	Label       string
	Content     map[string]string
	Owner       string
	Description string
	Rotate      string
	Expire      time.Time
}

type DataBag map[string]string

type UnitID string

type Relation struct {
	Endpoint        string
	Interface       string
	ID              string
	RemoteAppName   string
	LocalAppData    DataBag
	LocalUnitData   DataBag
	RemoteAppData   DataBag
	RemoteUnitsData map[UnitID]DataBag
	RemoteModelUUID string
}

type PeerRelation struct {
	Endpoint      string
	Interface     string
	ID            string
	LocalAppData  DataBag
	LocalUnitData DataBag
	PeersData     map[UnitID]DataBag // Does not include data for the unit under test
}

type Port struct {
	Port     int
	Protocol string
}

type Model struct {
	Name string
	UUID string
}

type StoredState map[string]string

type Layer struct {
	Summary     string                `yaml:"summary"`
	Description string                `yaml:"description"`
	Services    map[string]Service    `yaml:"services"`
	LogTargets  map[string]*LogTarget `yaml:"log-targets"`
}

type Mount struct {
	Location string
	Source   string
}

type Exec struct {
	Command    []string
	ReturnCode int
	Stdout     string
	Stderr     string
}

type Container struct {
	Name            string
	CanConnect      bool
	Layers          map[string]Layer
	ServiceStatuses map[string]client.ServiceStatus
	Mounts          map[string]Mount
	Execs           []Exec
	Notices         []client.Notice
	CheckInfos      []client.CheckInfo
}

type StatusName string

const (
	StatusUnknown     StatusName = "unknown"
	StatusError       StatusName = "error"
	StatusActive      StatusName = "active"
	StatusBlocked     StatusName = "blocked"
	StatusMaintenance StatusName = "maintenance"
	StatusWaiting     StatusName = "waiting"
)

type Status struct {
	Name    StatusName `json:"status"`
	Message string     `json:"message"`
}

type State struct {
	Leader             bool
	UnitStatus         Status
	AppStatus          Status
	Config             map[string]any
	Secrets            []Secret
	ApplicationVersion string
	Relations          []Relation
	PeerRelations      []PeerRelation
	Ports              []Port
	Model              Model
	StoredState        StoredState
	Containers         []Container
}
