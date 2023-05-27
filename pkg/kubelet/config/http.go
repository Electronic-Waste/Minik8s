package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/podmanager"
	"net/http"
)

// the source of http to apply the Pod

type HttpHandler func(http.ResponseWriter, *http.Request)

var (
	//PodAddWatchUrl string = "/Pod/Metrics"

	Port      string = ":3000"
	PodPrefix string = "/Pod"
	RunPodUrl string = PodPrefix + "/run"
	DelPodRul string = PodPrefix + "/del"
	PodMetricsUrl string = PodPrefix + "/metrics"
)

func HandlePodDel(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	pod := core.Pod{}
	json.Unmarshal(body, &pod)
	fmt.Println("in kubelet http server")
	fmt.Println(pod)
	err := podmanager.DelPod(pod.Name)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
}

func HandlePodRun(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	pod := core.Pod{}
	json.Unmarshal(body, &pod)
	fmt.Println("in kubelet http server")
	fmt.Println(pod)
	err := podmanager.RunPod(&pod)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	fmt.Printf("Run Pod Ip is %s\n", pod.Status.PodIp)
	resMes, err := json.Marshal(pod.Status.PodIp)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.Write(resMes)
}

func Run(PodMap map[string]HttpHandler) error {
	for url, handler := range PodMap {
		http.HandleFunc(url, handler)
	}

	http.ListenAndServe(Port, nil)
	return nil
}
