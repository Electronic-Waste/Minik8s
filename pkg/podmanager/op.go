package podmanager

import (
	"bytes"
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/runtime/restart"
	"github.com/docker/go-units"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/cli/remote_cli"
	"minik8s.io/pkg/idutil/containerwalker"
	"os/exec"
	"strings"
	"time"
	"minik8s.io/pkg/util/listwatch"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"encoding/json"
)

// here just finish some operation need by pod running and deleting

// need the image need by the Pod have been pull
func RunPod(pod *core.Pod) error {
	fmt.Println("podmanager run pod")
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
	genPodIp(pod)

	bytes,_ := json.Marshal(pod.Name)
	listwatch.Publish(apiurl.PodStatusGetMetricsUrl, bytes)

	return nil
}

func genPodIp(pod *core.Pod) error {
	// use cmd to generate a Ip for a Pod
	// run cmd : nerdctl inspect -f '{{.NetworkSettings.IPAddress}}' test
	cmd := exec.Command("nerdctl", "inspect", "-f", "`{{.NetworkSettings.IPAddress}}`", pod.Name)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	ip := outStr[1 : len(outStr)-2]
	pod.Status.PodIp = ip
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	return nil
}

func DelSimpleContainer(name string) error {
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
	return nil
}

func DelPod(name string) error {
	err := DelSimpleContainer(name)
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
		return false
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

func IsCrashContainer(name string) bool {
	// just determine the pause container is running or not
	// find all container labeled with the 'name'
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return false
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
	_, err = walker.WalkPod(ctx, name)
	return is_find
}

func GetPods() ([]core.Pod, error) {
	// get pod name and core pod status is ok
	var podSet []core.Pod
	// get all pause container
	// find all container labeled with the 'name'
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	// !!! : need to specify the namespace of finding container
	ctx := namespaces.WithNamespace(context.Background(), "default")
	if err != nil {
		return nil, err
	}
	walker := &containerwalker.ContainerWalker{
		Client: cli.Client(),
		OnFound: func(ctx context.Context, found containerwalker.Found) error {
			//fmt.Println("find a Pod")
			labels, err := found.Container.Labels(ctx)
			if err != nil {
				return err
			}
			status, exitStatus := GetContainerStatus(ctx, found.Container)
			//fmt.Printf("container %s's status is %s and exit status is %d\n", labels["nerdctl/name"], status, exitStatus)
			pod := core.Pod{}
			pod.Name = labels["nerdctl/name"]
			if strings.Contains(status, "Exited") {
				if exitStatus != 0 {
					pod.Status.Phase = core.PodFailed
				} else {
					pod.Status.Phase = core.PodUnknown
				}
			} else {
				pod.Status.Phase = core.PodUnknown
			}
			podSet = append(podSet, pod)
			return nil
		},
	}
	_, err = walker.WalkPause(ctx, "pause")

	// check all container status
	is_running := false
	is_pending := false
	for i, pod := range podSet {
		// walk for Pod's container
		is_running = false
		walker = &containerwalker.ContainerWalker{
			Client: cli.Client(),
			OnFound: func(ctx context.Context, found containerwalker.Found) error {
				//fmt.Println("find a Container")
				if err != nil {
					return err
				}
				status, exitStatus := GetContainerStatus(ctx, found.Container)
				//fmt.Printf("container %s's status is %s and exit status is %d\n", found.Container.ID(), status, exitStatus)
				if strings.Contains(status, "Exited") {
					if exitStatus != 0 {
						pod.Status.Phase = core.PodFailed
					} else {
						pod.Status.Phase = core.PodUnknown
					}
				} else {
					pod.Status.Phase = core.PodUnknown
				}
				if strings.Compare(status, "Up") == 0 {
					is_running = true
				} else if strings.Compare(status, "Create") == 0 {
					is_pending = true
				}
				return nil
			},
		}

		_, err := walker.WalkPod(ctx, pod.Name)
		if err != nil {
			return nil, err
		}
		if is_pending {
			pod.Status.Phase = core.PodPending
		}
		if is_running && pod.Status.Phase == core.PodUnknown {
			pod.Status.Phase = core.PodRunning
		}
		if !is_running && pod.Status.Phase == core.PodUnknown {
			pod.Status.Phase = core.PodSucceeded
		}
		podSet[i] = pod
	}
	return podSet, nil
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

func GetContainerStatus(ctx context.Context, c containerd.Container) (string, uint32) {
	// Just in case, there is something wrong in server.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	titleCaser := cases.Title(language.English)

	task, err := c.Task(ctx, nil)
	if err != nil {
		// NOTE: NotFound doesn't mean that container hasn't started.
		// In docker/CRI-containerd plugin, the task will be deleted
		// when it exits. So, the status will be "created" for this
		// case.
		if errdefs.IsNotFound(err) {
			return titleCaser.String(string(containerd.Created)), 100
		}
		return titleCaser.String(string(containerd.Unknown)), 100
	}

	status, err := task.Status(ctx)
	if err != nil {
		return titleCaser.String(string(containerd.Unknown)), 100
	}
	labels, err := c.Labels(ctx)
	if err != nil {
		return titleCaser.String(string(containerd.Unknown)), 100
	}

	switch s := status.Status; s {
	case containerd.Stopped:
		if labels[restart.StatusLabel] == string(containerd.Running) && restart.Reconcile(status, labels) {
			return fmt.Sprintf("Restarting (%v) %s", status.ExitStatus, TimeSinceInHuman(status.ExitTime)), status.ExitStatus
		}
		return fmt.Sprintf("Exited (%v) %s", status.ExitStatus, TimeSinceInHuman(status.ExitTime)), status.ExitStatus
	case containerd.Running:
		return "Up", 100 // TODO: print "status.UpTime" (inexistent yet)
	default:
		return titleCaser.String(string(s)), 100
	}
}

func TimeSinceInHuman(since time.Time) string {
	return fmt.Sprintf("%s ago", units.HumanDuration(time.Since(since)))
}
