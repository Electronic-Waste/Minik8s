package etcd

import (
	"time"
	"context"
	"fmt"
	"errors"
	"encoding/json"

	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s.io/pkg/apis/core"
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
func Get(key string, val interface{}) (interface{}, error) {
	getResp, err := client.KV.Get(context.Background(), key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if len(getResp.Kvs) != 1 {
		return nil, errors.New("Should and should only get one value")
	}
	switch val.(type) {
	case core.Pod:
		var retVal core.Pod
		json.Unmarshal(getResp.Kvs[0].Value, &retVal)
		return retVal, nil
	case string:
		var retVal string
		json.Unmarshal(getResp.Kvs[0].Value, &retVal)
		return retVal, nil
	default:
		return nil, errors.New("etcd: Unsupported type!")
	}
	
}

// GetWithPrefix reads all value of which key starts with keyPrefix
// key: string type; val: a pointer of the type you desire
func GetWithPrefix(keyPrefix string, val interface{}) ([]interface{}, error) {
	getResp, err := client.KV.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	retVals := make([]interface{}, getResp.Count)
	switch val.(type) {
	case []core.Pod:
		var value core.Pod
		for _, kv := range getResp.Kvs {
			json.Unmarshal(kv.Value, &value)
			_ = append(retVals, value)
		}
		return retVals, nil
	
	case []string:
		var value string
		for _, kv := range getResp.Kvs {
			json.Unmarshal(kv.Value, &value)
			_ = append(retVals, value)
		}
		return retVals, nil
	default:
		return nil, errors.New("etcd: Unsupported type!")
	}
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

