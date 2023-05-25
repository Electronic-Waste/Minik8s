package nginx

import (
	"os"
	"strings"
	"io/ioutil"
	"os/exec"
	"fmt"

	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/kubeproxy/path"
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
	lines = append(lines, "worker_processes 1")
	lines = append(lines, "events {\n\tworker_connections 1024;\n}")
	lines = append(lines, "http {\n}")
	output := strings.Join(lines, "\n")
	return ioutil.WriteFile(controller.NginxConfPath, []byte(output), 0644)

}

func (controller *NginxController) ApplyNginxServer(hostName string, paths []core.DNSSubpath) error {
	// Read original nginx.conf file
	content, err := ioutil.ReadFile(controller.NginxConfPath)
	if err != nil {
		return fmt.Errorf("Error in ApplyNginxServer: %v", err)
	}

	// Parse file content to []stiring & Find inserting position
	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines) - 1] == "" {
		lines = lines[:len(lines) - 1]
	}
	begin_pos := len(lines) - 1
	for ; begin_pos > 0; begin_pos-- {
		if lines[begin_pos] == "}" {
			break
		}
	}

	// Generate new lines
	newLines := []string{}
	newLines = append(newLines, "\tserver {")
	newLines = append(newLines, "\t\tlisten\t80;")
	newLines = append(newLines, fmt.Sprintf("\t\tserver_name %s", hostName))
	for _, subpath := range paths {
		newLines = append(newLines, fmt.Sprintf("\t\tlocation %s {", subpath.Path))
		newLines = append(newLines, fmt.Sprintf("\t\t\t\tproxy_pass http://%s:%s;", subpath.ClusterIP, subpath.ServicePort))
		newLines = append(newLines, "\t\t}")
	}
	newLines = append(newLines, "\t}")
	newLines = append(newLines, "}")

	// Combine them and write to file
	lines = append(lines[:begin_pos], newLines...)
	output := strings.Join(lines, "\n")
	return ioutil.WriteFile(controller.NginxConfPath, []byte(output), 0644)
}

func (controller *NginxController) DelNginxServer(hostName string) error {
	// Read original nginx.conf file
	content, err := ioutil.ReadFile(controller.NginxConfPath)
	if err != nil {
		return fmt.Errorf("Error in DelNginxServer: %v", err)
	}

	// Parse file content to []stiring & Find target server block
	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines) - 1] == "" {
		lines = lines[:len(lines) - 1]
	}
	targetLine := fmt.Sprintf("\t\tserver_name %s", hostName)
	targetLineNum := 0
	for i := range lines {
		if lines[i] == targetLine {
			targetLineNum = i
			break
		}
	}
	if targetLineNum == 0 {
		return fmt.Errorf("Error in DelNginxServer: can't find target server")
	}
	var startLineNum, endLineNum int
	for i := targetLineNum; i > 0; i-- {
		if lines[i] == "\tserver {" {
			startLineNum = i
			break
		}
	}
	for i := targetLineNum; i < len(lines); i++ {
		if lines[i] == "\t}" {
			endLineNum = i
			break
		}
	}

	// Delete lines in [startLineNum, endLineNum] & Write to file
	lines = append(lines[:startLineNum], lines[endLineNum + 1:]...)
	output := strings.Join(lines, "\n")
	return ioutil.WriteFile(controller.NginxConfPath, []byte(output), 0644)
}

func (controller *NginxController)ReloadNginxService() error {
	// `nginx -s reload`
	cmd := exec.Command(path.NginxExecutableFile, path.NginxActionFlag, path.NginxReloadAction)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error in executing `nginx -s reload`: %v", err)
	}
	fmt.Printf("nginx reload: %s", string(output))
	return nil
}