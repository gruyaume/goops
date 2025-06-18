package goopstest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops"
	"gopkg.in/yaml.v3"
)

type FakePebbleClient struct {
	Containers    []*Container
	ContainerName string
}

func (f *FakePebbleClient) getContainer() *Container {
	var container *Container

	for _, c := range f.Containers {
		if c.Name == f.ContainerName {
			container = c
			break
		}
	}

	if container == nil {
		return nil
	}

	return container
}

func (f *FakePebbleClient) Exec(*client.ExecOptions) (goops.PebbleExecProcess, error) {
	container := f.getContainer()
	if container == nil {
		return nil, fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

	return nil, nil
}

func (f *FakePebbleClient) Pull(opts *client.PullOptions) error {
	container := f.getContainer()
	if container == nil {
		return fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return fmt.Errorf("cannot connect to Pebble")
	}

	if opts.Target == nil {
		return fmt.Errorf("target file cannot be nil")
	}

	for mountName, mount := range container.Mounts {
		if mount.Location != opts.Path {
			continue
		}

		safePath := filepath.Join(mount.Source, filepath.Clean(mount.Location))

		// Validate that safePath is within mount.Source
		rel, err := filepath.Rel(mount.Source, safePath)
		if err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
			return fmt.Errorf("refusing to read outside of mount source: %s", safePath)
		}

		sourceFile, err := os.Open(safePath) // #nosec G304 -- path validated above
		if err != nil {
			return fmt.Errorf("cannot open mount %s at %s: %w", mountName, safePath, err)
		}
		defer sourceFile.Close()

		if _, err := io.Copy(opts.Target, sourceFile); err != nil {
			return fmt.Errorf("failed to copy file contents from %s: %w", safePath, err)
		}
	}

	return nil
}

func (f *FakePebbleClient) Push(opts *client.PushOptions) error {
	container := f.getContainer()
	if container == nil {
		return fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return fmt.Errorf("cannot connect to Pebble")
	}

	for mountName, mount := range container.Mounts {
		if mount.Location != opts.Path {
			continue
		}

		if err := f.pushToMount(mountName, mount, opts); err != nil {
			return err
		}
	}

	return nil
}

func (f *FakePebbleClient) pushToMount(mountName string, mount Mount, opts *client.PushOptions) error {
	safePath := filepath.Join(mount.Source, filepath.Clean(mount.Location))

	rel, err := filepath.Rel(mount.Source, safePath)
	if err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return fmt.Errorf("refusing to write outside of mount source: %s", safePath)
	}

	if err := os.MkdirAll(filepath.Dir(safePath), 0o750); err != nil {
		return fmt.Errorf("cannot create directory for mount %s at %s: %w", mountName, filepath.Dir(safePath), err)
	}

	destFile, err := os.OpenFile(safePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600) // #nosec G304 -- validated path
	if err != nil {
		return fmt.Errorf("cannot open mount %s at %s: %w", mountName, safePath, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, opts.Source); err != nil {
		return fmt.Errorf("failed to copy file contents to %s: %w", safePath, err)
	}

	return nil
}

func (f *FakePebbleClient) Restart(*client.ServiceOptions) (string, error) {
	container := f.getContainer()
	if container == nil {
		return "", fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return "", fmt.Errorf("cannot connect to Pebble")
	}

	return "123", nil
}

func (f *FakePebbleClient) Replan(*client.ServiceOptions) (string, error) {
	container := f.getContainer()
	if container == nil {
		return "", fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return "", fmt.Errorf("cannot connect to Pebble")
	}

	return "123", nil
}

func (f *FakePebbleClient) Start(opts *client.ServiceOptions) (string, error) {
	container := f.getContainer()
	if container == nil {
		return "", fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return "", fmt.Errorf("cannot connect to Pebble")
	}

	for _, name := range opts.Names {
		if container.ServiceStatuses == nil {
			container.ServiceStatuses = make(map[string]client.ServiceStatus)
		}

		container.ServiceStatuses[name] = client.StatusActive
	}

	return "123", nil
}

func (f *FakePebbleClient) Stop(opts *client.ServiceOptions) (string, error) {
	container := f.getContainer()
	if container == nil {
		return "", fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return "", fmt.Errorf("cannot connect to Pebble")
	}

	for _, name := range opts.Names {
		if container.ServiceStatuses == nil {
			container.ServiceStatuses = make(map[string]client.ServiceStatus)
		}

		container.ServiceStatuses[name] = client.StatusInactive
	}

	return "123", nil
}

func (f *FakePebbleClient) Services(opts *client.ServicesOptions) ([]*client.ServiceInfo, error) {
	container := f.getContainer()
	if container == nil {
		return nil, fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

	var services []*client.ServiceInfo

	for _, layer := range container.Layers {
		for name, service := range layer.Services {
			if opts.Names != nil && !contains(opts.Names, name) {
				continue
			}

			services = append(services, &client.ServiceInfo{
				Name:    name,
				Startup: client.ServiceStartup(service.Startup),
				Current: container.ServiceStatuses[name],
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
	container := f.getContainer()
	if container == nil {
		return nil, fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

	return nil, nil
}

func (f *FakePebbleClient) WaitChange(string, *client.WaitChangeOptions) (*client.Change, error) {
	container := f.getContainer()
	if container == nil {
		return nil, fmt.Errorf("container not found")
	}

	if !container.CanConnect {
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
	container := f.getContainer()
	if container == nil {
		return nil, fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

	plan := pebblePlan{
		Services:   make(map[string]serviceConfig),
		Checks:     make(map[string]check),
		LogTargets: make(map[string]logTarget),
	}

	for _, layer := range container.Layers {
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
	container := f.getContainer()
	if container == nil {
		return fmt.Errorf("container not found")
	}

	if !container.CanConnect {
		return fmt.Errorf("cannot connect to Pebble")
	}

	var layer Layer
	if err := yaml.Unmarshal(opts.LayerData, &layer); err != nil {
		return fmt.Errorf("cannot unmarshal layer data: %w", err)
	}

	if container.Layers == nil {
		container.Layers = make(map[string]*Layer)
	}

	container.Layers[opts.Label] = &layer

	return nil
}

type fakePebbleGetter struct {
	Containers []*Container
}

func (f *fakePebbleGetter) Pebble(name string) goops.PebbleClient {
	return &FakePebbleClient{
		Containers:    f.Containers,
		ContainerName: name,
	}
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
