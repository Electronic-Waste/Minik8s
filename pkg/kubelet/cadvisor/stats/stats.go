package stats

import (
	"sync"
	"time"
)

// StatsEntry represents the statistics data collected from a container
type StatsEntry struct {
	Container        string
	Name             string
	ID               string
	CPUPercentage    float64
	Memory           float64
	MemoryLimit      float64
	MemoryPercentage float64
}

type PodStats struct {
	Name             string
	CPUPercentage    float64
	MemoryPercentage float64
}

// Stats represents an entity to store containers statistics synchronously
type Stats struct {
	mutex sync.RWMutex
	StatsEntry
	err error
}

// ContainerStats represents the runtime container stats
type ContainerStats struct {
	Time                        time.Time
	CgroupCPU, Cgroup2CPU       uint64
	CgroupSystem, Cgroup2System uint64
}

// NewStats is from https://github.com/docker/cli/blob/3fb4fb83dfb5db0c0753a8316f21aea54dab32c5/cli/command/container/formatter_stats.go#L113-L116
func NewStats(container string) *Stats {
	return &Stats{StatsEntry: StatsEntry{Container: container}}
}

// SetStatistics is from https://github.com/docker/cli/blob/3fb4fb83dfb5db0c0753a8316f21aea54dab32c5/cli/command/container/formatter_stats.go#L87-L93
func (cs *Stats) SetStatistics(s StatsEntry) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	s.Container = cs.Container
	cs.StatsEntry = s
}

// GetStatistics is from https://github.com/docker/cli/blob/3fb4fb83dfb5db0c0753a8316f21aea54dab32c5/cli/command/container/formatter_stats.go#L95-L100
func (cs *Stats) GetStatistics() StatsEntry {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	return cs.StatsEntry
}

// GetError is from https://github.com/docker/cli/blob/3fb4fb83dfb5db0c0753a8316f21aea54dab32c5/cli/command/container/formatter_stats.go#L51-L57
func (cs *Stats) GetError() error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	return cs.err
}
