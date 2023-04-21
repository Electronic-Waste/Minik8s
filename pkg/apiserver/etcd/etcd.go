package etcd

import (
	"time"
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var client *clientv3.Client

const (
	etcdTimeout =  2 * time.Second
)

// Initialize a new etcd client
func InitializeEtcdKVStore() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: etcdTimeout,
	})
	client = cli
	return err
}

// Put write a single key value pair to etcd.
func Put(key, value string) error {
	_, err := client.KV.Put(context.Background(), key, value)
	return err
}

// Get reads a single value for a given key.
func Get(key string) (string, error) {
	getResp, err := client.KV.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	kvs := getResp.Kvs
	if len(kvs) != 1 {
		return "", fmt.Errorf("expected exactly on value for key %s but got %d", key, len(kvs))
	}

	return string(kvs[0].Value), nil
}

// Del delete a key value pair for a given key.
func Del(key string) error {
	_, err := client.KV.Delete(context.Background(), key)
	return err
}

func Watch(key string) {
	watchCh := client.Watch(context.Background(), key)
	go func() {
		for res := range watchCh {
			key := res.Events[0].Kv.Key
			value := string(res.Events[0].Kv.Value)
			fmt.Printf("Watch key %s's value changed to %s\n", key, value)
		}
	}()
}

