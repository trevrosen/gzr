package controllers

import (
	"net/http"

	"github.com/bypasslane/gzr/comms"
	"github.com/bypasslane/gzr/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func App(k8sConn *comms.K8sConnection) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/deployments", listDeploymentsHandler(k8sConn)).Methods("GET")
	router.HandleFunc("/deployments/{name}", getDeploymentHandler(k8sConn)).Methods("GET")
	router.HandleFunc("/deployments/{name}", updateDeploymentHandler(k8sConn)).Methods("PUT")
	n := negroni.Classic()
	n.Use(middleware.NewContentType()) // Ensure response Content-Type header is always "application/json"
	n.UseHandler(router)
	return n
}

// homeHandler handles requests to the root of the server
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("CHOOSE THE FORM"))
}
