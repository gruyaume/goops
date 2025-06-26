package goopstest_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
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
		t.Run(tt.name, func(t *testing.T) {
			ctx := goopstest.NewContext(tt.fn)

			stateIn := goopstest.State{
				Containers: []goopstest.Container{
					{
						Name:       "example",
						CanConnect: false,
					},
				},
			}

			_ = ctx.Run("install", stateIn)

			if ctx.CharmErr.Error() != "cannot connect to Pebble" {
				t.Errorf("Run should have returned 'cannot connect to Pebble', got: %v", ctx.CharmErr)
			}
		})
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
	ctx := goopstest.NewContext(ContainerCanConnect)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
			},
		},
	}

	_ = ctx.Run("install", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
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
	ctx := goopstest.NewContext(ContainerGetPebblePlan)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]goopstest.Layer{
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

	_ = ctx.Run("install", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
	}
}

func TestContainerUnexistantGetPebblePlan(t *testing.T) {
	ctx := goopstest.NewContext(ContainerGetPebblePlan)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers:     map[string]goopstest.Layer{},
			},
		},
	}

	_ = ctx.Run("install", stateIn)

	if ctx.CharmErr.Error() != "service 'my-service' not found in plan" {
		t.Fatalf("Run should have returned 'service 'my-service' not found in plan', got: %v", ctx.CharmErr)
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
	ctx := goopstest.NewContext(ContainerAddPebbleLayer)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
			},
		},
	}

	stateOut := ctx.Run("install", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
	}

	if len(stateOut.Containers) != 1 {
		t.Fatalf("Expected 1 container in stateOut, got %d", len(stateOut.Containers))
	}

	layer := stateOut.Containers[0].Layers["example-log-forwarding"]

	expectedLogTarget := &goopstest.LogTarget{
		Type:     "loki",
		Location: "tcp://loki:3100",
		Labels: map[string]string{
			"juju-model":       "example-model",
			"juju-application": "example-app",
		},
		Override: "replace",
		Services: []string{"all"},
	}

	actualLogTarget, ok := layer.LogTargets["my-service/0"]
	if !ok {
		t.Fatal("Expected log target 'my-service/0' to be present, but it was not found")
	}

	if !reflect.DeepEqual(actualLogTarget, expectedLogTarget) {
		t.Errorf("Log target 'my-service/0' does not match expected configuration.\nExpected: %+v\nActual: %+v", expectedLogTarget, actualLogTarget)
	}

	if len(layer.Services) != 0 {
		t.Fatalf("Expected no services in Pebble layer 'example-log-forwarding', got %d", len(layer.Services))
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
	ctx := goopstest.NewContext(ConatainerStartPebbleService)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]goopstest.Layer{
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

	stateOut := ctx.Run("install", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
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
	ctx := goopstest.NewContext(ContainerGetPebbleServiceStatus)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]goopstest.Layer{
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

	stateOut := ctx.Run("install", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
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
	ctx := goopstest.NewContext(ContainerStopPebbleService)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]goopstest.Layer{
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

	stateOut := ctx.Run("install", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
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
	ctx := goopstest.NewContext(ContainerRestartPebbleService)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Layers: map[string]goopstest.Layer{
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

	stateOut := ctx.Run("install", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
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

	err := pebble.Push(&client.PushOptions{
		Source: strings.NewReader(`# Example configuration file`),
		Path:   "/etc/config.yaml",
	})
	if err != nil {
		return fmt.Errorf("could not push file: %w", err)
	}

	return nil
}

func TestContainerPushFile(t *testing.T) {
	ctx := goopstest.NewContext(ContainerPushFile)

	dname, err := os.MkdirTemp("", "sampledir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(dname)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Mounts: map[string]goopstest.Mount{
					"config": {
						Location: "/etc",
						Source:   dname,
					},
				},
			},
		},
	}

	_ = ctx.Run("install", stateIn)

	content, err := os.ReadFile(dname + "/etc/config.yaml")
	if err != nil {
		t.Fatalf("Failed to read pushed file: %v", err)
	}

	expectedContent := "# Example configuration file"
	if string(content) != expectedContent {
		t.Errorf("Expected file content '%s', got '%s'", expectedContent, string(content))
	}
}

func ContainerPullFile() error {
	pebble := goops.Pebble("example")

	target := &bytes.Buffer{}

	err := pebble.Pull(&client.PullOptions{
		Path:   "/etc/config.yaml",
		Target: target,
	})
	if err != nil {
		return fmt.Errorf("could not push file: %w", err)
	}

	if target.String() != "# Example configuration file" {
		return fmt.Errorf("expected file content '# Example configuration file', got '%s'", target.String())
	}

	return nil
}

func TestContainerPullFile(t *testing.T) {
	ctx := goopstest.NewContext(ContainerPullFile)

	dname, err := os.MkdirTemp("", "sampledir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(dname)

	tempLocation := dname + "/etc/config.yaml"

	err = os.MkdirAll(filepath.Dir(tempLocation), 0o755)
	if err != nil {
		t.Fatalf("cannot create directory for mount %s at %s: %v", "config", filepath.Dir(tempLocation), err)
	}

	destFile, err := os.OpenFile(tempLocation, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		t.Fatalf("cannot open mount %s at %s: %v", "config", tempLocation, err)
	}

	defer destFile.Close()

	source := strings.NewReader(`# Example configuration file`)

	if _, err := io.Copy(destFile, source); err != nil {
		t.Fatalf("failed to copy file contents to %s: %v", tempLocation, err)
	}

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
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

	_ = ctx.Run("install", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
	}
}

func MultiContainer() error {
	pebble1 := goops.Pebble("example1")
	pebble2 := goops.Pebble("example2")

	_, err := pebble1.SysInfo()
	if err != nil {
		return fmt.Errorf("could not connect to example1 Pebble: %w", err)
	}

	_, err = pebble2.SysInfo()
	if err != nil {
		return fmt.Errorf("could not connect to example2 Pebble: %w", err)
	}

	return nil
}

func TestMultiContainer(t *testing.T) {
	tests := []struct {
		name               string
		example1CanConnect bool
		example2CanConnect bool
		expectError        bool
		expectedError      string
	}{
		{
			name:               "CanConnectBothContainers",
			example1CanConnect: true,
			example2CanConnect: true,
			expectError:        false,
			expectedError:      "",
		},
		{
			name:               "CantConnectExample1",
			example1CanConnect: false,
			example2CanConnect: true,
			expectError:        true,
			expectedError:      "could not connect to example1 Pebble: cannot connect to Pebble",
		},
		{
			name:               "CantConnectExample2",
			example1CanConnect: true,
			example2CanConnect: false,
			expectError:        true,
			expectedError:      "could not connect to example2 Pebble: cannot connect to Pebble",
		},
		{
			name:               "CantConnectBothContainers",
			example1CanConnect: false,
			example2CanConnect: false,
			expectError:        true,
			expectedError:      "could not connect to example1 Pebble: cannot connect to Pebble",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := goopstest.NewContext(MultiContainer)

			stateIn := goopstest.State{
				Containers: []goopstest.Container{
					{
						Name:       "example1",
						CanConnect: tt.example1CanConnect,
					},
					{
						Name:       "example2",
						CanConnect: tt.example2CanConnect,
					},
				},
			}

			_ = ctx.Run("install", stateIn)

			if tt.expectError {
				if ctx.CharmErr == nil {
					t.Errorf("Charm should have returned an error, but got none")
				} else if ctx.CharmErr.Error() != tt.expectedError {
					t.Errorf("Charm returned error %v, expected %v", ctx.CharmErr.Error(), tt.expectedError)
				}
			}
		})
	}
}
