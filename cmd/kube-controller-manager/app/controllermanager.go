package app

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

func NewControllerManagerCommand() *cobra.Command {
	//opts, err := options.NewKubeControllerManagerOptions()
	cmd := &cobra.Command{
		Use: "kube-controller-manager",
		Run: func(cmd *cobra.Command, args []string) {
			//c := opts.Config()
			Run(context.Background())
		},
	}
	print("new controller manager\n")
	return cmd
}

func Run(ctx context.Context) error {

	//controllerContext, err := CreateControllerContext(c)
	//if err != nil {
	//	return err
	//}
	if err := StartControllers(ctx, NewControllerInitializers()); err != nil {
		fmt.Printf("error starting controllers: %v\n", err)
	}
	return nil
}

// InitFunc is used to launch a particular controller. It returns a controller
// that can optionally implement other interfaces so that the controller manager
// can support the requested features.
// The returned controller may be nil, which will be considered an anonymous controller
// that requests no additional features from the controller manager.
// Any error returned will cause the controller process to `Fatal`
// The bool indicates whether the controller was enabled.
type InitFunc func(ctx context.Context) (err error)

func NewControllerInitializers() map[string]InitFunc {
	controller := map[string]InitFunc{}
	controller["deployment"] = StartDeploymentController
	controller["autoscaler"] = StartAutoSclaerController
	controller["job"] = StartJobController
	return controller
}

func StartControllers(ctx context.Context, controllers map[string]InitFunc) error {
	for _, initFunc := range controllers {
		err := initFunc(ctx)
		if err != nil {
			return err
		}
	}
	<-ctx.Done()
	return nil
}
