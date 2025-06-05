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
	commandRunner := GetRunner()

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

func GetNetworkBindAddresses(bindingName string) ([]BindAddress, error) {
	commandRunner := GetRunner()

	var args []string

	if bindingName == "" {
		return nil, fmt.Errorf("binding name cannot be empty")
	}

	args = append(args, "--bind-address")

	args = append(args, bindingName, "--format=json")

	output, err := commandRunner.Run(networkGetCommand, args...)
	if err != nil {
		return nil, err
	}

	var network Network

	err = json.Unmarshal(output, &network)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network data: %w", err)
	}

	return network.BindAddresses, nil
}

func GetNetworkIngressAddresses(bindingName string) ([]string, error) {
	commandRunner := GetRunner()

	var args []string

	if bindingName == "" {
		return nil, fmt.Errorf("binding name cannot be empty")
	}

	args = append(args, "--ingress-address")

	args = append(args, bindingName, "--format=json")

	output, err := commandRunner.Run(networkGetCommand, args...)
	if err != nil {
		return nil, err
	}

	var network Network

	err = json.Unmarshal(output, &network)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network data: %w", err)
	}

	return network.IngressAddresses, nil
}

func GetNetworkEgressSubnets(bindingName string) ([]string, error) {
	commandRunner := GetRunner()

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

	var network Network

	err = json.Unmarshal(output, &network)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network data: %w", err)
	}

	return network.EgressSubnets, nil
}
