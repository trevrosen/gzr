package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/bypasslane/gzr/comms"
)

func emptyDeploymentsList() (*comms.GzrDeploymentList, error) {
	return nil, nil
}

func populatedDeploymentsList() (*comms.GzrDeploymentList, error) {
	return &comms.GzrDeploymentList{}, nil
}

func emptyGetDeployment(deploymentName string) (*comms.GzrDeployment, error) {
	return nil, comms.ErrDeploymentNotFound
}

func populatedGetDeployment(deploymentName string) (*comms.GzrDeployment, error) {
	return &comms.GzrDeployment{}, nil
}

func successfulUpdateDeployment(dci *comms.DeploymentContainerInfo) (*comms.GzrDeployment, error) {
	return &comms.GzrDeployment{}, nil
}

func failUpdateDeploymentNoDeployment(dci *comms.DeploymentContainerInfo) (*comms.GzrDeployment, error) {
	return &comms.GzrDeployment{}, comms.ErrDeploymentNotFound
}

func failUpdateDeploymentNoContainer(dci *comms.DeploymentContainerInfo) (*comms.GzrDeployment, error) {
	return &comms.GzrDeployment{}, comms.ErrContainerNotFound
}

// Sends an HTTP request to provided server:
// GET /deployments
func getDeploymentsList(server *httptest.Server) (*http.Response, error) {
	client := new(http.Client)
	req, _ := http.NewRequest("GET", server.URL+"/deployments", nil)
	return client.Do(req)
}

// Sends an HTTP request to provided server:
// GET /deployments/{name}
func getDeployment(server *httptest.Server) (*http.Response, error) {
	client := new(http.Client)
	req, _ := http.NewRequest("GET", server.URL+"/deployments/name", nil)
	return client.Do(req)
}

// Sends an HTTP request to provided server:
// PUT /deployments/{name}
func updateDeployment(server *httptest.Server) (*http.Response, error) {
	client := new(http.Client)
	payloadSource := `{"container_name": "foobaricus", "image": "foobar:1.2.3"}`
	reader := strings.NewReader(payloadSource)
	req, _ := http.NewRequest("PUT", server.URL+"/deployments/name", reader)
	return client.Do(req)
}
