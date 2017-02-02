package comms

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coreos/etcd/client"
	"github.com/spf13/viper"
)

func createClient() (client.KeysAPI, error) {
	cxnString := fmt.Sprintf("%s:%s", viper.GetString("datastore.host"), viper.GetString("datastore.port"))
	cfg := client.Config{
		Endpoints: []string{cxnString},
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	kAPI := client.NewKeysAPI(c)
	return kAPI, nil
}

func storeEtcd(imageId string, imageMetadata ImageMetadata) error {
	data, err := json.Marshal(imageMetadata)
	if err != nil {
		return err
	}
	kAPI, err := createClient()
	if err != nil {
		return err
	}
	_, err = kAPI.Create(context.Background(), fmt.Sprintf("/%s", imageId), string(data))
	if err != nil {
		return err
	}

	fmt.Println("Storing in etcd")
	return nil
}

// func deleteEtcd(imageId string) error {
// 	_, err = kAPI.Delete(context.Background(), "/foo", &client.DeleteOptions{PrevValue: "bar"})
// 	if err != nil {
// 		// handle error
// 	}
//      return err
// }

// func getEtcd(imageId string) error {
// 	val, err := kAPI.Get(context.Background(), "/foo", &client.GetOptions{Recursive: true, Sort: false})
// 	if err != nil {
// 		// handle error
// 	}
// 	fmt.Printf("%+v\n", val)
// 	return err
// }
