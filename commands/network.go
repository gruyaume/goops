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

func (command Command) NetworkGet(bindingName string, bindAddress bool, egressSubnets bool, ingressAddress bool, primaryAddress bool, relation string) (*Network, error) {
	var args []string

	if bindingName == "" {
		return nil, fmt.Errorf("binding name cannot be empty")
	}

	if bindAddress {
		args = append(args, "--bind-address")
	}

	if egressSubnets {
		args = append(args, "--egress-subnets")
	}

	if ingressAddress {
		args = append(args, "--ingress-address")
	}

	if primaryAddress {
		args = append(args, "--primary-address")
	}

	if relation != "" {
		args = append(args, "--relation="+relation)
	}

	args = append(args, bindingName, "--format=json")

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
