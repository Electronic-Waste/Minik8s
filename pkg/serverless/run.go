package serverless

import(
	"net/http"

	"minik8s.io/pkg/serverless/util/url"
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
	http.HandleFunc(url.FunctionRegisterURL, k.HandleFuncRegister)
	http.HandleFunc(url.FunctionTriggerURL, k.HandleFuncTrigger)
	http.HandleFunc(url.WorkflowTriggerURL, k.HandleWorkflowTrigger)

	// Start Server
	http.ListenAndServe(url.ManagerPort, nil)
}


