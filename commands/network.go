package commands

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

func (command Command) NetworkGet(opts *NetworkGetOptions) (*Network, error) {
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

	output, err := command.Runner.Run(networkGetCommand, args...)
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
