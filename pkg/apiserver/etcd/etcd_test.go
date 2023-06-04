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
	/* Test: Update & Read */
	err = Put("vmeet", exp_value)
	if err != nil {
		t.Error("Put error")
	}
	actual_value, err = Get("vmeet")
	t.Log(actual_value, exp_value)
	if exp_value != actual_value {
		t.Error("Actual value mismatches with expected value")
	}
	/* Test: GetAll */
	var values []string
	values, err = GetWithPrefix("vmeet")
	t.Log(values)
	if len(values) != 1 {
		t.Errorf("Should only have one value but its len is %d", len(values))
	}
	/* Test: Delete */
	err = Del("vmeet")
	if err != nil {
		t.Error("Del error")
	}
	actual_value, err = Get("vmeet")
	/* Test end */
	DelAll()
	t.Log("Pass CRUD Test!")
}
