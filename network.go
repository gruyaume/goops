package goops

import (
	"encoding/json"
	"fmt"
)

const (
	networkGetCommand = "network-get"
)

type Address struct {
	Value string `json:"value"`
	CIDR  string `json:"cidr"`
}

type BindAddress struct {
	InterfaceName string    `json:"interface-name"`
	Addresses     []Address `json:"addresses"`
}

type Network struct {
	BindAddresses    []BindAddress `json:"bind-addresses"`
	IngressAddresses []string      `json:"ingress-addresses"`
	EgressSubnets    []string      `json:"egress-subnets"`
}

type NetworkGetOptions struct {
	BindingName    string
	BindAddress    bool
	EgressSubnets  bool
	IngressAddress bool
	PrimaryAddress bool
	Relation       string
}

func GetNetwork(opts *NetworkGetOptions) (*Network, error) {
	commandRunner := GetCommandRunner()

	var args []string

	if opts.BindingName == "" {
		return nil, fmt.Errorf("binding name cannot be empty")
	}

	if opts.BindAddress {
		args = append(args, "--bind-address")
	}

	if opts.EgressSubnets {
		args = append(args, "--egress-subnets")
	}

	if opts.IngressAddress {
		args = append(args, "--ingress-address")
	}

	if opts.PrimaryAddress {
		args = append(args, "--primary-address")
	}

	if opts.Relation != "" {
		args = append(args, "--relation="+opts.Relation)
	}

	args = append(args, opts.BindingName, "--format=json")

	output, err := commandRunner.Run(networkGetCommand, args...)
	if err != nil {
		return nil, err
	}

	var network Network

	err = json.Unmarshal(output, &network)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network data: %w", err)
	}

	return &network, nil
}

func GetNetworkBindAddress(bindingName string) (string, error) {
	commandRunner := GetCommandRunner()

	var args []string

	args = append(args, "--bind-address")

	args = append(args, bindingName, "--format=json")

	output, err := commandRunner.Run(networkGetCommand, args...)
	if err != nil {
		return "", err
	}

	var bindAddress string

	err = json.Unmarshal(output, &bindAddress)
	if err != nil {
		return "", fmt.Errorf("failed to parse network data: %w", err)
	}

	return bindAddress, nil
}

func GetNetworkIngressAddress(bindingName string) (string, error) {
	commandRunner := GetCommandRunner()

	var args []string

	if bindingName == "" {
		return "", fmt.Errorf("binding name cannot be empty")
	}

	args = append(args, "--ingress-address")

	args = append(args, bindingName, "--format=json")

	output, err := commandRunner.Run(networkGetCommand, args...)
	if err != nil {
		return "", err
	}

	var ingressAddress string

	err = json.Unmarshal(output, &ingressAddress)
	if err != nil {
		return "", fmt.Errorf("failed to parse network data: %w", err)
	}

	return ingressAddress, nil
}

func GetNetworkEgressSubnets(bindingName string) ([]string, error) {
	commandRunner := GetCommandRunner()

	var args []string

	if bindingName == "" {
		return nil, fmt.Errorf("binding name cannot be empty")
	}

	args = append(args, "--egress-subnets")

	args = append(args, bindingName, "--format=json")

	output, err := commandRunner.Run(networkGetCommand, args...)
	if err != nil {
		return nil, err
	}

	var egressSubnets []string

	err = json.Unmarshal(output, &egressSubnets)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network data: %w", err)
	}

	return egressSubnets, nil
}
