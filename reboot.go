package goops

const (
	jujuRebootCommand = "juju-reboot"
)

func Reboot(now bool) error {
	commandRunner := GetRunner()

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
