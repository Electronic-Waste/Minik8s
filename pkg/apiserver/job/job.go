package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apiserver/etcd"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/controller"
	"net/http"
	"path"
)

func HandleApplyJob(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	job := core.Job{}
	json.Unmarshal(body, &job)
	jobName := job.Meta.Name
	etcdURL := path.Join(apiurl.JobApplyUrl, jobName)
	err := etcd.Put(etcdURL, string(body))
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// start a New Pod to commit the Job
	// send jc listening on 9000 port
	err, str := clientutil.HttpPlus("Job", job, apiurl.HttpScheme+apiurl.HostURL+controller.JCPORT+controller.RunJobUrl)
	fmt.Printf("receive output is %s\n", str)
	resp.WriteHeader(http.StatusOK)
}

func HandleGetJob(resp http.ResponseWriter, req *http.Request) {
	jobs, err := etcd.GetWithPrefix(controller.JOBMAP)
	actJob := core.Job{}
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	maps := core.JobMaps{}
	for idx, _ := range jobs {
		json.Unmarshal([]byte(jobs[idx]), &actJob)
		maps.Maps = append(maps.Maps, core.Job2Pod{
			JobName: actJob.Meta.Name,
			PodName: actJob.Status.PodName,
		})
	}
	jsonVal, err := json.Marshal(maps)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonVal)
}

func HandleMapJob(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	job := core.Job{}
	json.Unmarshal(body, &job)
	fmt.Printf("get pob name is %s\n", job.Status.PodName)
	jobName := job.Meta.Name
	etcdURL := path.Join(apiurl.JobApplyUrl, jobName)
	err := etcd.Put(etcdURL, string(body))
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)

}
