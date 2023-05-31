package jobserver

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

var (
	GPUDIR     string = "/mnt/scripts/"
	GPUSH      string = GPUDIR + "gpu.sh"
	GPUTRANSH  string = GPUDIR + "gputran.sh"
	GPUQUEUESH string = GPUDIR + "gpuqueue.sh"
	GPUOUTSH   string = GPUDIR + "gpuout.sh"
	GPUDATA    string = "/mnt/data/"
)

type JobServer struct {
	JobId string
}

func NewJobServer() *JobServer {
	return &JobServer{}
}

func (js *JobServer) Run(isLocal bool, file, scripts, result string) error {
	if isLocal {
		GPUDIR = "/root/minik8s/minik8s/scripts/gpuscripts/"
		GPUSH = GPUDIR + "gpu.sh"
		GPUTRANSH = GPUDIR + "gputran.sh"
		GPUQUEUESH = GPUDIR + "gpuqueue.sh"
		GPUOUTSH = GPUDIR + "gpuout.sh"
		GPUDATA = "/root/minik8s/minik8s/scripts/data/"
	}
	// transform data first
	js.tranData(file, scripts)
	js.batchJob(scripts)
	for {
		if !js.getJobStatus() {
			break
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Job have finished")
	js.receiveResult(result)
	return nil
}

func (js *JobServer) tranData(file string, scripts string) error {
	cmd := exec.Command(GPUTRANSH, GPUDATA+file)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	cmd = exec.Command(GPUTRANSH, GPUDATA+scripts)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	return nil
}

func (js *JobServer) batchJob(scripts string) error {
	cmd := exec.Command(GPUSH, scripts)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	id := outStr[len(outStr)-9 : len(outStr)-1]
	fmt.Printf("id is %s and len is %d\n", id, len(id))
	js.JobId = id
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	return nil
}

func (js *JobServer) getJobStatus() bool {
	cmd := exec.Command(GPUQUEUESH)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return true
	}
	return strings.Contains(outStr, js.JobId)
}

func (js *JobServer) receiveResult(result string) error {
	cmd := exec.Command(GPUOUTSH, result, GPUDATA)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	return nil
}
