package controller

import (
	//"encoding/json"
	//"context"
	//"minik8s.io/pkg/apis/core"
	//util "minik8s.io/pkg/util/listwatch"
	//"time"
	//"time"
	//"github.com/go-redis/redis/v8"
	//apiurl "minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/podmanager"

	"testing"
)

func TestDeletePod(t *testing.T) {
	podmanager.DelPod("deployment_test-de7342ed")
}