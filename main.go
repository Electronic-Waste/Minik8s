package main

import (
	"fmt"

	"vmeet.io/minik8s/pkg/apiserver/etcd"
)

func main() {
	_ = etcd.InitializeEtcdKVStore()
	/* Test for put */
	etcd.Put("vmeet", "3")
	/* Test for get */
	value, _ := etcd.Get("vmeet")
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
	/* Test for del */
	etcd.Del("vmeet")
	value, _ = etcd.Get("vmeet")
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
}