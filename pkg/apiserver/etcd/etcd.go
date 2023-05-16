package etcd

import (
	"context"
<<<<<<< HEAD
	"encoding/json"
	"errors"
	"fmt"
	"time"
=======
	"fmt"
	"errors"
>>>>>>> apiserver

	clientv3 "go.etcd.io/etcd/client/v3"
)

var client *clientv3.Client = nil

const (
	etcdTimeout = 2 * time.Second
)

// Initialize a new etcd client
func InitializeEtcdKVStore() error {
	if client != nil {
		fmt.Println("etcd client has already existed!")
		return nil
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: etcdTimeout,
	})
	client = cli
	// TODO(Shao Wang): Close this client.
	return err
}

// Put write a single key value pair to etcd.
func Put(key, value string) error {
	_, err := client.KV.Put(context.Background(), key, value)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Get reads a single value for a given key.
// key: string type
func Get(key string) (string, error) {
	getResp, err := client.KV.Get(context.Background(), key)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if len(getResp.Kvs) == 0 {
		return "", errors.New("No value for this key is stored yet!")
	}
	return string(getResp.Kvs[0].Value), nil
}

// GetWithPrefix reads all value of which key starts with keyPrefix
// key: string type; val: a pointer of the type you desire
func GetWithPrefix(keyPrefix string) ([]string, error) {
	getResp, err := client.KV.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	retVals := make([]string, getResp.Count)
	for _, kv := range getResp.Kvs {
		retVals = append(retVals, string(kv.Value))
	}
	return retVals, nil
}

// Del delete a key value pair for a given key.
func Del(key string) error {
	_, err := client.KV.Delete(context.Background(), key)
	return err
}

// Delete k-v pairs of which key starts with string keyPrefix
func DelWithPrefix(keyPrefix string) error {
	_, err := client.KV.Delete(context.Background(), keyPrefix, clientv3.WithPrefix())
	return err
}

// Delete all records
func DelAll() error {
	return DelWithPrefix("")
}
// // Watch invoke a handler function on the change of a given key
// func Watch(key string) {
// 	watchCh := client.Watch(context.Background(), key)
// 	// TODO(Shao Wang): Replace the following handler function.
// 	go func() {
// 		for res := range watchCh {
// 			key := res.Events[0].Kv.Key
// 			value := string(res.Events[0].Kv.Value)
// 			fmt.Printf("Watch key %s's value changed to %s\n", key, value)
// 		}
// 	}()
// }

type WatchResult struct {
	ObjectType string
	ActionType int //0 for apply, 1 for modify, 2 for delete
	Payload    []byte
}
