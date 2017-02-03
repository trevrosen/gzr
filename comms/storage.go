package comms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ImageStorer
type ImageStorageInterface interface {
	StoreImage(string, string)
}

type ImageMetadata struct {
	GitCommit     string `json:"git-commit"`
	GitTag        string `json:"git-tag"`
	GitAnnotation string `json:"git-annotation"`
	CreatedAt     string `json:"created-at"`
}

type EtcdImageStorer struct{}

func (storer *EtcdImageStorer) StoreImage(imageId string, dataPath string) {
	file, err := ioutil.ReadFile(dataPath)
	if err != nil {
		fmt.Println("Could not read metadata file")
		os.Exit(1)
	}

	var meta ImageMetadata
	err = json.Unmarshal(file, &meta)
	if err != nil {
		fmt.Println("Could not read metadata file")
		os.Exit(1)
	}

	err = storeEtcd(imageId, meta)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
