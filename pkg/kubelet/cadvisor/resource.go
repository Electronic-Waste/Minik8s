package cadvisor

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetFreeMem() (int, error) {
	cmd := exec.Command("free", "-m")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, _ := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return -1, err
	}
	mesList := strings.Split(outStr, "\n")
	if len(mesList) != 4 {
		fmt.Println("output format error")
		return -1, errors.New("output format error")
	}
	numListWithZero := strings.Split(mesList[1], " ")
	var numList []string
	for _, str := range numListWithZero {
		if len(str) != 0 {
			numList = append(numList, str)
		}
	}

	free, err := strconv.Atoi(numList[3])
	if err != nil {
		return -1, err
	}
	available, err := strconv.Atoi(numList[6])
	if err != nil {
		return -1, err
	}
	return free + available, nil
}
