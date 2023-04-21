package main

import (
	"fmt"

	"vmeet.io/minik8s/pkg/apiserver/etcd"
)

func main() {
	_ = etcd.InitializeEtcdKVStore()
	/* Test for watch */
	etcd.Watch("vmeet")
	/* Test for put */
	etcd.Put("vmeet", "3")
	value, _ := etcd.Get("vmeet")
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
	/* Test for put & watch */
	etcd.Put("vmeet", "4")
	value, _ = etcd.Get("vmeet")
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
	/* Test for del */
	etcd.Del("vmeet")
	value, _ = etcd.Get("vmeet")
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
}