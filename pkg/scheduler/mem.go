package scheduler

import (
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/clientutil"
	"minik8s.io/pkg/kubelet/config"
	"strconv"
)

func (s *Scheduler) MemSchedule(nodes []core.Node, pod core.Pod) core.Node {
	// get all memory usage from all nodes
	var numList []int
	for _, node := range nodes {
		err, str := clientutil.HttpPlus("Mem", "", url.HttpScheme+node.Spec.NodeIp+config.Port+config.MemoryUrl)
		if err != nil {
			numList = append(numList, 0)
			continue
		}
		num, err := strconv.Atoi(str)
		if err != nil {
			numList = append(numList, 0)
			continue
		}
		numList = append(numList, num)
	}
	MaxMem, MaxIdx := 0, 0
	for idx, val := range numList {
		if val > MaxMem {
			MaxIdx = idx
		}
		MaxMem = val
	}

	return nodes[MaxIdx]
}
