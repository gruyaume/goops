package goopstest

import (
	"fmt"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops"
	"gopkg.in/yaml.v3"
)

type FakePebbleClient struct {
	CanConnect bool
	Layers     map[string]*Layer
}

func (f *FakePebbleClient) Exec(*client.ExecOptions) (goops.PebbleExecProcess, error) {
	return nil, nil
}

func (f *FakePebbleClient) Pull(*client.PullOptions) error {
	return nil
}

func (f *FakePebbleClient) Push(*client.PushOptions) error {
	return nil
}

func (f *FakePebbleClient) Restart(*client.ServiceOptions) (string, error) {
	return "", nil
}

func (f *FakePebbleClient) Start(*client.ServiceOptions) (string, error) {
	return "", nil
}

func (f *FakePebbleClient) Stop(*client.ServiceOptions) (string, error) {
	return "", nil
}

func (f *FakePebbleClient) SysInfo() (*client.SysInfo, error) {
	if !f.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

	return nil, nil
}

func (f *FakePebbleClient) WaitChange(string, *client.WaitChangeOptions) (*client.Change, error) {
	return nil, nil
}

type serviceConfig struct {
	Override string `yaml:"override"`
	Summary  string `yaml:"summary"`
	Command  string `yaml:"command"`
	Startup  string `yaml:"startup"`
}

type check struct {
	Override  string `yaml:"override"`
	Level     string `yaml:"level"`
	Startup   string `yaml:"startup"`
	Period    string `yaml:"period"`
	Timeout   string `yaml:"timeout"`
	Threshold string `yaml:"threshold"`
	HTTP      string `yaml:"http"`
	TCP       string `yaml:"tcp"`
	Exec      string `yaml:"exec"`
}

type logTarget struct {
	Override string            `yaml:"override"`
	Type     string            `yaml:"type"`
	Location string            `yaml:"location"`
	Services []string          `yaml:"services"`
	Labels   map[string]string `yaml:"labels"`
}

type pebblePlan struct {
	Services   map[string]serviceConfig `yaml:"services"`
	Checks     map[string]check         `yaml:"checks"`
	LogTargets map[string]logTarget     `yaml:"log-targets"`
}

func (f *FakePebbleClient) PlanBytes(_ *client.PlanOptions) ([]byte, error) {
	plan := pebblePlan{
		Services:   make(map[string]serviceConfig),
		Checks:     make(map[string]check),
		LogTargets: make(map[string]logTarget),
	}

	for _, layer := range f.Layers {
		for serviceName, service := range layer.Services {
			plan.Services[serviceName] = serviceConfig(service)
		}
	}

	dataBytes, err := yaml.Marshal(&plan)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal Pebble plan: %w", err)
	}

	return dataBytes, nil
}

func (f *FakePebbleClient) AddLayer(opts *client.AddLayerOptions) error {
	var layer Layer
	if err := yaml.Unmarshal(opts.LayerData, &layer); err != nil {
		return fmt.Errorf("cannot unmarshal layer data: %w", err)
	}

	if f.Layers == nil {
		f.Layers = make(map[string]*Layer)
	}

	f.Layers[opts.Label] = &layer

	return nil
}

func (f *FakePebbleClient) Pebble(string) goops.PebbleClient {
	return f
}

type Service struct {
	Override string
	Summary  string
	Command  string
	Startup  string
}

type LogTarget struct {
	Override string            `yaml:"override"`
	Type     string            `yaml:"type"`
	Location string            `yaml:"location"`
	Services []string          `yaml:"services"`
	Labels   map[string]string `yaml:"labels"`
}

type Layer struct {
	Summary     string                `yaml:"summary"`
	Description string                `yaml:"description"`
	Services    map[string]Service    `yaml:"services"`
	LogTargets  map[string]*LogTarget `yaml:"log-targets"`
}

type Container struct {
	Name            string
	CanConnect      bool
	Layers          map[string]*Layer
	ServiceStatuses map[string]client.ServiceStatus
	Mounts          map[string]Mount
	Execs           []Exec
	Notices         []client.Notice
	CheckInfos      []client.CheckInfo
}
