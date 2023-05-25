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
	"sync"
	"time"
)

type ContainerListener struct {
	// actual container stats
	DataStore map[string]*stats.Stats
	// mid data useful in compute the container stats
	PreStats map[string]stats.ContainerStats

	mutex sync.RWMutex
}

func GetNewListener() *ContainerListener {
	return &ContainerListener{
		DataStore: make(map[string]*stats.Stats),
		PreStats:  make(map[string]stats.ContainerStats),
	}
}

func (c *ContainerListener) RegisterContainer(id string) error {
	// start a new thread to listen to the container state
	// collect prestats first
	conStats, _, err := c.GetContainerStats(id)
	if err != nil {
		return err
	}
	c.PreStats[id] = conStats
	go c.collect(id)
	return nil
}

// stats which has been deal with
func (c *ContainerListener) GetStats(id string) stats.StatsEntry {
	for {
		if _, ok := c.DataStore[id]; ok {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return c.DataStore[id].GetStatistics()
}

func (c *ContainerListener) collect(id string) error {
	// using a for loop to collect the stats
	for {
		newPrevious, any, err := c.GetContainerStats(id)
		if err != nil {
			return err
		}
		var (
			data  *v1.Metrics
			data2 *v2.Metrics
		)

		switch v := any.(type) {
		case *v1.Metrics:
			data = v
			tmp := c.PreStats[id]
			Constats, err := stats.GetStatsEntryV1(data, &(tmp))
			if err != nil {
				return err
			}
			if _, ok := c.DataStore[id]; !ok {
				c.DataStore[id] = stats.NewStats(id)
			}
			c.DataStore[id].SetStatistics(Constats)
			c.mutex.RLock()
			c.PreStats[id] = newPrevious
			c.mutex.RUnlock()
		case *v2.Metrics:
			data2 = v
			tmp := c.PreStats[id]
			Constats, err := stats.GetStatsEntryV2(data2, &(tmp))
			if err != nil {
				return err
			}
			if _, ok := c.DataStore[id]; !ok {
				c.DataStore[id] = stats.NewStats(id)
			}
			c.DataStore[id].SetStatistics(Constats)
			c.mutex.RLock()
			c.PreStats[id] = newPrevious
			c.mutex.RUnlock()
		default:
			err := errors.New("cannot convert metric data to cgroups.Metrics")
			return err
		}

		// get data per second
		time.Sleep(1 * time.Second)
	}
}

// that raw stats
func (c *ContainerListener) GetContainerStats(id string) (stats.ContainerStats, interface{}, error) {
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return stats.ContainerStats{}, nil, err
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")
	container, err := cli.Client().LoadContainer(ctx, id)
	task, err := container.Task(ctx, nil)
	if err != nil {
		return stats.ContainerStats{}, nil, err
	}

	metric, err := task.Metrics(ctx)
	if err != nil {
		return stats.ContainerStats{}, nil, err
	}
	anydata, err := typeurl.UnmarshalAny(metric.Data)
	if err != nil {
		return stats.ContainerStats{}, nil, err
	}
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
		return stats.ContainerStats{}, nil, err
	}

	conStats := stats.ContainerStats{
		Time: time.Now(),
	}

	if data != nil {
		conStats.CgroupCPU = data.CPU.Usage.Total
		conStats.CgroupSystem = data.CPU.Usage.Kernel
	} else if data2 != nil {
		conStats.Cgroup2CPU = data2.CPU.UsageUsec * 1000
		conStats.Cgroup2System = data2.CPU.SystemUsec * 1000
	}
	return conStats, anydata, nil
}

type CAdvisor struct {
	DataStore map[string]stats.PodStats
	// use to listen to container status
	containerListener *ContainerListener
}

func GetCAdvisor() *CAdvisor {
	return &CAdvisor{
		DataStore:         make(map[string]stats.PodStats),
		containerListener: GetNewListener(),
	}
}

func (c *CAdvisor) RegisterAllPod() error {
	// register all the container in the Pod to containerListener
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return err
	}
	// !!! : need to specify the namespace of finding container
	ctx := namespaces.WithNamespace(context.Background(), "default")
	walker := &containerwalker.ContainerWalker{
		Client: cli.Client(),
		OnFound: func(ctx context.Context, found containerwalker.Found) error {
			labels, err := found.Container.Labels(ctx)
			if err != nil {
				return err
			}
			fmt.Printf("try to register pod %s\n", labels["nerdctl/name"])
			err = c.RegisterPod(labels["nerdctl/name"])
			return err
		},
	}

	_, err = walker.WalkPod(ctx, "pause")
	if err != nil {
		return err
	}
	return nil
}

