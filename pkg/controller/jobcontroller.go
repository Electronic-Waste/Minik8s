package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"minik8s.io/pkg/apis/core"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/util/tools/uid"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

const (
	JCPORT    string = ":9000"
	RunJobUrl string = "/JCRUN"
	JOBMAP    string = "/JOB/"
)

type JobController struct {
	FileCount int
}

func NewJobController() (*JobController, error) {
	return &JobController{
		FileCount: 0,
	}, nil
}

func (jc *JobController) Run(ctx context.Context) {
	fmt.Println("jc running")
	go jc.RunHttp()
	<-ctx.Done()
	return
}

func (jc *JobController) RunHttp() {
	var postHandlerMap = map[string]HttpHandler{
		RunJobUrl: jc.HandleRunJob,
	}
	// Bind POST request with handler
	for url, handler := range postHandlerMap {
		http.HandleFunc(url, handler)
	}
	// Start Server
	http.ListenAndServe(JCPORT, nil)
}

func (jc *JobController) HandleRunJob(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	job := core.Job{}
	json.Unmarshal(body, &job)
	fmt.Printf("get job is \n%v\n", job)
	dir, filePath, codeFile, err := jc.genConfig(job)
	if err != nil {
		fmt.Println(err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// config the Pod message
	var pod core.Pod
	pod.Name = "job-" + uid.NewUid()
	pod.Labels = map[string]string{"node": "vmeet2"}
	job.Status.PodName = pod.Name
	// send back to api-server
	err, _ = clientutil.HttpPlus("Job", job, apiurl.Prefix+apiurl.JobMapUrl)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	pod.Kind = "Pod"
	pod.Spec.Volumes = []core.Volume{{
		Name:     "shared-data",
		HostPath: dir,
	},
		{
			Name:     "shared-scripts",
			HostPath: "/root/minik8s/minik8s/scripts/gpuscripts",
		}}
	pod.Spec.Containers = []core.Container{
		{
			Name:  "t1" + uid.NewUid(),
			Image: "docker.io/library/jobserver:latest",
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "shared-data",
					MountPath: "/mnt/data",
				},
				{
					Name:      "shared-scripts",
					MountPath: "/mnt/scripts",
				},
			},
			Ports:   []core.ContainerPort{},
			Command: []string{"/mnt/scripts/jobserver"},
			Args:    []string{"remote", "--file=" + codeFile, "--scripts=" + filePath, "--result=" + job.Spec.FileName},
		},
	}
	err = pod.ContainerConvert()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	fmt.Println(pod)
	// send message to apiserver to apply a Pod
	err = clientutil.HttpApply("Pod", pod)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func parseCodePath(path string) (dir string, file string, err error) {
	idx := strings.LastIndex(path, "/")
	if idx < 0 {
		dir = ""
		file = ""
		err = errors.New("file format error")
		return
	}
	dir = path[:idx]
	file = path[idx+1:]
	err = nil
	return
}

func genSlurm(job core.Job) ([]string, string, string, error) {
	var fileThing []string
	dir, file, err := parseCodePath(job.Spec.CodePath)
	if err != nil {
		return fileThing, "", "", err
	}
	fileThing = append(fileThing, "#!/bin/bash\n\n")
	// not consider error handling
	fileThing = append(fileThing, "#SBATCH --job-name="+job.Meta.Name+"\n")
	fileThing = append(fileThing, "#SBATCH --partition="+job.Spec.Partition+"\n")
	fileThing = append(fileThing, "#SBATCH -N "+strconv.Itoa(job.Spec.ThreadNum)+"\n")
	fileThing = append(fileThing, "#SBATCH --ntasks-per-node="+strconv.Itoa(job.Spec.TaskPerNode)+"\n")
	fileThing = append(fileThing, "#SBATCH --cpus-per-task="+strconv.Itoa(job.Spec.CPUPerTask)+"\n")
	fileThing = append(fileThing, "#SBATCH --gres=gpu:"+strconv.Itoa(job.Spec.GPUNum)+"\n")
	fileThing = append(fileThing, "#SBATCH --output="+job.Spec.FileName+".out\n")
	fileThing = append(fileThing, "#SBATCH --error="+job.Spec.FileName+".out\n")
	fileThing = append(fileThing, "module load gcc/8.3.0 cuda/10.2\n")
	fileThing = append(fileThing, "nvcc "+file+" -o test -lcublas\n")
	fileThing = append(fileThing, "./test")
	return fileThing, dir, file, nil
}

func (jc *JobController) genConfig(job core.Job) (string, string, string, error) {
	// generate the slurm scripts
	filePath, err := jc.genFileName()
	if err != nil {
		return "", "", "", err
	}
	fileThing, dir, codeFile, err := genSlurm(job)
	file, err := os.OpenFile(dir+"/"+filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	if err != nil {
		return "", "", "", err
	}
	for _, str := range fileThing {
		writer.WriteString(str)
	}
	writer.Flush()
	return dir, filePath, codeFile, nil
}

func (jc *JobController) genFileName() (string, error) {
	jc.FileCount++
	if jc.FileCount < 0 {
		return "", errors.New("int overflow\n")
	}
	return "test" + strconv.Itoa(jc.FileCount) + ".slurm", nil
}
