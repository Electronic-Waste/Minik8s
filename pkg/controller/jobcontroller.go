package controller

import (
	"context"
	"fmt"
	"minik8s.io/pkg/controller/jobutil"
	"net/http"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

const (
	JCPORT    string = "9000"
	RunJobUrl string = "/JCRUN"
)

var postHandlerMap = map[string]HttpHandler{
	RunJobUrl: jobutil.HandleRunJob,
}

type JobController struct {
}

func NewJobController() (*JobController, error) {
	return &JobController{}, nil
}

func (jc *JobController) Run(ctx context.Context) {
	fmt.Println("jc running")
	go jc.RunHttp()
	<-ctx.Done()
	return
}

func (jc *JobController) RunHttp() {
	// Bind POST request with handler
	for url, handler := range postHandlerMap {
		http.HandleFunc(url, handler)
	}
	// Start Server
	http.ListenAndServe(JCPORT, nil)
}
