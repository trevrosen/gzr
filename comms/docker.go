package comms

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/viper"
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
	tag, err := GetDockerTag()
	if err != nil {
		return err
	}
	args = append([]string{"build", "-t", tag}, args...)
	build := exec.Command("docker", args...)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	err = build.Run()
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

func GetDockerTag() (string, error) {
	hash, err := getCommitHash()
	if err != nil {
		return "", err
	}
	imageName, err := os.Getwd()
	if err != nil {
		return "", err
	}
	splitImageName := strings.Split(imageName, "/")
	name := splitImageName[len(splitImageName)-1]
	return fmt.Sprintf("%s/%s:%s.%s", viper.GetString("repository"), name, time.Now().Format("20060102"), hash), nil
}
