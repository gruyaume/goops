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

func GetNetwork(bindingName string) (*Network, error) {
	commandRunner := GetCommandRunner()

	var args []string

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

	return &network, nil
}
