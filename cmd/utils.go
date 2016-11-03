package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type GozerConfig struct {
	ApiVersion  string `yaml:"apiVersion"`
	Deployments []Deployment
	Services    []Service
}

type Deployment struct {
	Name         string
	DefaultImage string   `yaml:"default-image"`
	TemplePath   string   `yaml:"temple-path"`
	BuildSteps   []string `yaml:"build-steps"`
}

type Service struct {
	Name         string
	ExternalPort int `yaml:"externalPort"`
	InternalPort int `yaml:"internalPort"`
	Type         string
	Selector     map[string]string
	Params       map[string]interface{}
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(-1)
}

// notify sends a formatted information line to stdout
func notify(msg string) {
	fmt.Printf("[-] %s\n", msg)
}

// kubeCtlExists returns true if the kubectl bin can be found and
// reports a supported version
func kubeCtlExists() bool {
	// Check $PATH for kubectl
	// if not found, return false
	// if found, grab version, check against supported version declaration
	return false
}

func gozerConfigFromFile() GozerConfig {
	var err error
	cwd, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/.gozer.yml", cwd)
	yamlBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		er(err)
	}
	config, err := gozerConfigFromBytes(yamlBytes)

	if err != nil {
		er(err)
	}
	return config
}

// gozerConfigFromBytes returns a GozerConfig unmarshalld from the bytes
func gozerConfigFromBytes(data []byte) (GozerConfig, error) {
	var err error
	config := GozerConfig{}
	err = yaml.Unmarshal(data, &config)
	return config, err
}
