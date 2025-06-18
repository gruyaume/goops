package goopstest_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
	"gopkg.in/yaml.v3"
)

func ContainerCantConnectAddLayer() error {
	pebble := goops.Pebble("example")

	err := pebble.AddLayer(&client.AddLayerOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ContainerCantConnectPush() error {
	pebble := goops.Pebble("example")

	err := pebble.Push(&client.PushOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ContainerCantConnectPull() error {
	pebble := goops.Pebble("example")

	err := pebble.Pull(&client.PullOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ContainerCantConnectExec() error {
	pebble := goops.Pebble("example")

	_, err := pebble.Exec(&client.ExecOptions{
		Command: []string{"echo", "hello"},
	})
	if err != nil {
		return err
	}

	return nil
}

func TestContainerCantConnect(t *testing.T) {
	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "ContainerCantConnectAddLayer",
			fn:   ContainerCantConnectAddLayer,
		},
		{
			name: "ContainerCantConnectPush",
			fn:   ContainerCantConnectPush,
		},
		{
			name: "ContainerCantConnectPull",
			fn:   ContainerCantConnectPull,
		},
		{
			name: "ContainerCantConnectExec",
			fn:   ContainerCantConnectExec,
		},
	}
	for _, tt := range tests {
		ctx := goopstest.Context{
			Charm: tt.fn,
		}

		stateIn := &goopstest.State{
			Containers: []*goopstest.Container{
				{
					Name:       "example",
					CanConnect: false,
				},
			},
		}

		_, err := ctx.Run("install", stateIn)
		if err.Error() != "failed to run charm: cannot connect to Pebble" {
			t.Errorf("Run should have returned 'failed to run charm: cannot connect to Pebble', got: %v", err)
		}
	}
}

func ContainerCanConnect() error {
	pebble := goops.Pebble("example")

	_, err := pebble.SysInfo()
	if err != nil {
		return err
	}

	return nil
}

func TestContainerCanConnect(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerCanConnect,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
			},
		},
	}

	_, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
}

type ServiceConfig struct {
	Override string `yaml:"override"`
	Summary  string `yaml:"summary"`
	Command  string `yaml:"command"`
	Startup  string `yaml:"startup"`
}

type PebbleLayer struct {
	Summary     string                   `yaml:"summary"`
	Description string                   `yaml:"description"`
	Services    map[string]ServiceConfig `yaml:"services"`
	LogTargets  map[string]LogTarget     `yaml:"log-targets"`
}

type Check struct {
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

type LogTarget struct {
	Override string            `yaml:"override"`
	Type     string            `yaml:"type"`
	Location string            `yaml:"location"`
	Services []string          `yaml:"services"`
	Labels   map[string]string `yaml:"labels"`
}

type PebblePlan struct {
	Services   map[string]ServiceConfig `yaml:"services"`
	Checks     map[string]Check         `yaml:"checks"`
	LogTargets map[string]LogTarget     `yaml:"log-targets"`
}

func ContainerGetPebblePlan() error {
	pebble := goops.Pebble("example")

	dataBytes, err := pebble.PlanBytes(nil)
	if err != nil {
		return err
	}

	var plan PebblePlan

	err = yaml.Unmarshal(dataBytes, &plan)
	if err != nil {
		return err
	}

	service, exists := plan.Services["my-service"]
	if !exists {
		return fmt.Errorf("service 'my-service' not found in plan")
	}

	if service.Command != "/bin/my-service --config /etc/my-service/config.yaml" {
		return fmt.Errorf("unexpected command for 'my-service': %s", service.Command)
	}

	return nil
}

func TestContainerGetPebblePlan(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerGetPebblePlan,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]*goopstest.Layer{
					"my-service": {
						Summary:     "My service layer",
						Description: "This layer configures my service",
						Services: map[string]goopstest.Service{
							"my-service": {
								Startup:  "enabled",
								Override: "replace",
								Command:  "/bin/my-service --config /etc/my-service/config.yaml",
							},
						},
					},
				},
			},
		},
	}

	_, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
}

func TestContainerUnexistantGetPebblePlan(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerGetPebblePlan,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers:     map[string]*goopstest.Layer{},
			},
		},
	}

	_, err := ctx.Run("install", stateIn)
	if err.Error() != "failed to run charm: service 'my-service' not found in plan" {
		t.Fatalf("Run should have returned 'failed to run charm: service 'my-service' not found in plan', got: %v", err)
	}
}

