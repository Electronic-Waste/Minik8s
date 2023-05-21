package config

import (
	"fmt"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/kubelet/types"
	"minik8s.io/pkg/util/config"
	"sync"
)

// need to implement the merge interface(Merger)
type PodStorage struct {
	update chan types.PodUpdate

	// map from source to a map from name to Pod
	// !!! : we assume the name is unique in this app, but if we consider different namespace that concludion is wrong
	storage map[string]map[string]*core.Pod
	// lock protect the storage
	storeLock sync.RWMutex
}

// this method don't check for update, just check for add and remove
func (p *PodStorage) SendUpdate(oldPods map[string]*core.Pod, newPods []*core.Pod, update chan types.PodUpdate, source string) error {
	// use name to make sure the Pod is the same
	var samePods map[string]*core.Pod
	var addUpdate types.PodUpdate
	addUpdate.Op = types.ADD
	addUpdate.Source = source
	for _, p := range newPods {
		if _, ok := oldPods[p.Name]; ok {
			samePods[p.Name] = p
		} else {
			// get the PodUpdate
			// add first
			addUpdate.Pods = append(addUpdate.Pods, p)
		}
	}
	update <- addUpdate

	// TODO(wjl): check update in the samePods

	// construct remove
	var delUpdate types.PodUpdate
	delUpdate.Op = types.DELETE
	delUpdate.Source = source
	for k, p := range oldPods {
		if _, ok := samePods[k]; ok {
			continue
		} else {
			delUpdate.Pods = append(delUpdate.Pods, p)
		}
	}
	update <- delUpdate
	p.storeLock.RUnlock()
	p.storeLock.Lock()
	// update podStorage
	p.storage[source] = map[string]*core.Pod{}
	for _, np := range newPods {
		p.storage[source][np.Name] = np
	}
	p.storeLock.Unlock()
	return nil
}

// we need to combine all of the update of Pod and send to Kubelet
func (p *PodStorage) Merge(source string, update interface{}) error {
	mes := update.(types.PodUpdate)
	switch mes.Op {
	case types.SET:
		{
			// check Pod is running or not first

		}
	case types.ADD:
		{
			// run a new Pod
			fmt.Println("not support types")
		}
	case types.DELETE:
		{
			// delete a Pod
			fmt.Println("not support types")
		}
	case types.REMOVE:
		{
			fmt.Println("not support types")
		}
	case types.UPDATE:
		{
			fmt.Println("not support types")
		}
	}
	return nil
}

func NewPodStorage(ch chan types.PodUpdate) *PodStorage {
	return &PodStorage{
		update:  ch,
		storage: map[string]map[string]*core.Pod{},
	}
}

type PodConfig struct {
	// that is the last channel used to hold all Pod message
	update chan types.PodUpdate

	storage *PodStorage

	mux *config.Mux
}

func (p *PodConfig) Updates() chan types.PodUpdate {
	return p.update
}

func NewPodConfig() *PodConfig {
	ch := make(chan types.PodUpdate, 50)
	ps := NewPodStorage(ch)
	return &PodConfig{
		update:  ch,
		storage: ps,
		mux:     config.NewMux(ps),
	}
}

func (p *PodConfig) Channel(source string) chan interface{} {
	return p.mux.BuildNewChan(source)
}
