package nginx

import (
	"os"

	"minik8s.io/pkg/apis/core"
)

type NginxFunInterface interface {
	// InitNginxConf creates initial nginx.conf
	InitNginxConf() error

	// ApplyNginxServer adds new server block in nginx.conf
	ApplyNginxServer(hostName string, paths []core.DNSSubpath) error

	// DeleteNginxServer dels a server block in nginx.conf by hostName
	DelNginxServer(hostName string) error

	// ReloadNginxSerice apply new nginx config into use
	ReloadNginxService()  error
	
}

type NginxController struct {
	NginxConfPath string
}

func NewNginxController(nginxConfPath string) (*NginxController, error) {
	// Check if file exists
	_, err := os.Stat(nginxConfPath)
	if err != nil {
		return nil, err
	}
	return &NginxController{
		NginxConfPath : nginxConfPath,
	}, nil
}

func (controller *NginxController) InitNginxConf() error {
	lines := []string{}
	lines = append(lines, "worker_processes 1\n")
	lines = append(lines, "events {\nworker_connections 1024;\n}\n")
}

func (controller *NginxController) ApplyNginxServer(hostName string, paths []core.DNSSubpath) error {

}

func (controller *NginxController) DelNginxServer(hostName string) error {

}

func (controller *NginxController)ReloadNginxService() error {

}