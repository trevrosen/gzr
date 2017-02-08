package comms

import (
	"encoding/json"
	"fmt"
	"io"
)

// GozerMetadataStore is the interface that should be implemented for any
// backend that needs to handle image data storage
type GozerMetadataStore interface {
	Store(string, ImageMetadata) error
	List(string) ([]Image, error)
	Cleanup()
	Delete(string) error
}

// Image is a struct unifying an image name with its metadata
type Image struct {
	Name string        `json:"name"`
	Meta ImageMetadata `json:"metadata"`
}

// ImageMetadata is a struct containing necessary metadata about a particular image
type ImageMetadata struct {
	GitCommit     string `json:"git-commit"`
	GitTag        string `json:"git-tag"`
	GitAnnotation string `json:"git-annotation"`
	GitOrigin     string `json:"git-origin"`
	CreatedAt     string `json:"created-at"`
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
