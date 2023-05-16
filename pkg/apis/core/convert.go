package core

import (
	"errors"
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
