package comms

import (
	"context"
	"encoding/json"
	"fmt"

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
	return cli, nil

}

func storeEtcd(imageId string, imageMetadata ImageMetadata) error {
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

	resp, err := kv.Put(context.Background(), fmt.Sprintf("/%s", imageId), string(data))
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil

}

func getEtcd(imageId string) (string, error) {
	etcdClient, err := createClient()
	defer etcdClient.Close()

	kv := clientv3.NewKV(etcdClient)
	if err != nil {
		return "", err
	}

	resp, err := kv.Get(context.Background(), fmt.Sprintf("/%s", imageId))
	if err != nil {
		return "", err
	}
	val, err := extractVal(resp)
	return val, err
}

func extractVal(resp *clientv3.GetResponse) (string, error) {
	if resp.Count > 1 {
		return "", fmt.Errorf("Error in retrieving information")
	}

	kv := resp.Kvs[0]
	return string(kv.Value), nil
}

func deleteEtcd(imageId string) error {
	etcdClient, err := createClient()
	defer etcdClient.Close()

	kv := clientv3.NewKV(etcdClient)
	if err != nil {
		return err
	}

	_, err = kv.Delete(context.Background(), fmt.Sprintf("/%s", imageId))
	if err != nil {
		return err
	}
	return err
}
