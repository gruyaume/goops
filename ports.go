package goops

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

func SetPorts(ports []*Port) error {
	openedPorts, err := OpenedPorts()
	if err != nil {
		return fmt.Errorf("failed to get opened ports: %w", err)
	}

	desiredMap := make(map[string]*Port)

	for _, port := range ports {
		key := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		desiredMap[key] = port
	}

	openedMap := make(map[string]*Port)

	for _, port := range openedPorts {
		key := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		openedMap[key] = port
	}

	// Open ports that are desired but not currently opened.
	for key, port := range desiredMap {
		if _, exists := openedMap[key]; !exists {
			if err := OpenPort(port.Port, port.Protocol); err != nil {
				return fmt.Errorf("failed to open port %s: %w", key, err)
			}
		}
	}

	// Close ports that are currently opened but not desired.
	for key, port := range openedMap {
		if _, exists := desiredMap[key]; !exists {
			if err := ClosePort(port.Port, port.Protocol); err != nil {
				return fmt.Errorf("failed to close port %s: %w", key, err)
			}
		}
	}

	return nil
}

func OpenPort(port int, protocol string) error {
	commandRunner := GetRunner()

	if port < 0 || port > 65535 {
		return fmt.Errorf("port %d is out of range", port)
	}

	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("protocol %s is not supported", protocol)
	}

	args := []string{fmt.Sprintf("%d/%s", port, protocol)}

	_, err := commandRunner.Run(openPortCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to open port %d/%s: %w", port, protocol, err)
	}

	return nil
}

func ClosePort(port int, protocol string) error {
	commandRunner := GetRunner()

	if port < 0 || port > 65535 {
		return fmt.Errorf("port %d is out of range", port)
	}

	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("protocol %s is not supported", protocol)
	}

	args := []string{fmt.Sprintf("%d/%s", port, protocol)}

	_, err := commandRunner.Run(closePortCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to open port %d/%s: %w", port, protocol, err)
	}

	return nil
}

func OpenedPorts() ([]*Port, error) {
	commandRunner := GetRunner()

	args := []string{"--format=json"}

	output, err := commandRunner.Run(openedPortsCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get opened ports: %w", err)
	}

	var openedPortsString []string

	err = json.Unmarshal(output, &openedPortsString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse opened ports: %w", err)
	}

	var openedPorts []*Port

	for _, portString := range openedPortsString {
		var port Port

		_, err := fmt.Sscanf(portString, "%d/%s", &port.Port, &port.Protocol)
		if err != nil {
			return nil, fmt.Errorf("failed to parse port %s: %w", portString, err)
		}

		openedPorts = append(openedPorts, &port)
	}

	return openedPorts, nil
}
