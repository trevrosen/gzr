package comms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

type ImageMetadata struct {
	GitCommit     string `json:"git-commit"`
	GitTag        string `json:"git-tag"`
	GitAnnotation string `json:"git-annotation"`
	CreatedAt     string `json:"created-at"`
}

func StoreImage(args []string) {
	file, err := ioutil.ReadFile(args[1])
	if err != nil {
		fmt.Println("Could not read metadata file")
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	var meta ImageMetadata
	err = json.Unmarshal(file, &meta)
	if err != nil {
		fmt.Println("Could not read metadata file")
		os.Exit(1)
	}
	fmt.Printf("Results: %v\n", meta)

	switch viper.GetString("datastore.type") {
	case "etcd":
		err = storeEtcd(args[0], meta)
	default:
		err = fmt.Errorf("%s is not a valid datastore type", viper.GetString("datastore.type"))
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
