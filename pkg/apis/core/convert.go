package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func (p *Pod) getMountPath(name string) (string, error) {
	for _, v := range p.Spec.Volumes {
		if strings.Compare(name, v.Name) == 0 {
			return v.HostPath.Path, nil
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
	viper.SetConfigType("yaml")
	file, err := os.ReadFile(path)
	err = viper.ReadConfig(bytes.NewReader(file))
	if err != nil {
		//fmt.Println("error reading file, please use relative path\n for example: apply ./cmd/config/xxx.yml")
		return nil, err
	}
	pod := Pod{}
	err = viper.Unmarshal(&pod)
	if err != nil {
		return nil, err
	}
	return &pod, nil
}
