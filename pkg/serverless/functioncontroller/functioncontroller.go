package functioncontroller

import(
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apis/meta"
	"context"
	"minik8s.io/pkg/util/listwatch"
)

var(
	FunctionStatus	 string = "/function/status"
	FunctionApplyUrl string = FunctionStatus + "/apply"
)

type Function{
	//mock
	Metadata meta.ObjectMeta
	Pods []core.Pod
}

type FunctionController struct{
	functionMap	map[string]Function	//record function name to function
}

func NewFunctionController() (*FunctionController,error) {
	fc := &FunctionController{
		functionMap: make(map[string]Function)
	}
}

func (fc *FunctionController) Run (ctx context.Context) {
	go fc.worker(ctx)
	go fc.register()
	<-ctx.Done()
}

func (fc *FunctionController) register {
	go listwatch.Watch(FunctionApplyUrl, )
}

func (fc *FunctionController) listener (msg *redis.Message) {
	
}