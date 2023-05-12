package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var client *clientv3.Client

const (
	etcdTimeout = 2 * time.Second
)

// Initialize a new etcd client
func InitializeEtcdKVStore() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: etcdTimeout,
	})
	client = cli
	// TODO(Shao Wang): Close this client.
	return err
}

// Put write a single key value pair to etcd.
func Put(key string, value interface{}) error {
	jsonVal, _ := json.Marshal(value)
	_, err := client.KV.Put(context.Background(), key, string(jsonVal))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Get reads a single value for a given key.
// key: string type; val: a pointer of the type you desire
func Get(key string, val interface{}) error {
	getResp, err := client.KV.Get(context.Background(), key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if len(getResp.Kvs) != 1 {
		return errors.New("Should and should only get one value")
	}
	return json.Unmarshal(getResp.Kvs[0].Value, val)
}

// Del delete a key value pair for a given key.
func Del(key string) error {
	_, err := client.KV.Delete(context.Background(), key)
	return err
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
