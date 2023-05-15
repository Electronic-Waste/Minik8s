package pod

import (
	"testing"
	"encoding/json"

)

func TestStringArrToBytes(t *testing.T) {
	arr := []string{"id:1, type:1", "id:2, type:2"}
	t.Log(arr)
	jsonArr, err := json.Marshal(arr)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonArr))
	newArr := &[]string{}
	err = json.Unmarshal(jsonArr, newArr)
	if err != nil {
		t.Error(err)
	}
	t.Log(newArr)
}	

// Addtional test is done by postman, of which the testcase is named as test_postman.json,
// and stored under /testcases folder. You can download it to your local desktop and import it to your postman.