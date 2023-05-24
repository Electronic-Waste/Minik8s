package kubelet

import (
	"fmt"
	"minik8s.io/pkg/kubelet/config"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/util/listwatch"
	"encoding/json"
	kubetypes "minik8s.io/pkg/kubelet/types"
	"github.com/go-redis/redis/v8"
	"minik8s.io/pkg/podmanager"
	"os"
)

// that is a object that admin the control plane
// Bootstrap is a bootstrapping interface for kubelet, targets the initialization protocol
type Bootstrap interface {
	//GetConfiguration() kubeletconfiginternal.KubeletConfiguration
	//BirthCry()
	//StartGarbageCollection()
	//ListenAndServe()
	//ListenAndServeReadOnly(address net.IP, port uint)
	//ListenAndServePodResources()
	Run(chan kubetypes.PodUpdate)
	//RunOnce(<-chan kubetypes.PodUpdate) ([]RunPodResult, error)
}

type Kubelet struct {
	// TODO(wjl) : add some object need by kubelet to admin the Pod or Deployment
}

func (k *Kubelet) Run(update chan kubetypes.PodUpdate) {
	// wait for new event caused by listening source
	bindWatchHandler()
	k.syncLoop(update)
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
	return nil
}

func NewMainKubelet(podConfig **config.PodConfig) (*Kubelet, error) {
	// return a new Kubelet Object
	*podConfig = makePodSourceConfig()
	return &Kubelet{}, nil
}

func makePodSourceConfig() *config.PodConfig {
	// TODO(wjl) : add fileSource support here
	cfg := config.NewPodConfig()
	config.NewSourceFile(cfg.Channel(kubetypes.FileSource))
	return cfg
}

func bindWatchHandler() {
	go listwatch.Watch("/pods/status/apply", ApplyPodHanlder)
	go listwatch.Watch("/pods/status/del", DeletePodHandler)
	go listwatch.Watch("/pods/status/update", UpdatePodHandler)
}

func ApplyPodHanlder(msg *redis.Message) {
	var podSpec core.Pod
	json.Unmarshal([]byte(msg.Payload), &podSpec)
	podSpec.ContainerConvert()
	fmt.Printf("Kubelet receive msg: %s", msg.Payload)
	podmanager.RunPod(&podSpec)
}

func UpdatePodHandler(msg *redis.Message) {

}

func DeletePodHandler(msg *redis.Message) {
	podName := msg.Payload
	fmt.Printf("kubelet receive del msg: %s", podName)
	podmanager.DelPod(podName)
}
