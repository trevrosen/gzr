package comms

import (
	"os"
	"os/exec"
)

type ImageBuilder interface {
	Build(...string) error
	Push(string) error
}

type DockerBuilder struct{}

func NewDockerBuilder() *DockerBuilder {
	return &DockerBuilder{}
}

func (docker *DockerBuilder) Build(args ...string) error {
	args = append([]string{"build"}, args...)
	build := exec.Command("docker", args...)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	err := build.Run()
	if err != nil {
		return err
	}
	return nil
}

func (docker *DockerBuilder) Push(name string) error {
	push := exec.Command("docker", "push", name)
	push.Stdout = os.Stdout
	push.Stderr = os.Stderr

	err := push.Run()
	if err != nil {
		return err
	}
	return nil
}
