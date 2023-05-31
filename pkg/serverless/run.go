package serverless

import(
	"net/http"

	"minik8s.io/pkg/serverless/util/url"
	"minik8s.io/pkg/serverless/register"
	"minik8s.io/pkg/serverless/trigger"
)

type Knative struct {

}

func NewKnative() *Knative {
	return &Knative{}
}

func (knatvie *Knative) Run() {
	http.HandleFunc(url.FunctionRegisterURL, register.HandleFuncRegister)
	http.HandleFunc(url.FunctionTriggerURL, trigger.HandleFuncTrigger)

	// Start Server
	http.ListenAndServe(url.ManagerPort, nil)
}


