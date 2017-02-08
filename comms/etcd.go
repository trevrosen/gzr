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

// EtcdStorage implements GozerMetadataStore and has exported
// Etcd clients and KV accessors
type EtcdStorage struct {
	Client *clientv3.Client
	KV     clientv3.KV
}

func NewEtcdStorage() (*EtcdStorage, error) {
	newEtcd := &EtcdStorage{}
	cxnString := fmt.Sprintf("%s:%s", viper.GetString("datastore.host"), viper.GetString("datastore.port"))
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{cxnString},
	})
	if err != nil {
		return nil, err
	}
	newEtcd.Client = cli

	kv := clientv3.NewKV(cli)
	newEtcd.KV = kv
	return newEtcd, nil
}

// List queries the etcd store for all images stored under a particular name
func (storer *EtcdStorage) List(imageName string) ([]Image, error) {
	images, err := storer.getEtcd(imageName)
	if err != nil {
		return []Image{}, err
	}
	return images, nil
}

// Store stores the metadata about an image where the metadata is a path
// to a JSON-formatted file containing ImageMetadata fields
func (storer *EtcdStorage) Store(imageName string, meta ImageMetadata) error {
	err := storer.storeEtcd(imageName, meta)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func (storer *EtcdStorage) Cleanup() {
	storer.Client.Close()
}

func (storer *EtcdStorage) Delete(imageName string) error {
	_, err := storer.KV.Delete(context.Background(), imageName, clientv3.WithPrefix())
	return err
}

func (storer *EtcdStorage) createEtcdKey(imageName string) (string, error) {
	splitName := strings.Split(imageName, ":")
	if len(splitName) != 2 {
		return "", fmt.Errorf("IMAGE_NAME must be formatted as NAME:VERSION and must contain only the seperating colon")
	}
	now := time.Now()
	nowString := fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day())
	return fmt.Sprintf("%s:%s:%s", splitName[0], splitName[1], nowString), nil
}

func (storer *EtcdStorage) storeEtcd(imageName string, imageMetadata ImageMetadata) error {
	data, err := json.Marshal(imageMetadata)
	if err != nil {
		return err
	}

	key, err := storer.createEtcdKey(imageName)
	if err != nil {
		return err
	}
	_, err = storer.KV.Put(context.Background(), key, string(data))
	if err != nil {
		return err
	}

	return nil

}

func (storer *EtcdStorage) getEtcd(imageName string) ([]Image, error) {
	resp, err := storer.KV.Get(context.Background(), fmt.Sprintf("%s:", imageName), clientv3.WithPrefix())
	if err != nil {
		return []Image{}, err
	}
	val, err := storer.extractVal(resp)
	return val, err
}

func (storer *EtcdStorage) extractVal(resp *clientv3.GetResponse) ([]Image, error) {
	if len(resp.Kvs) < 1 {
		return []Image{}, fmt.Errorf("No results found")
	}
	var images []Image
	for _, kv := range resp.Kvs {
		var meta ImageMetadata
		json.Unmarshal(kv.Value, &meta) // TODO: Handle error
		images = append(images, Image{Name: string(kv.Key), Meta: meta})
	}
	return images, nil
}

func (storer *EtcdStorage) deleteEtcd(imageName string) error {
	_, err := storer.KV.Delete(context.Background(), fmt.Sprintf("/%s", imageName))
	if err != nil {
		return err
	}
	return err
}
