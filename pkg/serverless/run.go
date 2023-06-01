package serverless

import(
	"net/http"
	"minik8s.io/pkg/serverless/util/url"
	"context"
	"minik8s.io/pkg/serverless/functioncontroller"
)

type Knative struct {
	rrCount int
}

func NewKnative() *Knative {
	return &Knative{
		rrCount: 0,
	}
}

func (k *Knative) Run() {
	//start controller
	ctx := context.Background()
	functioncontroller,_ := functioncontroller.NewFunctionController()
	go functioncontroller.Run(ctx)


	http.HandleFunc(url.FunctionRegisterURL, k.HandleFuncRegister)
	http.HandleFunc(url.FunctionTriggerURL, k.HandleFuncTrigger)
	http.HandleFunc(url.FunctionUpdateURL, k.HandleFuncUpdate)
	http.HandleFunc(url.FunctionDelURL, k.HandleFuncDel)
	http.HandleFunc(url.WorkflowTriggerURL, k.HandleWorkflowTrigger)

	// Start Server
	go http.ListenAndServe(url.ManagerPort, nil)
	<-ctx.Done()
}


