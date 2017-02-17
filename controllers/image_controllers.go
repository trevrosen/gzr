package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/bypasslane/gzr/comms"
	"github.com/gorilla/mux"
)

func getImagesHandler(imageStore comms.GzrMetadataStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		name, err := url.QueryUnescape(name)
		if name == "" {
			log.Println("name param required for this path")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		images, err := imageStore.List(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(images.Images) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jsonData, err := images.SerializeForWire()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	})
}

func getImageHandler(imageStore comms.GzrMetadataStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		name, err := url.QueryUnescape(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		version := mux.Vars(r)["version"]
		if name == "" || version == "" {
			log.Println("name and version required for this path")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		image, err := imageStore.Get(fmt.Sprintf("%s:%s", name, version))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// If empty Image, one wasn't found
		if image == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jsonData, err := image.SerializeForWire()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	})
}
