package remote_cli

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	"minik8s.io/pkg/apis/core"
	"time"

	"github.com/containerd/containerd"
	constant "minik8s.io/pkg/const"
)

// remoteRuntimeService is a gRPC implementation of internalapi.RuntimeService.
type remoteRuntimeService struct {
	timeout       time.Duration
	runtimeClient *containerd.Client
}

func NewRemoteRuntimeService(connectionTimeout time.Duration) (*remoteRuntimeService, error) {
	// build a new cri client
	client, err := containerd.New(constant.Cri_uri)
	// need to call client.Close() to gc this object
	if err != nil {
		return nil, err
	}
	return &remoteRuntimeService{
		connectionTimeout,
		client,
	}, nil
}

func NewRemoteImageServiceByImageService(cli *remoteImageService) *remoteRuntimeService {
	return &remoteRuntimeService{
		timeout:       cli.timeout,
		runtimeClient: cli.imageClient,
	}
}

// set filter to nil and list all containers
func (cli *remoteRuntimeService) ListContainers(ctx context.Context, filters ...string) ([]containerd.Container, error) {
	res, err := cli.runtimeClient.Containers(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// we use a Container Object to start a container with our purpose
func (cli *remoteRuntimeService) StartContainer(ctx context.Context, containerMeta core.Container) error {
	// get image object first and construct the container
	// image_getted, err := cli.runtimeClient.GetImage(ctx, image)
	// if err != nil {
	// 	return nil, err
	// }
	image_getted, err := cli.runtimeClient.GetImage(ctx, containerMeta.Image)
	if err != nil {
		return err
	}
	// create a container
	var processArgs []string
	for _, cmd := range containerMeta.Command {
		processArgs = append(processArgs, cmd)
	}
	for _, arg := range containerMeta.Args {
		processArgs = append(processArgs, arg)
	}
	var s specs.Spec
	var opts []oci.SpecOpts
	opts = append(opts,
		oci.WithDefaultSpec(),
		oci.WithProcessArgs(processArgs...),
	)
	container, err := cli.runtimeClient.NewContainer(
		ctx,
		containerMeta.Name,
		containerd.WithImage(image_getted),
		containerd.WithNewSnapshot(containerMeta.Name+"-snapshot", image_getted),
		containerd.WithNewSpec(oci.WithImageConfig(image_getted)),
		containerd.WithSpec(&s, opts...),
	)
	if err != nil {
		return err
	}

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}

	defer task.Delete(ctx)

	// make sure we wait before calling start
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// call start on the task to execute the redis server
	if err := task.Start(ctx); err != nil {
		return err
	}
	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	fmt.Printf("%s exited with status: %d\n", containerMeta.Name, code)

	return nil
}
