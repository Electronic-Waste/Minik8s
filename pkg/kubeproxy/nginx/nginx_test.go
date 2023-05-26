package nginx

import (
	"testing"

	"minik8s.io/pkg/apis/core"
)

func TestInitNginxConf(t *testing.T) {
	controller, err := NewNginxController("/root/minik8s/conf/nginx.conf")
	if err != nil {
		t.Errorf("%v", err)
	}
	controller.InitNginxConf()
}

func TestApplyNginxServer(t *testing.T) {
	controller, err := NewNginxController("/root/minik8s/conf/nginx.conf")
	if err != nil {
		t.Errorf("%v", err)
	}
	hostName := "minik8s.io"
	paths := []core.DNSSubpath{
		{
			Path: "/test1",
			ServiceName: "test-service",
			ClusterIP: "222.222.0.1",
			ServicePort: 80,
		},
		{
			Path: "/test2",
			ServiceName: "test-service",
			ClusterIP: "222.222.0.1",
			ServicePort: 8080,
		},
	}
	controller.ApplyNginxServer(hostName, paths)

}

// func TestDelNginxServer(t *testing.T) {
// 	controller, err := NewNginxController("/root/minik8s/conf/nginx.conf")
// 	if err != nil {
// 		t.Errorf("%v", err)
// 	}
// 	hostName := "minik8s.io"
// 	controller.DelNginxServer(hostName)
// }