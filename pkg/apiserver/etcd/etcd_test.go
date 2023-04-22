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
	value, _ := Get("vmeet")
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
	/* Test for put & watch */
	Put("vmeet", "4")
	value, _ = Get("vmeet")
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
	/* Test for del */
	Del("vmeet")
	value, _ = Get("vmeet")
	fmt.Printf("Key: %s -> Value: %s\n", "vmeet", value)
}