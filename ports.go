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

type Protocol string

const (
	ProtocolTCP  Protocol = "tcp"
	ProtocolUDP  Protocol = "udp"
	ProtocolICMP Protocol = "icmp"
)

type Port struct {
	Port     int
	Protocol Protocol
}

// SetPorts sets the desired ports for the unit.
// It opens ports that are desired but not currently opened, and closes ports that are currently opened but not desired.
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

// OpenPort registers a request to open the specified port.
// The port must be between 0 and 65535, and the protocol must be one of tcp, udp, or icmp.
// If the protocol is icmp, the port argument is ignored.
func OpenPort(port int, protocol Protocol) error {
	commandRunner := GetCommandRunner()

	if port < 0 || port > 65535 {
		return fmt.Errorf("port %d is out of range", port)
	}

	arg := fmt.Sprintf("%d/%s", port, protocol)

	switch protocol {
	case ProtocolTCP, ProtocolUDP:
	case ProtocolICMP:
		arg = "icmp"
	default:
		return fmt.Errorf("invalid protocol: %s, must be one of tcp, udp, or icmp", protocol)
	}

	args := []string{arg}

	_, err := commandRunner.Run(openPortCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to open port %d/%s: %w", port, protocol, err)
	}

	return nil
}

// ClosePort registers a request to close the specified port.
// The port must be between 0 and 65535, and the protocol must be one of tcp, udp, or icmp.
// If the protocol is icmp, the port argument is ignored.
func ClosePort(port int, protocol Protocol) error {
	commandRunner := GetCommandRunner()

	if port < 0 || port > 65535 {
		return fmt.Errorf("port %d is out of range", port)
	}

	arg := fmt.Sprintf("%d/%s", port, protocol)

	switch protocol {
	case ProtocolTCP, ProtocolUDP:
	case ProtocolICMP:
		arg = "icmp"
	default:
		return fmt.Errorf("invalid protocol: %s, must be one of tcp, udp, or icmp", protocol)
	}

	args := []string{arg}

	_, err := commandRunner.Run(closePortCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to open port %d/%s: %w", port, protocol, err)
	}

	return nil
}

// List all ports opened by the unit.
func OpenedPorts() ([]*Port, error) {
	commandRunner := GetCommandRunner()

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
