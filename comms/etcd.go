package comms

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
)

func createClient() (*clientv3.Client, error) {
	cxnString := fmt.Sprintf("%s:%s", viper.GetString("datastore.host"), viper.GetString("datastore.port"))
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{cxnString},
	})
	if err != nil {
		return nil, err
	}
	clientv3.SetLogger(clientv3.GetLogger())
	return cli, nil
}

func createEtcdKey(imageName string) (string, error) {
	splitName := strings.Split(imageName, ":")
	if len(splitName) != 2 {
		return "", fmt.Errorf("IMAGE_NAME must be formatted as NAME:VERSION and must contain only the seperating colon")
	}
	now := time.Now()
	nowString := fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day())
	return fmt.Sprintf("%s:%s:%s", splitName[0], nowString, splitName[1]), nil
}

func storeEtcd(imageName string, imageMetadata ImageMetadata) error {
	data, err := json.Marshal(imageMetadata)
	if err != nil {
		return err
	}

	etcdClient, err := createClient()
	defer etcdClient.Close()

	kv := clientv3.NewKV(etcdClient)
	if err != nil {
		return err
	}

	key, err := createEtcdKey(imageName)
	if err != nil {
		return err
	}
	_, err = kv.Put(context.Background(), key, string(data))
	if err != nil {
		return err
	}

	return nil

}

func getEtcd(imageName string) ([]Image, error) {
	etcdClient, err := createClient()
	defer etcdClient.Close()

	kv := clientv3.NewKV(etcdClient)
	if err != nil {
		return []Image{}, err
	}

	resp, err := kv.Get(context.Background(), fmt.Sprintf("%s:", imageName), clientv3.WithPrefix())
	if err != nil {
		return []Image{}, err
	}
	val, err := extractVal(resp)
	return val, err
}

func extractVal(resp *clientv3.GetResponse) ([]Image, error) {
	if len(resp.Kvs) < 1 {
		return []Image{}, fmt.Errorf("No results found")
	}
	var images []Image
	for _, kv := range resp.Kvs {
		var meta ImageMetadata
		json.Unmarshal(kv.Value, &meta) // TODO: Handle error
		images = append(images, Image{ImageID: string(kv.Key), ImageMeta: meta})
	}
	return images, nil
}

func deleteEtcd(imageName string) error {
	etcdClient, err := createClient()
	defer etcdClient.Close()

	kv := clientv3.NewKV(etcdClient)
	if err != nil {
		return err
	}

	_, err = kv.Delete(context.Background(), fmt.Sprintf("/%s", imageName))
	if err != nil {
		return err
	}
	return err
}
