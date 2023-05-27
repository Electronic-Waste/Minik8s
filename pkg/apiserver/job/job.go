package job

import (
	"encoding/json"
	"io/ioutil"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apiserver/etcd"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/clientutil"
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
	clientutil.HttpPlus("Job", job, "")
	resp.WriteHeader(http.StatusOK)
}
