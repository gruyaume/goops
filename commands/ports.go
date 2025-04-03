package commands

import (
	"encoding/json"
	"fmt"
)

const (
	openPortCommand    = "open-port"
	closePortCommand   = "close-port"
	openedPortsCommand = "opened-ports"
)

type Port struct {
	Port     int
	Protocol string // allowed values: tcp, udp
}

func (command Command) SetPorts(ports []Port) error {
	openedPorts, err := command.OpenedPorts()
	if err != nil {
		return fmt.Errorf("failed to get opened ports: %w", err)
	}

	desiredMap := make(map[string]Port)

	for _, port := range ports {
		key := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		desiredMap[key] = port
	}

	openedMap := make(map[string]Port)

	for _, port := range openedPorts {
		key := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		openedMap[key] = port
	}

	// Open ports that are desired but not currently opened.
	for key, port := range desiredMap {
		if _, exists := openedMap[key]; !exists {
			if err := command.OpenPort(port); err != nil {
				return fmt.Errorf("failed to open port %s: %w", key, err)
			}
		}
	}

	// Close ports that are currently opened but not desired.
	for key, port := range openedMap {
		if _, exists := desiredMap[key]; !exists {
			if err := command.ClosePort(port); err != nil {
				return fmt.Errorf("failed to close port %s: %w", key, err)
			}
		}
	}

	return nil
}

func (command Command) OpenPort(port Port) error {
	if port.Port < 0 || port.Port > 65535 {
		return fmt.Errorf("port %d is out of range", port.Port)
	}

	if port.Protocol != "tcp" && port.Protocol != "udp" {
		return fmt.Errorf("protocol %s is not supported", port.Protocol)
	}

	args := []string{fmt.Sprintf("%d/%s", port.Port, port.Protocol)}

	_, err := command.Runner.Run(openPortCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to open port %d/%s: %w", port.Port, port.Protocol, err)
	}

	return nil
}

func (command Command) ClosePort(port Port) error {
	if port.Port < 0 || port.Port > 65535 {
		return fmt.Errorf("port %d is out of range", port.Port)
	}

	if port.Protocol != "tcp" && port.Protocol != "udp" {
		return fmt.Errorf("protocol %s is not supported", port.Protocol)
	}

	args := []string{fmt.Sprintf("%d/%s", port.Port, port.Protocol)}

	_, err := command.Runner.Run(closePortCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to open port %d/%s: %w", port.Port, port.Protocol, err)
	}

	return nil
}

func (command Command) OpenedPorts() ([]Port, error) {
	args := []string{"--format=json"}

	output, err := command.Runner.Run(openedPortsCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get opened ports: %w", err)
	}

	var openedPortsString []string

	err = json.Unmarshal(output, &openedPortsString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse opened ports: %w", err)
	}

	var openedPorts []Port

	for _, portString := range openedPortsString {
		var port Port

		_, err := fmt.Sscanf(portString, "%d/%s", &port.Port, &port.Protocol)
		if err != nil {
			return nil, fmt.Errorf("failed to parse port %s: %w", portString, err)
		}

		openedPorts = append(openedPorts, port)
	}

	return openedPorts, nil
}