func (c *CAdvisor) RegisterPod(name string) error {
	// register all the container in the Pod to containerListener
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return err
	}
	walker := &containerwalker.ContainerWalker{
		Client: cli.Client(),
		OnFound: func(ctx context.Context, found containerwalker.Found) error {
			fmt.Printf("try to register container %s\n", found.Container.ID())
			err := c.containerListener.RegisterContainer(found.Container.ID())
			return err
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

func (c *CAdvisor) GetAllPodMetric() map[string]stats.PodStats {
	// register all the container in the Pod to containerListener
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
<<<<<<< HEAD
		fmt.Println(err)
=======
>>>>>>> 2b2df8301f69d16c66e716990a86247423ba0861
		return nil
	}
	// !!! : need to specify the namespace of finding container
	ctx := namespaces.WithNamespace(context.Background(), "default")
	resMap := make(map[string]stats.PodStats)
	walker := &containerwalker.ContainerWalker{
		Client: cli.Client(),
		OnFound: func(ctx context.Context, found containerwalker.Found) error {
			labels, err := found.Container.Labels(ctx)
			if err != nil {
				return err
			}
			fmt.Printf("try to register pod %s\n", labels["nerdctl/name"])
			status, err := c.GetPodMetric(labels["nerdctl/name"])
			resMap[labels["nerdctl/name"]] = status
			return err
		},
	}

	_, err = walker.WalkPod(ctx, "pause")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return resMap
}

// given a Pod name
func (c *CAdvisor) GetPodMetric(name string) (stats.PodStats, error) {
	// we don't include the pause container's resource usage
	// find all container labeled with the 'name'
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return stats.PodStats{}, nil
	}
	cpuPer := 0.0
	mem := 0.0
	memLimit := 0.0
	walker := &containerwalker.ContainerWalker{
		Client: cli.Client(),
		OnFound: func(ctx context.Context, found containerwalker.Found) error {
			fmt.Println(found.Container.ID())
			ConStats := c.containerListener.GetStats(found.Container.ID())
			cpuPer += ConStats.CPUPercentage
			mem += ConStats.Memory
			memLimit += ConStats.MemoryLimit
			return nil
		},
	}

	// !!! : need to specify the namespace of finding container
	ctx := namespaces.WithNamespace(context.Background(), "default")
	_, err = walker.WalkPod(ctx, name)
	if err != nil {
		return stats.PodStats{}, err
	}
	return stats.PodStats{
		Name:             name,
		CPUPercentage:    cpuPer,
		MemoryPercentage: mem / memLimit,
	}, nil
}

func GetContainerMetric(id string) (stats.StatsEntry, error) {
	cli, err := remote_cli.NewRemoteRuntimeService(remote_cli.IdenticalErrorDelay)
	if err != nil {
		return stats.StatsEntry{}, err
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")
	container, err := cli.Client().LoadContainer(ctx, id)
	entry, err := getContainerMetric(container, ctx)
	return entry, err
}

// get a metric of container which id is provided
func getContainerMetric(container containerd.Container, ctx context.Context) (stats.StatsEntry, error) {
	//task is in the for loop to avoid nil task just after Container creation
	task, err := container.Task(ctx, nil)
	if err != nil {
		return stats.StatsEntry{}, err
	}

	metric, err := task.Metrics(ctx)
	if err != nil {
		return stats.StatsEntry{}, err
	}
	anydata, err := typeurl.UnmarshalAny(metric.Data)
	if err != nil {
		return stats.StatsEntry{}, err
	}
	err = printMetric(anydata)
	if err != nil {
		return stats.StatsEntry{}, err
	}
	return stats.StatsEntry{}, err
}

func getMetric(any interface{}) (stats.StatsEntry, error) {
	var (
		data  *v1.Metrics
		data2 *v2.Metrics
	)

	switch v := any.(type) {
	case *v1.Metrics:
		data = v
	case *v2.Metrics:
		data2 = v
	default:
		err := errors.New("cannot convert metric data to cgroups.Metrics")
		return stats.StatsEntry{}, err
	}

	if data != nil {

	} else if data2 != nil {

	}

	return stats.StatsEntry{}, nil
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
