package comms

import (
	"encoding/json"
	"fmt"
	"io"
)

// GozerMetadataStore is the interface that should be implemented for any
// backend that needs to handle image data storage
type GozerMetadataStore interface {
	// Store stores image metadata with a name
	Store(string, ImageMetadata) error
	// List lists all of the images under a name
	List(string) ([]Image, error)
	// Cleanup allows the storage backend to clean up any connections, etc
	Cleanup()
	// Delete deletes all images under a nmae
	Delete(string) error
	// Get gets a single image with a version
	Get(string) (Image, error)
}

// Image is a struct unifying an image name with its metadata
type Image struct {
	// Name is the image's full name
	Name string `json:"name"`
	// Meta is the metadata related to the image
	Meta ImageMetadata `json:"metadata"`
}

// ImageMetadata is a struct containing necessary metadata about a particular image
type ImageMetadata struct {
	// GitCommit is the commit related to the built image
	GitCommit string `json:"git-commit"`
	// GitTag is the tag related to the built image if it exists
	GitTag string `json:"git-tag"`
	// GitAnnotation is the annotation related to the built image if it exists
	GitAnnotation string `json:"git-annotation"`
	// GitOrigin is the remote origin for the git repo associated with the image
	GitOrigin string `json:"git-origin"`
	// CreatedAt is the time the metadata was stored, with day granularity
	CreatedAt string `json:"created-at"`
}

// CreateMeta takes a ReadWriter and returns an instance of ImageMetadata
// after parsing
func CreateMeta(reader io.ReadWriter) (ImageMetadata, error) {
	var meta ImageMetadata
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&meta)
	if err != nil {
		return ImageMetadata{}, fmt.Errorf("Could not read metadata file")
	}
	return meta, nil
}