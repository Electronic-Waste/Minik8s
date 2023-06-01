package serverless

import(
	"net/http"
	//"context"
	"minik8s.io/pkg/serverless/util/url"
	//"minik8s.io/pkg/serverless/functioncontroller"
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
	//ctx := context.Background()
	//start controller
	//functioncontroller,_ := functioncontroller.NewFunctionController()
	//go functioncontroller.Run(ctx)


	http.HandleFunc(url.FunctionRegisterURL, k.HandleFuncRegister)
	http.HandleFunc(url.FunctionTriggerURL, k.HandleFuncTrigger)
	// Start Server
	http.ListenAndServe(url.ManagerPort, nil)
	//<-ctx.Done()
}


