package cadvisor

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/containerd/cgroups/stats/v1"
	v2 "github.com/containerd/cgroups/v2/stats"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/typeurl/v2"
	"minik8s.io/pkg/cli/remote_cli"
	"minik8s.io/pkg/idutil/containerwalker"
	"minik8s.io/pkg/kubelet/cadvisor/stats"
)

type CAdvisor struct {
}

func GetContainerMetric(id string) error {
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return err
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")
	container, err := cli.Client().LoadContainer(ctx, id)
	err = getContainerMetric(container, ctx)
	return err
}

// get a metric of container which id is provided
func getContainerMetric(container containerd.Container, ctx context.Context) error {
	//task is in the for loop to avoid nil task just after Container creation
	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	metric, err := task.Metrics(ctx)
	if err != nil {
		return err
	}
	anydata, err := typeurl.UnmarshalAny(metric.Data)
	if err != nil {
		return err
	}
	err = printMetric(anydata)
	if err != nil {
		return err
	}
	return nil
}

// the interface for debug use
func printMetric(anydata interface{}) error {
	var (
		data  *v1.Metrics
		data2 *v2.Metrics
	)

	switch v := anydata.(type) {
	case *v1.Metrics:
		data = v
	case *v2.Metrics:
		data2 = v
	default:
		err := errors.New("cannot convert metric data to cgroups.Metrics")
		return err
	}

	if data != nil {
		stats.PrintCpuStatusV1(data)
		stats.PrintCgroupMemUsage(data)
		stats.PrintMemoryLimitV1(data)
	} else if data2 != nil {
		stats.PrintCpuStatusV2(data2)
		stats.PrintCgroup2MemUsage(data2)
		stats.PrintMemoryLimitV2(data2)
	}

	fmt.Println("finish one print")

	return nil
}

func (c *CAdvisor) GetAllPodMetric() map[string]stats.StatsEntry {
	return nil
}

// given a Pod name
func (c *CAdvisor) GetPodMetric(name string) (stats.StatsEntry, error) {
	// we don't include the pause container's resource usage
	// find all container labeled with the 'name'
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return err
	}
	walker := &containerwalker.ContainerWalker{
		Client: cli.Client(),
		OnFound: func(ctx context.Context, found containerwalker.Found) error {
			fmt.Println(found.Container.ID())

			return nil
		},
	}

	// !!! : need to specify the namespace of finding container
	ctx := namespaces.WithNamespace(context.Background(), "default")
	_, err = walker.WalkPod(ctx, name)
	if err != nil {
		return stats.StatsEntry{}, err
	}
	return stats.StatsEntry{}, nil
}
