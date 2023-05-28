package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"minik8s.io/pkg/apis/core"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	// "minik8s.io/pkg/podmanager"
	"minik8s.io/pkg/util/listwatch"
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
	fmt.Printf("Scheduler receive msg: %s\n", msg.Payload)
	node := s.RRSchedule(Param.NodeList, Param.RunPod)
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

// func (s *Scheduler) DeletePodHandler(msg *redis.Message) {
// 	podName := msg.Payload
// 	fmt.Printf("kubelet receive del msg: %s", podName)
// 	podmanager.DelPod(podName)
// }

func (s *Scheduler) BindWatchHandler() {
	go listwatch.Watch("/pods/status/apply", s.ApplyPodHanlder)
	// go listwatch.Watch("/pods/status/del", s.DeletePodHandler)
	go listwatch.Watch("/pods/status/update", s.UpdatePodHandler)
}

func (s *Scheduler) Run() {
	// Start Server
	stop := make(chan int)
	scheduler := GetNewScheduler()
	scheduler.BindWatchHandler()
	<-stop
}