func ContainerAddPebbleLayer() error {
	pebble := goops.Pebble("example")

	layerData, err := yaml.Marshal(PebbleLayer{
		LogTargets: map[string]LogTarget{
			"my-service/0": {
				Override: "replace",
				Services: []string{"all"},
				Type:     "loki",
				Location: "tcp://loki:3100",
				Labels: map[string]string{
					"juju-model":       "example-model",
					"juju-application": "example-app",
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("could not marshal layer data to YAML: %w", err)
	}

	err = pebble.AddLayer(&client.AddLayerOptions{
		Combine:   true,
		Label:     "example" + "-log-forwarding",
		LayerData: layerData,
	})
	if err != nil {
		return fmt.Errorf("could not add Pebble layer: %w", err)
	}

	return nil
}

func TestContainerAddPebbleLayer(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerAddPebbleLayer,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
			},
		},
	}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Containers) != 1 {
		t.Fatalf("Expected 1 container in stateOut, got %d", len(stateOut.Containers))
	}

	if len(stateOut.Containers[0].Layers) != 1 {
		t.Fatalf("Expected 1 Pebble layer in container, got %d", len(stateOut.Containers[0].Layers))
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"] == nil {
		t.Fatal("Expected Pebble layer 'example-log-forwarding' to be present, but it was not found")
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets == nil {
		t.Fatal("Expected Pebble layer 'example-log-forwarding' to have log targets, but none were found")
	}

	if len(stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets) != 1 {
		t.Fatalf("Expected 1 log target in Pebble layer 'example-log-forwarding', got %d", len(stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets))
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"] == nil {
		t.Fatal("Expected log target 'my-service/0' in Pebble layer 'example-log-forwarding', but it was not found")
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Type != "loki" {
		t.Errorf("Expected log target type 'loki', got '%s'", stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Type)
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Location != "tcp://loki:3100" {
		t.Errorf("Expected log target location 'tcp://loki:3100', got '%s'", stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Location)
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Labels["juju-model"] != "example-model" {
		t.Errorf("Expected log target label 'juju-model' to be 'example-model', got '%s'", stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Labels["juju-model"])
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Labels["juju-application"] != "example-app" {
		t.Errorf("Expected log target label 'juju-application' to be 'example-app', got '%s'", stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Labels["juju-application"])
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Override != "replace" {
		t.Errorf("Expected log target override to be 'replace', got '%s'", stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Override)
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Services == nil {
		t.Fatal("Expected log target 'my-service/0' to have services, but none were found")
	}

	if len(stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Services) != 1 {
		t.Fatalf("Expected 1 service in log target 'my-service/0', got %d", len(stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Services))
	}

	if stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Services[0] != "all" {
		t.Errorf("Expected service 'all' in log target 'my-service/0', got '%s'", stateOut.Containers[0].Layers["example-log-forwarding"].LogTargets["my-service/0"].Services[0])
	}

	if len(stateOut.Containers[0].Layers["example-log-forwarding"].Services) != 0 {
		t.Fatalf("Expected no services in Pebble layer 'example-log-forwarding', got %d", len(stateOut.Containers[0].Layers["example-log-forwarding"].Services))
	}
}

func ConatainerStartPebbleService() error {
	pebble := goops.Pebble("example")

	changeID, err := pebble.Start(&client.ServiceOptions{
		Names: []string{"my-service"},
	})
	if err != nil {
		return fmt.Errorf("could not start Pebble service: %w", err)
	}

	if changeID == "" {
		return fmt.Errorf("expected non-empty change ID after starting service")
	}

	return nil
}

func TestContainerStartPebbleService(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ConatainerStartPebbleService,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]*goopstest.Layer{
					"my-service": {
						Summary:     "My service layer",
						Description: "This layer configures my service",
						Services: map[string]goopstest.Service{
							"my-service": {
								Startup:  "enabled",
								Override: "replace",
								Command:  "/bin/my-service --config /etc/my-service/config.yaml",
							},
						},
					},
				},
				ServiceStatuses: map[string]client.ServiceStatus{
					"my-service": client.StatusInactive,
				},
			},
		},
	}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.Containers[0].ServiceStatuses["my-service"] != client.StatusActive {
		t.Errorf("Expected service 'my-service' to be active, got %s", stateOut.Containers[0].ServiceStatuses["my-service"])
	}
}

func ContainerGetPebbleServiceStatus() error {
	pebble := goops.Pebble("example")

	services, err := pebble.Services(&client.ServicesOptions{
		Names: []string{"my-service"},
	})
	if err != nil {
		return fmt.Errorf("could not get Pebble services: %w", err)
	}

	if len(services) != 1 {
		return fmt.Errorf("expected exactly one service, got %d", len(services))
	}

	if services[0].Name != "my-service" {
		return fmt.Errorf("expected service name 'my-service', got '%s'", services[0].Name)
	}

	if services[0].Current != client.StatusError {
		return fmt.Errorf("expected service status 'error', got '%s'", services[0].Current)
	}

	if services[0].Startup != client.StartupEnabled {
		return fmt.Errorf("expected service startup 'enabled', got '%s'", services[0].Startup)
	}

	return nil
}

func TestContainerGetPebbleServiceStatus(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerGetPebbleServiceStatus,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]*goopstest.Layer{
					"my-service": {
						Summary:     "My service layer",
						Description: "This layer configures my service",
						Services: map[string]goopstest.Service{
							"my-service": {
								Startup:  "enabled",
								Override: "replace",
								Command:  "/bin/my-service --config /etc/my-service/config.yaml",
							},
						},
					},
				},
				ServiceStatuses: map[string]client.ServiceStatus{
					"my-service": client.StatusError,
				},
			},
		},
	}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Containers) != 1 {
		t.Fatalf("Expected 1 container in stateOut, got %d", len(stateOut.Containers))
	}

	if len(stateOut.Containers[0].ServiceStatuses) != 1 {
		t.Fatalf("Expected 1 service status in container, got %d", len(stateOut.Containers[0].ServiceStatuses))
	}

	if stateOut.Containers[0].ServiceStatuses["my-service"] != client.StatusError {
		t.Errorf("Expected service 'my-service' to be error, got %s", stateOut.Containers[0].ServiceStatuses["my-service"])
	}
}

