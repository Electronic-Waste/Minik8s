package remote_cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"minik8s.io/pkg/network"

	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/oci"
	"minik8s.io/pkg/apis/core"

	"github.com/containerd/containerd"
	"minik8s.io/pkg/constant"
)

// remoteRuntimeService is a gRPC implementation of internalapi.RuntimeService.
type remoteRuntimeService struct {
	timeout       time.Duration
	runtimeClient *containerd.Client
}

func (r *remoteRuntimeService) Client() *containerd.Client {
	return r.runtimeClient
}

func NewRemoteRuntimeService(connectionTimeout time.Duration) (*remoteRuntimeService, error) {
	// build a new cri client
	client, err := containerd.New(constant.Cli_uri)
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

func propagateContainerdLabelsToOCIAnnotations() oci.SpecOpts {
	return func(ctx context.Context, oc oci.Client, c *containers.Container, s *oci.Spec) error {
		return oci.WithAnnotations(c.Labels)(ctx, oc, c, s)
	}
}

// we use a Container Object to start a container with our purpose
func (cli *remoteRuntimeService) StartContainer(ctx context.Context, containerMeta core.Container, Namespace string) error {
	// get image object first and construct the container
	image_getted, err := cli.runtimeClient.GetImage(ctx, containerMeta.Image)
	if err != nil {
		return err
	}
	// create a container
	var processArgs []string
	flag := len(containerMeta.Command) == 0

	// init the command and core args
	for _, cmd := range containerMeta.Command {
		processArgs = append(processArgs, cmd)
	}
	for _, arg := range containerMeta.Args {
		processArgs = append(processArgs, arg)
	}

	var opts []oci.SpecOpts
	var cOpts []containerd.NewContainerOpts

	// the code to prepare the port map
	portMap := make(map[string]string)
	if len(containerMeta.Ports) > 0 {
		portsJSON, err := json.Marshal(containerMeta.Ports)
		if err != nil {
			return err
		}
		portMap["ports"] = string(portsJSON)
	}

	netConfig := network.DefaultNetOpt()
	nameMap := make(map[string]string)
	if Namespace != "" {
		// with shared network namespace
		// format container:<containerid>
		// TODO : add the checking logic here to check for the format
		netConfig.NetworkSlice = []string{Namespace}

		// init the label to use namespace to find all container
		// parse the Name here
		arr := strings.Split(Namespace, ":")
		if len(arr) < 2 {
			return errors.New("wrong namespace format")
		}
		nameMap["minik8s/podName"] = arr[1]
	}
	network_manager := network.ConstructNetworkManager(*(network.New()), netConfig)

	netOpts, netNewContainerOpts, err := network_manager.ContainerNetworkingOpts(ctx, containerMeta.Name)
	if err != nil {
		fmt.Println("err in network setting")
		panic(err)
	}
	fmt.Printf("the length of NetOpts is %d\n", len(netOpts))
	if flag {
		//opts = append(opts,
		//	oci.WithDefaultSpec(),
		//	propagateContainerdLabelsToOCIAnnotations(),
		//)

		cOpts = append(cOpts, containerd.WithImage(image_getted))
		cOpts = append(cOpts, containerd.WithNewSnapshot(containerMeta.Name+"-snapshot", image_getted))
		cOpts = append(cOpts, containerd.WithNewSpec(oci.WithImageConfig(image_getted)))
		opts = append(opts, netOpts...)
		cOpts = append(cOpts, netNewContainerOpts...)

	} else {
		opts = append(opts,
			oci.WithDefaultSpec(),
			oci.WithImageConfig(image_getted),
			oci.WithProcessArgs(processArgs...),
			oci.WithMounts(core.ConvertMounts(containerMeta.Mounts)),
			propagateContainerdLabelsToOCIAnnotations(),
		)
		opts = append(opts, netOpts...)

		cOpts = append(cOpts, containerd.WithImage(image_getted))
		cOpts = append(cOpts, containerd.WithNewSnapshot(containerMeta.Name+"-snapshot", image_getted))
		cOpts = append(cOpts, containerd.WithAdditionalContainerLabels(portMap))
		cOpts = append(cOpts, containerd.WithNewSpec(opts...))
	}
	container, err := cli.runtimeClient.NewContainer(
		ctx,
		containerMeta.Name,
		cOpts...,
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
	// exitStatusC, err := task.Wait(ctx)
	_, err = task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// call start on the task to execute the redis server
	if err := task.Start(ctx); err != nil {
		return err
	}

	// status := <-exitStatusC
	// code, _, err := status.Result()
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("%s exited with status: %d\n", containerMeta.Name, code)

	return nil
}

// input the name of RunSandBox
func (cli *remoteRuntimeService) RunSandBox(name string) error {
	// use cmd to build a pause container
	// run cmd : nerdctl run -d  --name fake_k8s_pod_pause   registry.aliyuncs.com/google_containers/pause:3.9
	cmd := exec.Command("nerdctl", "run", "-d", "--name", name, "--net", "flannel", constant.SandBox_Image)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	return nil
}
