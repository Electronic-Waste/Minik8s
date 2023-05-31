package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s.io/pkg/apis/core"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	// "minik8s.io/pkg/podmanager"
	"minik8s.io/pkg/util/listwatch"
	"strings"
)

type Scheduler struct {
	// the counter used in the rr
	rrCount int
}

func GetNewScheduler() *Scheduler {
	return &Scheduler{
		rrCount: 0,
	}
}

func (s *Scheduler) ApplyPodHanlder(msg *redis.Message) {
	var Param core.ScheduleParam
	json.Unmarshal([]byte(msg.Payload), &Param)
	Param.RunPod.ContainerConvert()
	fmt.Printf("Scheduler receive msg: %s", msg.Payload)
	node := s.Schedule(Param.NodeList, Param.RunPod)
	Param.RunPod.Spec.RunningNode = node
	// send back to api-server
	body, err := json.Marshal(Param.RunPod)
	if err != nil {
		fmt.Println(err)
	}
	listwatch.Publish(apiurl.SchedApplyURL, string(body))
}

func (s *Scheduler) UpdatePodHandler(msg *redis.Message) {

}

func (s *Scheduler) BindWatchHandler() {
	go listwatch.Watch("/pods/status/apply", s.ApplyPodHanlder)
	go listwatch.Watch("/pods/status/update", s.UpdatePodHandler)
}

func (s *Scheduler) Run() {
	// Start Server
	stop := make(chan int)
	scheduler := GetNewScheduler()
	scheduler.BindWatchHandler()
	<-stop
}

func (s *Scheduler) Schedule(nodes []core.Node, pod core.Pod) core.Node {
	if node, ok := s.MatchSchedule(nodes, pod); ok {
		fmt.Println("match schedule")
		return node
	} else {
		if _, ok := pod.Labels["resourcepolicy"]; ok && strings.Compare(pod.Labels["resourcepolicy"], "on") == 0 {
			fmt.Println("call memory schedule")
			return s.MemSchedule(nodes, pod)
		}
		return s.RRSchedule(nodes, pod)
	}
}
