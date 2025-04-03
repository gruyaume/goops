package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestNetworkGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"bind-addresses":[{"mac-address":"","interface-name":"","addresses":[{"hostname":"","value":"10.1.107.220","cidr":""}]}],"egress-subnets":["10.152.183.78/32"],"ingress-addresses":["10.152.183.78"]}`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	result, err := command.NetworkGet("whatever", false, false, false, false, "")
	if err != nil {
		t.Fatalf("NetworkGet returned an error: %v", err)
	}

	expected := commands.Network{
		BindAddresses: []commands.BindAddress{
			{
				InterfaceName: "",
				Addresses: []commands.Address{
					{
						Value: "10.1.107.220",
						CIDR:  "",
					},
				},
			},
		},
		EgressSubnets: []string{
			"10.152.183.78/32",
		},
		IngressAddresses: []string{
			"10.152.183.78",
		},
	}
	if len(result.BindAddresses) != len(expected.BindAddresses) {
		t.Fatalf("Expected %d bind addresses, got %d", len(expected.BindAddresses), len(result.BindAddresses))
	}

	if len(result.EgressSubnets) != len(expected.EgressSubnets) {
		t.Fatalf("Expected %d egress subnets, got %d", len(expected.EgressSubnets), len(result.EgressSubnets))
	}

	if len(result.IngressAddresses) != len(expected.IngressAddresses) {
		t.Fatalf("Expected %d ingress addresses, got %d", len(expected.IngressAddresses), len(result.IngressAddresses))
	}

	if result.BindAddresses[0].InterfaceName != expected.BindAddresses[0].InterfaceName {
		t.Fatalf("Expected %q, got %q", expected.BindAddresses[0].InterfaceName, result.BindAddresses[0].InterfaceName)
	}

	if result.BindAddresses[0].Addresses[0].Value != expected.BindAddresses[0].Addresses[0].Value {
		t.Fatalf("Expected %q, got %q", expected.BindAddresses[0].Addresses[0].Value, result.BindAddresses[0].Addresses[0].Value)
	}

	if result.BindAddresses[0].Addresses[0].CIDR != expected.BindAddresses[0].Addresses[0].CIDR {
		t.Fatalf("Expected %q, got %q", expected.BindAddresses[0].Addresses[0].CIDR, result.BindAddresses[0].Addresses[0].CIDR)
	}

	if fakeRunner.Command != commands.NetworkGetCommand {
		t.Errorf("Expected command %q, got %q", commands.NetworkGetCommand, fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "whatever" {
		t.Errorf("Expected argument %q, got %q", "whatever", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}
}
