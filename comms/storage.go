package comms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ImageStorer
type ImageStorageInterface interface {
	StoreImage(string, string) error
	GetImage(string) ([]Image, error)
}

type Image struct {
	ImageID   string        `json:"image_id"`
	ImageMeta ImageMetadata `json:"image_metadata"`
}

type ImageMetadata struct {
	GitCommit     string `json:"git-commit"`
	GitTag        string `json:"git-tag"`
	GitAnnotation string `json:"git-annotation"`
	CreatedAt     string `json:"created-at"`
}

type EtcdImageStorer struct{}

func (storer *EtcdImageStorer) GetImage(imageName string) ([]Image, error) {
	images, err := getEtcd(imageName)
	if err != nil {
		return []Image{}, err
	}
	return images, nil
}

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
