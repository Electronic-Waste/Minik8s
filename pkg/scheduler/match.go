package scheduler

import (
	"minik8s.io/pkg/apis/core"
	"strings"
	"fmt"
)

func (s *Scheduler) MatchSchedule(nodes []core.Node, pod core.Pod) (core.Node, bool) {
	fmt.Println("match scheduler")
	if _, ok := pod.ObjectMeta.Labels["node"]; !ok {
		fmt.Println("match scheduler: labels not exist")
		return core.Node{}, false
	}
	for _, node := range nodes {
		if strings.Compare(node.MetaData.Name, pod.ObjectMeta.Labels["node"]) == 0 {
			fmt.Println("match scheduler: match!")
			return node, true
		}
	}
	fmt.Println("match scheduler: node not found")
	return core.Node{}, false
}
