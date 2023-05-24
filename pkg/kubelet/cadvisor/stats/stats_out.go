package stats

import (
	"fmt"
	v1 "github.com/containerd/cgroups/stats/v1"
	v2 "github.com/containerd/cgroups/v2/stats"
	"time"
)

// Useful interface
func GetStatsEntryV1(data *v1.Metrics, previous *ContainerStats) (StatsEntry, error) {
	mem := calculateCgroupMemUsage(data)
	memLimit := float64(data.Memory.Usage.Limit)
	return StatsEntry{
		CPUPercentage:    calculateCgroupCPUPercent(previous, data),
		Memory:           mem,
		MemoryLimit:      memLimit,
		MemoryPercentage: calculateMemPercent(memLimit, mem),
	}, nil
}

// Useful interface
func GetStatsEntryV2(data *v2.Metrics, previous *ContainerStats) (StatsEntry, error) {
	fmt.Println("have not support v2 yet hhhhhh")
	return StatsEntry{}, nil
}

func calculateCgroupCPUPercent(previousStats *ContainerStats, metrics *v1.Metrics) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(metrics.CPU.Usage.Total) - float64(previousStats.CgroupCPU)
		// calculate the change for the entire system between readings
		//UserDelta = float64(metrics.CPU.Usage.User) - float64(previousStats.CgroupSystem)
		UserDelta = float64(time.Since(previousStats.Time).Nanoseconds())
	)

	if UserDelta > 0.0 && cpuDelta > 0.0 {
		// for the reason that we only hace 2 cpu cores, so we use 2 to replace the online cpu cores
		cpuPercent = (cpuDelta / UserDelta) * 100.0
	}
	return cpuPercent
}

func calculateCgroupMemUsage(metrics *v1.Metrics) float64 {
	if v := metrics.Memory.TotalInactiveFile; v < metrics.Memory.Usage.Usage {
		return float64(metrics.Memory.Usage.Usage - v)
	}
	return float64(metrics.Memory.Usage.Usage)
}

func calculateMemPercent(limit float64, usedNo float64) float64 {
	// Limit will never be 0 unless the container is not running and we haven't
	// got any data from cgroup
	if limit != 0 {
		return usedNo / limit * 100.0
	}
	return 0
}

// ----------------------------- Above is Usage interface ---------------------------------------//

func PrintCpuStatusV1(data *v1.Metrics) {
	fmt.Printf("total  cpu usage is %v\n", data.CPU.Usage.Total)
	fmt.Printf("kernel cpu usage is %v\n", data.CPU.Usage.Kernel)
	fmt.Printf("user   cpu usage is %v\n", data.CPU.Usage.User)
	fmt.Printf("Percpu cpu len   is %v\n", len(data.CPU.Usage.PerCPU))
	fmt.Printf("Percpu cpu usage is %v\n", (data.CPU.Usage.PerCPU))
	return
}

func PrintCpuStatusV2(data *v2.Metrics) {
	fmt.Printf("total  cpu usage is %v\n", data.CPU.UsageUsec)
	fmt.Printf("kernel cpu usage is %v\n", data.CPU.SystemUsec)
	fmt.Printf("user   cpu usage is %v\n", data.CPU.UserUsec)
}

func PrintCgroupMemUsage(metrics *v1.Metrics) {
	if v := metrics.Memory.TotalInactiveFile; v < metrics.Memory.Usage.Usage {
		fmt.Printf("Memory usage is %v\n", float64(metrics.Memory.Usage.Usage-v))
	}
	fmt.Printf("Memory usage is %v\n", float64(metrics.Memory.Usage.Usage))
}

func PrintCgroup2MemUsage(metrics *v2.Metrics) {
	if v := metrics.Memory.InactiveFile; v < metrics.Memory.Usage {
		fmt.Printf("Memory usage is %v\n", float64(metrics.Memory.Usage-v))
	}
	fmt.Printf("Memory usage is %v\n", float64(metrics.Memory.Usage))
}

func PrintMemoryLimitV1(metrics *v1.Metrics) {
	fmt.Printf("Memory Limit is %v\n", metrics.Memory.Usage.Limit)
}

func PrintMemoryLimitV2(metrics *v2.Metrics) {
	fmt.Printf("Memory Limit is %v\n", metrics.Memory.UsageLimit)
}
