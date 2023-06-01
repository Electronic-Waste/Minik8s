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
	listwatch.Publish(FunctionTriggerUrl, bytes)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	time.Sleep(15 * time.Second)
	listwatch.Publish(FunctionTriggerUrl, bytes)
	time.Sleep(20 * time.Second)
}