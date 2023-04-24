package app

import (
	"context"
	"fmt"
	"k8s/cmd/kube-controller-manager/app/config"
	"k8s/cmd/kube-controller-manager/app/controllercontext"
	"k8s/cmd/kube-controller-manager/app/options"

	"github.com/spf13/cobra"
)

// InitFunc is used to launch a particular controller.
// Any error returned will cause the controller process to `Fatal`
type InitFunc func(ctx context.Context, controllerCtx controllercontext.ControllerContext) (err error)

// a tiny case of run
// actually config makes little difference to options
// TODO: for now, there is no client created, which is the core of controller
// TODO: add flags
func NewControllerManagerCommand() *cobra.Command {
	opts, err := options.NewKubeControllerManagerOptions()
	cmd := &cobra.Command{
		Use: "kube-controller-manager",
		Run: func(cmd *cobra.Command, args []string) {
			c := opts.Config()
			Run(context.Background(), c.Complete())
		},
	}
	return cmd
}

// context is not used
func Run(ctx context.Context, c *config.CompletedConfig) error {

	controllerContext, err := CreateControllerContext(c)
	if err != nil {
		return err
	}
	if err := StartControllers(ctx, controllerContext, NewControllerInitializers()); err != nil {
		fmt.Printf("error starting controllers: %v\n", err)
	}
	// TODO: give each controller a new unique ls
	select {}
}

// not completed yet
func CreateControllerContext(c *config.CompletedConfig) (*controllercontext.ControllerContext, error) {

	//sharedInformers := informers.NewSharedInformerFactory(versionedClient, ResyncPeriod(s)())
	//ls, err := listerwatcher.NewListerWatcher(listerwatcher.DefaultConfig())
	//if err != nil {
	//	return nil, err
	//}
	controllerContext := &context.ControllerContext{Config: c}
	return controllerContext, nil
}

func NewControllerInitializers() map[string]InitFunc {
	controller := map[string]InitFunc{}
	// TODO : Initialize the map with controller name and InitFunc
	controller["replicaset"] = StartReplicaSetController
	controller["deployment"] = StartDeploymentController
	return controller
}

func StartControllers(ctx context.Context, controllerContext *controllercontext.ControllerContext, controllers map[string]InitFunc) error {
	for controllerName, initFunc := range controllers {
		err := initFunc(ctx, *controllerContext)
		if err != nil {
			return err
		}
	}
	return nil
}
