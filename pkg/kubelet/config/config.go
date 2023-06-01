package config

import (
	"fmt"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/kubelet/types"
	"minik8s.io/pkg/podmanager"
	"minik8s.io/pkg/util/config"
	"strings"
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

// !!! attention : this function can only be call in the lock protected
func (p *PodStorage) IsPodExist(name string) bool {
	for _, m := range p.storage {
		for p_name, _ := range m {
			if strings.Compare(p_name, name) == 0 {
				return true
			}
		}
	}
	return false
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
			// running all Pod when types is SET
			p.storeLock.Lock()
			for _, p := range mes.Pods {
				if !podmanager.IsPodRunning(p.Name) && !podmanager.IsCrashContainer(p.Name) {

				} else if podmanager.IsPodRunning(p.Name) {
					podmanager.DelPod(p.Name)
				} else if podmanager.IsCrashContainer(p.Name) {
					podmanager.DelSimpleContainer(p.Name)
				}
			}

			// change the storage and Pod status
			p.storage[mes.Source] = make(map[string]*core.Pod)
			for _, pod := range mes.Pods {
				// running core Pod
				fmt.Printf("pod name is %s\n", pod.Name)
				err := podmanager.RunSysPod(pod)
				if err != nil {
					fmt.Println(err)
					return err
				}

				// update storage
				p.storage[mes.Source][pod.Name] = pod
			}
			p.storeLock.Unlock()
		}
	case types.ADD:
		{
			// run a new Pod
			p.storeLock.Lock()
			// do some error handling
			// check for exist Pod name
			for _, pod := range mes.Pods {
				if p.IsPodExist(pod.Name) {
					fmt.Println("Pod name crash in Pod adding in file source")
				} else {
					// run the pod
					err := podmanager.RunSysPod(pod)
					if err != nil {
						fmt.Println(err)
						return err
					}
					p.storage[mes.Source][pod.Name] = pod
				}
			}
			p.storeLock.Unlock()
		}
	case types.DELETE:
		{
			// delete a Pod
			p.storeLock.Lock()
			for _, pod := range mes.Pods {
				if podmanager.IsPodRunning(pod.Name) {
					podmanager.DelPod(pod.Name)
				}
				delete(p.storage[mes.Source], pod.Name)
			}
			p.storeLock.Unlock()
		}
	case types.REMOVE:
		{
			fmt.Println("not support types")
		}
	case types.UPDATE:
		{
			// a Pod be updated
			fmt.Println("unsupport function")
		}
	case types.CHECK:
		{
			p.storeLock.Lock()
			fmt.Println("check the controller plane")
			pods, err := podmanager.GetPods()
			if err != nil {
				return err
			}
			is_find := false
			for _, p := range mes.Pods {
				for _, podVal := range pods {
					if strings.Compare(p.Name, podVal.Name) == 0 {
						is_find = true
						if strings.Compare(string(podVal.Status.Phase), "Running") == 0 {
							continue
						} else {
							podmanager.DelPod(p.Name)
							podmanager.RunSysPod(p)
						}
					}
				}
				if !is_find {
					podmanager.DelSimpleContainer(p.Name)
					podmanager.RunSysPod(p)
				}
				is_find = false
			}
			fmt.Println("finish check")
			p.storeLock.Unlock()
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
