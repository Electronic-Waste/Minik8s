package scheduler

import (
	"minik8s.io/pkg/apis/core"
	"strings"
)

func (s *Scheduler) MatchSchedule(nodes []core.Node, pod core.Pod) (core.Node, bool) {
	if _, ok := pod.ObjectMeta.Labels["node"]; !ok {
		return core.Node{}, false
	}
	for _, node := range nodes {
		if strings.Compare(node.MetaData.Name, pod.ObjectMeta.Labels["node"]) == 0 {
			return node, true
		}
	}
	return core.Node{}, false
}
