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
	mockImageStore := &comms.MockStore{}

	server := httptest.NewServer(App(mockK8sConn, mockImageStore))
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
	mockImageStore := &comms.MockStore{}

	server := httptest.NewServer(App(mockK8sConn, mockImageStore))
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
	mockImageStore := &comms.MockStore{}

	server := httptest.NewServer(App(mockK8sConn, mockImageStore))
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
	mockImageStore := &comms.MockStore{}

	server := httptest.NewServer(App(mockK8sConn, mockImageStore))
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
	mockImageStore := &comms.MockStore{}

	server := httptest.NewServer(App(mockK8sConn, mockImageStore))
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
	mockImageStore := &comms.MockStore{}

	server := httptest.NewServer(App(mockK8sConn, mockImageStore))
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
	mockImageStore := &comms.MockStore{}

	server := httptest.NewServer(App(mockK8sConn, mockImageStore))
	defer server.Close()
	res, err := updateDeployment(server)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected %v, but received %v", http.StatusNotFound, res.Status)
	}
}
