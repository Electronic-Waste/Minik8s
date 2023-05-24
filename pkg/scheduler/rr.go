package scheduler

import "minik8s.io/pkg/apis/core"

// the package rr scheduler
func (s *Scheduler) RRSchedule(nodes []core.Node, pod core.Pod) core.Node {
	s.rrCount++
	s.rrCount = s.rrCount % len(nodes)
	return nodes[s.rrCount]
}
