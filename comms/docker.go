package comms

import (
	"os"
	"os/exec"
)

func BuildDocker(args ...string) error {
	args = append([]string{"build"}, args...)
	docker := exec.Command("docker", args...)
	docker.Stdout = os.Stdout
	docker.Stderr = os.Stderr
	err := docker.Run()
	if err != nil {
		return err
	}
	return nil
}

func PushDocker(name string) error {
	push := exec.Command("docker", "push", name)
	push.Stdout = os.Stdout
	push.Stderr = os.Stderr

	err := push.Run()
	if err != nil {
		return err
	}
	return nil
}
