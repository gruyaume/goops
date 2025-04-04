package commands

const (
	jujuRebootCommand = "juju-reboot"
)

type JujuRebootOptions struct {
	Now bool
}

func (command Command) JujuReboot(opts *JujuRebootOptions) error {
	var args []string
	if opts.Now {
		args = append(args, "--now")
	}

	_, err := command.Runner.Run(jujuRebootCommand, args...)
	if err != nil {
		return err
	}

	return nil
}
