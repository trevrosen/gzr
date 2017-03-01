package comms

import (
	"fmt"
	"log"
	"os/exec"
)

func BuildDocker(args ...string) error {
	docker := exec.Command("docker", args...)
	stdoutStderr, err := docker.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
	return nil
}