func ContainerStopPebbleService() error {
	pebble := goops.Pebble("example")

	changeID, err := pebble.Stop(&client.ServiceOptions{
		Names: []string{"my-service"},
	})
	if err != nil {
		return fmt.Errorf("could not stop Pebble service: %w", err)
	}

	if changeID == "" {
		return fmt.Errorf("expected non-empty change ID after stopping service")
	}

	return nil
}

func TestContainerStopPebbleService(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerStopPebbleService,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]*goopstest.Layer{
					"my-service": {
						Summary:     "My service layer",
						Description: "This layer configures my service",
						Services: map[string]goopstest.Service{
							"my-service": {
								Startup:  "enabled",
								Override: "replace",
								Command:  "/bin/my-service --config /etc/my-service/config.yaml",
							},
						},
					},
				},
				ServiceStatuses: map[string]client.ServiceStatus{
					"my-service": client.StatusActive,
				},
			},
		},
	}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.Containers[0].ServiceStatuses["my-service"] != client.StatusInactive {
		t.Errorf("Expected service 'my-service' to be inactive, got %s", stateOut.Containers[0].ServiceStatuses["my-service"])
	}
}

func ContainerRestartPebbleService() error {
	pebble := goops.Pebble("example")

	changeID, err := pebble.Restart(&client.ServiceOptions{
		Names: []string{"my-service"},
	})
	if err != nil {
		return fmt.Errorf("could not restart Pebble service: %w", err)
	}

	if changeID == "" {
		return fmt.Errorf("expected non-empty change ID after restarting service")
	}

	return nil
}

func TestContainerRestartPebbleService(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerRestartPebbleService,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]*goopstest.Layer{
					"my-service": {
						Summary:     "My service layer",
						Description: "This layer configures my service",
						Services: map[string]goopstest.Service{
							"my-service": {
								Startup:  "enabled",
								Override: "replace",
								Command:  "/bin/my-service --config /etc/my-service/config.yaml",
							},
						},
					},
				},
				ServiceStatuses: map[string]client.ServiceStatus{
					"my-service": client.StatusActive,
				},
			},
		},
	}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.Containers[0].ServiceStatuses["my-service"] != client.StatusActive {
		t.Errorf("Expected service 'my-service' to be active, got %s", stateOut.Containers[0].ServiceStatuses["my-service"])
	}
}

func ContainerReplanPebbleService() error {
	pebble := goops.Pebble("example")

	changeID, err := pebble.Replan(&client.ServiceOptions{
		Names: []string{"my-service"},
	})
	if err != nil {
		return fmt.Errorf("could not replan Pebble service: %w", err)
	}

	if changeID == "" {
		return fmt.Errorf("expected non-empty change ID after replanning service")
	}

	return nil
}

func ContainerPushFile() error {
	pebble := goops.Pebble("example")
	content := `# Example configuration file`
	path := "/etc/config.yaml"

	source := strings.NewReader(content)

	err := pebble.Push(&client.PushOptions{
		Source: source,
		Path:   path,
	})
	if err != nil {
		return fmt.Errorf("could not push file: %w", err)
	}

	return nil
}

func TestContainerPushFile(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerPushFile,
	}

	dname, err := os.MkdirTemp("", "sampledir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(dname)

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Mounts: map[string]goopstest.Mount{
					"config": {
						Location: "/etc/config.yaml",
						Source:   dname,
					},
				},
			},
		},
	}

	_, err = ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	content, err := os.ReadFile(dname + "/etc/config.yaml")
	if err != nil {
		t.Fatalf("Failed to read pushed file: %v", err)
	}

	expectedContent := "# Example configuration file"
	if string(content) != expectedContent {
		t.Errorf("Expected file content '%s', got '%s'", expectedContent, string(content))
	}
}
