package goopstest_test

import (
	"fmt"

	"github.com/gruyaume/goops"
)

func Configure() error {
	isLeader, err := goops.IsLeader()
	if err != nil {
		return fmt.Errorf("could not check if unit is leader: %w", err)
	}

	if !isLeader {
		goops.SetUnitStatus(goops.StatusBlocked, "Unit is not leader")
		return nil
	}

	err = goops.SetPorts([]*goops.Port{
		{
			Port:     2111,
			Protocol: "tcp",
		},
	})
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	goops.SetUnitStatus(goops.StatusActive, fmt.Sprintf("Port %d/tcp is set", 2111))

	return nil
}
