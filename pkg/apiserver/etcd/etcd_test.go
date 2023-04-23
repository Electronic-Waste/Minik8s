package etcd

import (
	"fmt"
)

func EtcdTest() {
	_ = InitializeEtcdKVStore()
	/* Test for watch */
	Watch("vmeet")
	/* Test for put */
	Put("vmeet", "3")
	var value string
	Get("vmeet", &value)
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
	/* Test for put & watch */
	Put("vmeet", "4")
	Get("vmeet", &value)
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
	/* Test for del */
	Del("vmeet")
	Get("vmeet", &value)
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
}