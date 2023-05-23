package stats

import (
	"fmt"
	v1 "github.com/containerd/cgroups/stats/v1"
	v2 "github.com/containerd/cgroups/v2/stats"
)

// Useful interface
func GetStatsEntryV1(data *v1.Metrics) (StatsEntry, error) {
	return StatsEntry{
		CPUPercentage:    0,
		Memory:           0,
		MemoryLimit:      0,
		MemoryPercentage: 0,
	}, nil
}

// Useful interface
func GetStatsEntryV2(data *v1.Metrics) (StatsEntry, error) {
	fmt.Println("have not support v2 yet hhhhhh")
	return StatsEntry{}, nil
}

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
