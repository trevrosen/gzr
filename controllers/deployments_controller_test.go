package controllers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bypasslane/gzr/comms"
)

func TestListDeploymentsExist(t *testing.T) {
}

func TestListDeploymentsNone(t *testing.T) {
	mockK8sConn := &comms.MockK8sCommunicator{
		OnListDeployments: emptyDeploymentsList,
	}

	server := httptest.NewServer(App(mockK8sConn))
	defer server.Close()
	res, err := getDeploymentsList(server)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %v, but received %v", http.StatusNotFound, res.Status)
	}
}

func TestGetDeploymentFound(t *testing.T) {
}

func TestGetDeploymentNotFound(t *testing.T) {
}

func TestUpdateDeploymentFound(t *testing.T) {
}

func TestUpdateDeploymentNotFound(t *testing.T) {
}
