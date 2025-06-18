package goopstest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops"
	"gopkg.in/yaml.v3"
)

type FakePebbleClient struct {
	CanConnect      bool
	Layers          map[string]*Layer
	ServiceStatuses map[string]client.ServiceStatus
	Mounts          map[string]Mount
}

func (f *FakePebbleClient) Exec(*client.ExecOptions) (goops.PebbleExecProcess, error) {
	if !f.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

	return nil, nil
}

func (f *FakePebbleClient) Pull(opts *client.PullOptions) error {
	if !f.CanConnect {
		return fmt.Errorf("cannot connect to Pebble")
	}

	if opts.Target == nil {
		return fmt.Errorf("target file cannot be nil")
	}

	for mountName, mount := range f.Mounts {
		if mount.Location == opts.Path {
			tempLocation := mount.Source + mount.Location

			sourceFile, err := os.Open(tempLocation)
			if err != nil {
				return fmt.Errorf("cannot open mount %s at %s: %w", mountName, tempLocation, err)
			}

			defer sourceFile.Close()

			_, err = io.Copy(opts.Target, sourceFile)
			if err != nil {
				return fmt.Errorf("failed to copy file contents from %s: %w", tempLocation, err)
			}
		}
	}

	return nil
}

func (f *FakePebbleClient) Push(opts *client.PushOptions) error {
	if !f.CanConnect {
		return fmt.Errorf("cannot connect to Pebble")
	}

	for mountName, mount := range f.Mounts {
		if mount.Location == opts.Path {
			tempLocation := mount.Source + mount.Location

			err := os.MkdirAll(filepath.Dir(tempLocation), 0o750)
			if err != nil {
				return fmt.Errorf("cannot create directory for mount %s at %s: %w", mountName, filepath.Dir(tempLocation), err)
			}

			destFile, err := os.OpenFile(tempLocation, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
			if err != nil {
				return fmt.Errorf("cannot open mount %s at %s: %w", mountName, tempLocation, err)
			}

			defer destFile.Close()

			if _, err := io.Copy(destFile, opts.Source); err != nil {
				return fmt.Errorf("failed to copy file contents to %s: %w", tempLocation, err)
			}
		}
	}

	return nil
}

func (f *FakePebbleClient) Restart(*client.ServiceOptions) (string, error) {
	if !f.CanConnect {
		return "", fmt.Errorf("cannot connect to Pebble")
	}

	return "123", nil
}

func (f *FakePebbleClient) Replan(*client.ServiceOptions) (string, error) {
	if !f.CanConnect {
		return "", fmt.Errorf("cannot connect to Pebble")
	}

	return "123", nil
}

func (f *FakePebbleClient) Start(opts *client.ServiceOptions) (string, error) {
	if !f.CanConnect {
		return "", fmt.Errorf("cannot connect to Pebble")
	}

	for _, name := range opts.Names {
		if f.ServiceStatuses == nil {
			f.ServiceStatuses = make(map[string]client.ServiceStatus)
		}

		f.ServiceStatuses[name] = client.StatusActive
	}

	return "123", nil
}

func (f *FakePebbleClient) Stop(opts *client.ServiceOptions) (string, error) {
	if !f.CanConnect {
		return "", fmt.Errorf("cannot connect to Pebble")
	}

	for _, name := range opts.Names {
		if f.ServiceStatuses == nil {
			f.ServiceStatuses = make(map[string]client.ServiceStatus)
		}

		f.ServiceStatuses[name] = client.StatusInactive
	}

	return "123", nil
}

func (f *FakePebbleClient) Services(opts *client.ServicesOptions) ([]*client.ServiceInfo, error) {
	if !f.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

	var services []*client.ServiceInfo

	for _, layer := range f.Layers {
		for name, service := range layer.Services {
			if opts.Names != nil && !contains(opts.Names, name) {
				continue
			}

			services = append(services, &client.ServiceInfo{
				Name:    name,
				Startup: client.ServiceStartup(service.Startup),
				Current: f.ServiceStatuses[name],
			})
		}
	}

	return services, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

func (f *FakePebbleClient) SysInfo() (*client.SysInfo, error) {
	if !f.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

	return nil, nil
}

func (f *FakePebbleClient) WaitChange(string, *client.WaitChangeOptions) (*client.Change, error) {
	if !f.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

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
	if !f.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

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
	if !f.CanConnect {
		return fmt.Errorf("cannot connect to Pebble")
	}

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
