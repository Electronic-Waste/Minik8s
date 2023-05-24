package scheduler

import (
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"net/http"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

var scheduMap = map[string]HttpHandler{
	apiurl.SchedPostURL: HandleSchedu,
}

type Scheduler struct {
	// the counter used in the rr
	rrCount int
}

func GetNewScheduler() *Scheduler {
	return &Scheduler{
		rrCount: 0,
	}
}

func (s *Scheduler) Run() {
	// Start Server
	for url, handler := range scheduMap {
		http.HandleFunc(url, handler)
	}
	http.ListenAndServe(apiurl.SchedulerPort, nil)
}
