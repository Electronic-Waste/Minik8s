package kubelet

import (
	"fmt"
	"minik8s.io/pkg/kubelet/config"
	// "minik8s.io/pkg/apis/core"
	// "minik8s.io/pkg/util/listwatch"
	"encoding/json"
	"minik8s.io/pkg/kubelet/cadvisor"
	kubetypes "minik8s.io/pkg/kubelet/types"
	"minik8s.io/pkg/kubeproxy"
	"net/http"
	"os"
	"time"
	"github.com/go-redis/redis/v8"
	"minik8s.io/pkg/util/listwatch"
	apiurl "minik8s.io/pkg/apiserver/util/url"
)

// that is a object that admin the control plane
// Bootstrap is a bootstrapping interface for kubelet, targets the initialization protocol
type Bootstrap interface {
	Run(chan kubetypes.PodUpdate)
}

type Kubelet struct {
	// TODO : add some object need by kubelet to admin the Pod or Deployment

	Cadvisor *cadvisor.CAdvisor
	kubeProxyManager *kubeproxy.KubeproxyManager

}
func (k *Kubelet) Run(update chan kubetypes.PodUpdate) {
	// wait for new event caused by listening source

	k.kubeProxyManager, _ = kubeproxy.NewKubeProxy()
	k.kubeProxyManager.Run()

	//bindWatchHandler()
	PodMap := map[string]config.HttpHandler{
		config.RunPodUrl: 		config.HandlePodRun,
		config.DelPodRul: 		config.HandlePodDel,
	    config.PodMetricsUrl:	k.HandlePodGetMetrics,
		config.GetAllPodUrl:	config.HandleGetAllPod,
		config.MemoryUrl:     config.HandleMemGet,

	}
	go listwatch.Watch(apiurl.PodStatusRegisterMetricsUrl, k.PodRegister)
	go listwatch.Watch(apiurl.PodStatusUnregisterMetricsUrl, k.PodUnregister)
	//go k.PodRegister()
	go config.Run(PodMap)
	k.syncLoop(update)
}

func (k *Kubelet) HandlePodGetMetrics(resp http.ResponseWriter, req *http.Request) {
	//fmt.Println("kubelet get pod metrics")
	vars := req.URL.Query()
	podName := vars.Get("name")
	fmt.Println("get pod: ", podName)
	stats,err := k.Cadvisor.GetPodMetric(podName)
	fmt.Println(stats)
	if err != nil{
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	data, err := json.Marshal(stats)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(data))
}

func (k *Kubelet) PodRegister (msg *redis.Message) {
	//fmt.Println("pod register")
	time.Sleep(time.Millisecond * 500)
	bytes := []byte(msg.Payload)
	var podname string
	err := json.Unmarshal(bytes, &podname)
	fmt.Println("register pod: ",podname)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = k.Cadvisor.RegisterPod(podname)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (k *Kubelet) PodUnregister (msg *redis.Message) {
	//fmt.Println("pod unregister")
	time.Sleep(time.Millisecond * 500)
	bytes := []byte(msg.Payload)
	var podname string
	err := json.Unmarshal(bytes, &podname)
	fmt.Println("unregister pod: ",podname)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = k.Cadvisor.UnRegisterPod(podname)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (k *Kubelet) syncLoop(update chan kubetypes.PodUpdate) {
	for {
		if err := k.syncLoopIteration(update); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func (k *Kubelet) syncLoopIteration(update chan kubetypes.PodUpdate) error {
	// add the logic to receive the message from channel and deal with it
	time.Sleep(1 * time.Second)
	return nil
}

func NewMainKubelet(podConfig **config.PodConfig) (*Kubelet, error) {
	// return a new Kubelet Object
	*podConfig = makePodSourceConfig()
	return &Kubelet{
		Cadvisor: cadvisor.GetCAdvisor(),
	}, nil
}

func makePodSourceConfig() *config.PodConfig {
	// TODO(wjl) : add fileSource support here
	cfg := config.NewPodConfig()
	config.NewSourceFile(cfg.Channel(kubetypes.FileSource))
	return cfg
}