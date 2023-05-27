package scheduler

import (
	"fmt"

	"minik8s.io/pkg/apis/core"
)

// the package rr scheduler
func (s *Scheduler) RRSchedule(nodes []core.Node, pod core.Pod) core.Node {
	s.rrCount++
	s.rrCount = s.rrCount % len(nodes)
	fmt.Printf("rrCount is %d\n", s.rrCount)
	for _, node := range nodes {
		fmt.Printf("node mes is %s\n", node.MetaData.Name)
	}
	return nodes[s.rrCount]
}
