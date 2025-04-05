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

type SetPortsOptions struct {
	Ports []*Port
}

type OpenPortOptions struct {
	Port     int
	Protocol string // allowed values: tcp, udp
}

type ClosePortOptions struct {
	Port     int
	Protocol string // allowed values: tcp, udp
}

func (command Command) SetPorts(opts *SetPortsOptions) error {
	openedPorts, err := command.OpenedPorts()
	if err != nil {
		return fmt.Errorf("failed to get opened ports: %w", err)
	}

	desiredMap := make(map[string]*Port)

	for _, port := range opts.Ports {
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
			openPortOpts := &OpenPortOptions{
				Port:     port.Port,
				Protocol: port.Protocol,
			}
			if err := command.OpenPort(openPortOpts); err != nil {
				return fmt.Errorf("failed to open port %s: %w", key, err)
			}
		}
	}

	// Close ports that are currently opened but not desired.
	for key, port := range openedMap {
		if _, exists := desiredMap[key]; !exists {
			closePortOpts := &ClosePortOptions{
				Port:     port.Port,
				Protocol: port.Protocol,
			}
			if err := command.ClosePort(closePortOpts); err != nil {
				return fmt.Errorf("failed to close port %s: %w", key, err)
			}
		}
	}

	return nil
}

func (command Command) OpenPort(opts *OpenPortOptions) error {
	if opts.Port < 0 || opts.Port > 65535 {
		return fmt.Errorf("port %d is out of range", opts.Port)
	}

	if opts.Protocol != "tcp" && opts.Protocol != "udp" {
		return fmt.Errorf("protocol %s is not supported", opts.Protocol)
	}

	args := []string{fmt.Sprintf("%d/%s", opts.Port, opts.Protocol)}

	_, err := command.Runner.Run(openPortCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to open port %d/%s: %w", opts.Port, opts.Protocol, err)
	}

	return nil
}

func (command Command) ClosePort(opts *ClosePortOptions) error {
	if opts.Port < 0 || opts.Port > 65535 {
		return fmt.Errorf("port %d is out of range", opts.Port)
	}

	if opts.Protocol != "tcp" && opts.Protocol != "udp" {
		return fmt.Errorf("protocol %s is not supported", opts.Protocol)
	}

	args := []string{fmt.Sprintf("%d/%s", opts.Port, opts.Protocol)}

	_, err := command.Runner.Run(closePortCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to open port %d/%s: %w", opts.Port, opts.Protocol, err)
	}

	return nil
}

func (command Command) OpenedPorts() ([]*Port, error) {
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
