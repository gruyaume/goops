package commands

const (
	JujuRebootCommand = "juju-reboot"
)

func (command Command) JujuReboot(now bool) error {
	var args []string
	if now {
		args = append(args, "--now")
	}

	_, err := command.Runner.Run(JujuRebootCommand, args...)
	if err != nil {
		return err
	}

	return nil
}
