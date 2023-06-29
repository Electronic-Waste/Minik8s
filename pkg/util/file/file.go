package file

import (
	"fmt"
	"io"
	"minik8s.io/pkg/constant"
	"os"
)

func GenConfigFile() error {
	err := ClearConfig()
	if err != nil {
		return err
	}
	err = genConfig()
	return err
}

func ClearConfig() error {
	err := clear(constant.SysPodDir)
	return err
}

func clear(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(dir + "/" + name)
		if err != nil {
			return err
		}
	}
	return nil
}

func genConfig() error {
	_, err := copyFile(constant.ConfigDir+constant.API, constant.SysPodDir+constant.API)
	if err != nil {
		return err
	}
	_, err = copyFile(constant.ConfigDir+constant.CONTROLLER, constant.SysPodDir+constant.CONTROLLER)
	if err != nil {
		return err
	}
	_, err = copyFile(constant.ConfigDir+constant.SCH, constant.SysPodDir+constant.SCH)
	if err != nil {
		return err
	}
	_, err = copyFile(constant.ConfigDir+constant.KNATIVE, constant.SysPodDir+constant.KNATIVE)
	if err != nil {
		return err
	}
	return nil
}
func copyFile(srcFile, destFile string) (int, error) {
	file1, err := os.Open(srcFile)
	if err != nil {
		return 0, err
	}
	file2, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer file1.Close()
	defer file2.Close()
	bs := make([]byte, 1024, 1024)
	n := -1
	total := 0
	for {
		n, err = file1.Read(bs)
		if err == io.EOF || n == 0 {
			fmt.Println("finish copy")
			break
		} else if err != nil {
			fmt.Println("error")
			return total, err
		}
		total += n
		file2.Write(bs[:n])
	}
	return total, nil
}
