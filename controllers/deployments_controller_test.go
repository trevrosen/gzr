package controllers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bypasslane/gzr/comms"
)

func TestListDeploymentsExist(t *testing.T) {
	mockK8sConn := &comms.MockK8sCommunicator{
		OnListDeployments: populatedDeploymentsList,
	}

	server := httptest.NewServer(App(mockK8sConn))
	defer server.Close()
	res, err := getDeploymentsList(server)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, but received %v", http.StatusOK, res.Status)
	}
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
	mockK8sConn := &comms.MockK8sCommunicator{
		OnGetDeployment: populatedGetDeployment,
	}

	server := httptest.NewServer(App(mockK8sConn))
	defer server.Close()
	res, err := getDeployment(server)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, but received %v", http.StatusOK, res.Status)
	}
}

func TestGetDeploymentNotFound(t *testing.T) {
	mockK8sConn := &comms.MockK8sCommunicator{
		OnGetDeployment: emptyGetDeployment,
	}

	server := httptest.NewServer(App(mockK8sConn))
	defer server.Close()
	res, err := getDeployment(server)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %v, but received %v", http.StatusNotFound, res.Status)
	}
}

func TestUpdateDeploymentAndCountainerFound(t *testing.T) {
	mockK8sConn := &comms.MockK8sCommunicator{
		OnGetDeployment:    populatedGetDeployment,
		OnUpdateDeployment: successfulUpdateDeployment,
	}

	server := httptest.NewServer(App(mockK8sConn))
	defer server.Close()
	res, err := updateDeployment(server)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, but received %v", http.StatusOK, res.Status)
	}
}

func TestUpdateDeploymentNotFound(t *testing.T) {
	mockK8sConn := &comms.MockK8sCommunicator{
		OnGetDeployment: emptyGetDeployment,
	}

	server := httptest.NewServer(App(mockK8sConn))
	defer server.Close()
	res, err := updateDeployment(server)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %v, but received %v", http.StatusNotFound, res.Status)
	}
}

func TestUpdateContainerNotFound(t *testing.T) {
	mockK8sConn := &comms.MockK8sCommunicator{
		OnGetDeployment:    populatedGetDeployment,
		OnUpdateDeployment: failUpdateDeploymentNoContainer,
	}

	server := httptest.NewServer(App(mockK8sConn))
	defer server.Close()
	res, err := updateDeployment(server)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %v, but received %v", http.StatusNotFound, res.Status)
	}
}
