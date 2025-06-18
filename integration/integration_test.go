package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/gruyaume/goops/integration/juju"
)

const (
	JujuModelName = "test-model"
)

func buildCharm() error {
	name := "charmcraft"
	args := []string{"pack", "--verbose", "--project-dir=../"}
	cmd := exec.Command(name, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func waitForActiveStatus(t *testing.T, client *juju.Client, timeout time.Duration) error {
	start := time.Now()

	for {
		if time.Since(start) > timeout {
			return fmt.Errorf("timed out waiting for active status")
		}

		status, err := client.Status()
		if err != nil {
			return err
		}

		if status.Applications["example"].ApplicationStatus.Current == "active" {
			return nil
		}

		t.Log("Waiting for active status, current status:", status)

		time.Sleep(1 * time.Second)
	}
}

func TestIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests, set environment variable INTEGRATION")
	}

	jujuClient := juju.New()

	jujuModels, err := jujuClient.ListModels()
	if err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}

	for _, model := range jujuModels {
		if model.ShortName == JujuModelName {
			t.Fatalf("Model %s already exists", JujuModelName)
		}
	}

	addModelOpts := &juju.AddModelOptions{
		Name: JujuModelName,
	}

	err = jujuClient.AddModel(addModelOpts)
	if err != nil {
		t.Fatalf("Failed to add model: %v", err)
	}

	t.Log("Model added successfully")

	err = buildCharm()
	if err != nil {
		t.Fatalf("Failed to build charm: %v", err)
	}

	t.Log("Charm built successfully")

	err = jujuClient.Deploy(&juju.DeployOptions{
		Charm: "./example_amd64.charm",
		Trust: true,
	})
	if err != nil {
		t.Fatalf("Failed to deploy charm: %v", err)
	}

	t.Log("Charm deployed successfully")

	err = waitForActiveStatus(t, jujuClient, 5*time.Minute)
	if err != nil {
		t.Fatalf("Failed to wait for active status: %v", err)
	}

	t.Log("Charm is active")
}
