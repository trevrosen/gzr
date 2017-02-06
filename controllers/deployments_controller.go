package controllers

import (
	"log"
	"net/http"

	"github.com/bypasslane/gzr/comms"
	"github.com/gorilla/mux"
)

// listDeploymentsHandler lists deployments in the Kubernetes instance
func listDeploymentsHandler(k8sConn *comms.K8sConnection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deployments, err := k8sConn.ListDeployments(k8sConn.Namespace)
		// TODO: differentiate between legit errors and unhandleable errors
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonData, err := deployments.SerializeForWire()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	})
}

// getDeploymentHandler gets a single Deployment by name from the Kubernetes instance
func getDeploymentHandler(k8sConn *comms.K8sConnection) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		if name == "" {
			log.Println("name param required for this path")
			w.WriteHeader(http.StatusBadRequest)
		}
		deployment, err := k8sConn.GetDeployment(k8sConn.Namespace, name)

		if err == comms.ErrContainerNotFound {
			w.WriteHeader(http.StatusNotFound)
		}

		// TODO: catch other kinds of errors

		jsonData, err := deployment.SerializeForWire()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	})
}
