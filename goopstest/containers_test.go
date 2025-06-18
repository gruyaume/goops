package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
	"gopkg.in/yaml.v3"
)

func ContainerCanConnect() error {
	pebble := goops.Pebble("example")

	_, err := pebble.SysInfo()
	if err != nil {
		return err
	}

	return nil
}

func TestContainerCantConnect(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerCanConnect,
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
	if err == nil {
		t.Fatalf("Run should have returned an error, but got nil")
	}
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
