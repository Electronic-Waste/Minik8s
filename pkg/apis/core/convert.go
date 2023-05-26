package core

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
	"strings"
)

func (p *Pod) getMountPath(name string) (string, error) {
	for _, v := range p.Spec.Volumes {
		if strings.Compare(name, v.Name) == 0 {
			return v.HostPath, nil
		}
	}
	return "", errors.New("not such volume")
}

func (p *Pod) ContainerConvert() error {
	// convert volume to mount
	for idx, con := range p.Spec.Containers {
		if len(con.VolumeMounts) != 0 {
			// need to convert
			for _, vo := range con.VolumeMounts {
				path, err := p.getMountPath(vo.Name)
				if err != nil {
					return err
				}
				p.Spec.Containers[idx].Mounts = append(p.Spec.Containers[idx].Mounts, Mount{
					SourcePath:      path,
					DestinationPath: vo.MountPath,
				})
			}
		}
	}
	return nil
}

func ParsePod(path string) (*Pod, error) {
	if !strings.HasSuffix(path, ".yaml") {
		//get yaml file content
		fmt.Println("error file type")
		return nil, errors.New("error file type")
	}
	file, err := os.ReadFile(path)
	pod := Pod{}
	err = yaml.Unmarshal(file, &pod)
	fmt.Println(pod)
	pod.ContainerConvert()
	fmt.Printf("pod name after parse is %s\n", pod.Name)
	if err != nil {
		return nil, err
	}
	return &pod, nil
}

func ParseNode(path string) (*Node, error) {
	if !strings.HasSuffix(path, ".yaml") {
		//get yaml file content
		fmt.Println("error file type")
		return nil, errors.New("error file type")
	}
	file, err := os.ReadFile(path)
	node := Node{}
	err = yaml.Unmarshal(file, &node)
	fmt.Println(node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}
