package yaml

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// We assume the yaml file have not --- string except the separator
const test_path = "./test_set/test.yaml"

func TestReadYaml(t *testing.T) {
	// test the function to read yaml and parse to a document
	buf, err := os.ReadFile(test_path)
	if err != nil {
		panic(err)
	}

	// construct a new buf
	dou_buf := bytes.NewBuffer(buf)
	reader := NewYAMLReader(bufio.NewReader(dou_buf))
	for {
		getting_buf, err := reader.Read()
		get_str := string(getting_buf)
		if err == io.EOF {
			fmt.Println("the len of getting str is ", len(get_str))
			if len(get_str) != 0 {
				if strings.Contains(get_str, "---") {
					t.Error("parse error in yaml")
				}
			}
			break
		}
		if err != nil {
			t.Error(err)
		}
		if len(get_str) == 0 || strings.Contains(get_str, "---") {
			t.Error("parse error in yaml")
		}

		fmt.Println("parse one yaml document")
		fmt.Println(get_str)
	}

}
