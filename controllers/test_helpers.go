package controllers

import (
	"net/http"
	"net/http/httptest"

	"github.com/bypasslane/gzr/comms"
)

func emptyDeploymentsList(namespace string) (*comms.GzrDeploymentList, error) {
	return nil, nil
}

// Sends an HTTP request to provided server:
// GET /deployments
func getDeploymentsList(server *httptest.Server) (*http.Response, error) {
	client := new(http.Client)
	req, _ := http.NewRequest("GET", server.URL+"/deployments", nil)
	return client.Do(req)
}
