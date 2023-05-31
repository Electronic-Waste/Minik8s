package functioncontroller

import(
	"testing"
	"minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apis/meta"
	"encoding/json"
	"time"
	"context"
	"fmt"
)

func TestCountingdown (t *testing.T) {
	fc,_ := NewFunctionController()
	go fc.Run(context.Background())
	time.Sleep(time.Second)
	deployment := core.Deployment{
		Metadata: meta.ObjectMeta{
			Name: "function_test",
		},
	}
	bytes,err := json.Marshal(deployment)
	if err != nil{
		fmt.Println(err)
		return
	}
	listwatch.Publish(FunctionApplyUrl, bytes)
	time.Sleep(5 * time.Second)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	time.Sleep(15 * time.Second)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	time.Sleep(20 * time.Second)
}