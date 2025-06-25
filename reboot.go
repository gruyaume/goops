package goops

const (
	jujuRebootCommand = "juju-reboot"
)

// Reboot causes the host machine to reboot, after stopping all containers hosted on the machine.
func Reboot(now bool) error {
	commandRunner := GetCommandRunner()

	var args []string
	if now {
		args = append(args, "--now")
	}

	_, err := commandRunner.Run(jujuRebootCommand, args...)
	if err != nil {
		return err
	}

	return nil
}
