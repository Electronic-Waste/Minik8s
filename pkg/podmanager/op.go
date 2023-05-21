package podmanager

import (
	"bytes"
	"context"
	"fmt"
	"github.com/containerd/containerd/namespaces"
	"log"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/cli/remote_cli"
	"minik8s.io/pkg/idutil/containerwalker"
	"os/exec"
)

// here just finish some operation need by pod running and deleting

// need the image need by the Pod have been pull
func RunPod(pod *core.Pod) error {
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return err
	}
	err = cli.RunSandBox(pod.Name)
	//time.Sleep(time.Second * 10)
	if err != nil {
		return err
	}
	// run core pod's container
	ctx := namespaces.WithNamespace(context.Background(), "default")
	for _, con := range pod.Spec.Containers {
		err = cli.StartContainer(ctx, con, "container:"+pod.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func DelPod(name string) error {
	// find all container labeled with the 'name'
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return err
	}
	walker := &containerwalker.ContainerWalker{
		Client: cli.Client(),
		OnFound: func(ctx context.Context, found containerwalker.Found) error {
			fmt.Println(found.Container.ID())
			stopContainer(found.Container.ID())
			delContainer(found.Container.ID())
			return nil
		},
	}

	// !!! : need to specify the namespace of finding container
	ctx := namespaces.WithNamespace(context.Background(), "default")
	_, err = walker.WalkPod(ctx, name)
	if err != nil {
		return err
	}

	// delete pause container
	stopContainer(name)
	delContainer(name)
	return nil
}

// judge a Pod is running or not
func IsPodRunning(name string) bool {
	// just determine the pause container is running or not
	// find all container labeled with the 'name'
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return err
	}
	is_find := false
	walker := &containerwalker.ContainerWalker{
		Client: cli.Client(),
		OnFound: func(ctx context.Context, found containerwalker.Found) error {
			is_find = true
			return nil
		},
	}

	ctx := namespaces.WithNamespace(context.Background(), "default")
	_, err = walker.Walk(ctx, name)
	return is_find
}

func stopContainer(id string) error {
	// use cmd to build a pause container
	// run cmd : nerdctl run -d  --name fake_k8s_pod_pause   registry.aliyuncs.com/google_containers/pause:3.9
	cmd := exec.Command("nerdctl", "stop", id)
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

func delContainer(id string) error {
	// use cmd to build a pause container
	// run cmd : nerdctl run -d  --name fake_k8s_pod_pause   registry.aliyuncs.com/google_containers/pause:3.9
	cmd := exec.Command("nerdctl", "rm", id)
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
