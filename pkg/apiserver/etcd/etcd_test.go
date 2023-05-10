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
	var exp_value string = "3"
	var actual_value string
	var interface_value interface{}
	/* Test: Update & Read */
	err = Put("vmeet", exp_value)
	if err != nil {
		t.Error("Put error")
	}
	interface_value, err = Get("vmeet", actual_value)
	actual_value = interface_value.(string)
	if exp_value != actual_value {
		t.Error("Actual value mismatches with expected value")
	}
	/* Test: Delete */
	err = Del("vmeet")
	if err != nil {
		t.Error("Put error")
	}
	interface_value, err = Get("vmeet", actual_value)
	t.Log(actual_value)
	if interface_value != nil {
		t.Error("The value should not exist")
	}
	/* Test end */
	t.Log("Pass CRUD Test!")
}
