package etcd

import (
	"testing"
)

func TestCrud(t *testing.T) {
	/* Test: Create */
	err := InitializeEtcdKVStore()
	if err != nil {
		t.Error(err)
	}
	var exp_value int = 3
	var actual_value int
	/* Test: Update & Read */
	err = Put("vmeet", exp_value)
	if err != nil {
		t.Error("Put error")
	}
	err = Get("vmeet", &actual_value)
	if exp_value != actual_value {
		t.Error("Actual value mismatches with expected value")
	}
	/* Test: Delete */
	err = Del("vmeet")
	if err != nil {
		t.Error("Put error")
	}
	err = Get("vmeet", &actual_value)
	if err != nil {
		t.Error("The value should not exist")
	}
	/* Test end */
	t.Log("Pass CRUD Test!")
}
