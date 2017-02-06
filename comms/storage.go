package comms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ImageStorageInterface is the interface that should be implemented for any
// backend that needs to handle image data storage
type ImageStorageInterface interface {
	StoreImage(string, string) error
	GetImages(string) ([]Image, error)
}

// Image is a struct unifying an image name with its metadata
type Image struct {
	ImageName string        `json:"image_name"`
	ImageMeta ImageMetadata `json:"image_metadata"`
}

// ImageMetadata is a struct containing necessary metadata about a particular image
type ImageMetadata struct {
	GitCommit     string `json:"git-commit"`
	GitTag        string `json:"git-tag"`
	GitAnnotation string `json:"git-annotation"`
	CreatedAt     string `json:"created-at"`
}

// EtcdImageStorer is simply an empty struct to implement ImageStorageInterface
type EtcdImageStorer struct{}

// GetImages queries the etcd store for all images stored under a particular name
func (storer *EtcdImageStorer) GetImages(imageName string) ([]Image, error) {
	images, err := getEtcd(imageName)
	if err != nil {
		return []Image{}, err
	}
	return images, nil
}

// StoreImage stores the metadata about an image where the metadata is a path
// to a JSON-formatted file containing ImageMetadata fields
func (storer *EtcdImageStorer) StoreImage(imageName string, dataPath string) error {
	file, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return fmt.Errorf("Could not read metadata file")
	}

	var meta ImageMetadata
	err = json.Unmarshal(file, &meta)
	if err != nil {
		return fmt.Errorf("Could not read metadata file")
	}

	err = storeEtcd(imageName, meta)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}
