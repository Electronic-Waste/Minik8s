package controller

import(
	"fmt"
	"testing"
)

func TestHttpGet(t *testing.T) {
	pods,_ := GetReplicaPods("deployment_test")
	fmt.Println(pods)
}