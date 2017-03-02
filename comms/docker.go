package comms

import (
	"log"
	"os"
	"os/exec"
)

func BuildDocker(args ...string) error {
	docker := exec.Command("docker", args...)
	docker.Stdout = os.Stdout
	docker.Stderr = os.Stderr
	err := docker.Run()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
