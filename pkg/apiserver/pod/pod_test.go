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
