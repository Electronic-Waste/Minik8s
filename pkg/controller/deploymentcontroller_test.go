package controller

import(
	"minik8s.io/pkg/podmanager"

	"testing"
)

func TestDeletePod(t *testing.T) {
	podmanager.DelPod("deployment_test-de7342ed")
}