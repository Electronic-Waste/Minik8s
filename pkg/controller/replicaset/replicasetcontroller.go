package replicaset

import (
	"context"
	"k8s/utils/concurrentmap"
	"k8s/utils/queue"
	"time"
)

// since client and listwatch are not completed yet, this file simply shows the workflow of controller
// all stuct and func related with listwatch are deleted since they are only for test

// TODO: add client, informer/listwatch
type ReplicaSetController struct {
	stopChannel <-chan struct{}

	// working queue holds key of RS
	queue queue.ConcurrentQueue
	// actually holds RS
	cp *concurrentmap.ConcurrentMap
	//Client
	//ls
}

type Pod struct {
}

// TODO: add client, informer/listwatch
// NewReplicaSetController configures a replica set controller with the specified event recorder
// k8s uses broadcaster to package listwatch, for simplicity, just use listwatch is ok
func NewReplicaSetController(burstReplicas int) *ReplicaSetController {
	//listwatch/informer
	return &ReplicaSetController{
		queue: queue.ConcurrentQueue(),
		cp:    concurrentmap.NewConcurrentMap(),
	}
}

func (rsc *ReplicaSetController) Run(ctx context.Context) {
	// TODO: start informer/listwatch
	go rsc.worker(ctx)
}

func (rsc *ReplicaSetController) worker(ctx context.Context) {
	// TODO:
	for {
		if !rsc.queue.Empty() {
			key := rsc.queue.Front()
			rsc.queue.Dequeue()
			rsc.syncReplicaSet(ctx, key.(string))
		} else {
			time.Sleep(time.Second)
		}
	}
}
func (rsc *ReplicaSetController) syncReplicaSet(ctx context.Context, key string) error {
	return nil
}

func (rsc *ReplicaSetController) syncReplicaSet(ctx context.Context, key string) error {
	// get expected replica set
	rs, _ := rsc.cp.Get(key).(*object.ReplicaSet)
	// get all actual pods of the rs
	allPods := GetAllPods()
	// filter all inactive pods
	activePods := FilterActivePods(allPods)
	// manage pods
	rsc.manageReplicas()
	// update status
	err = putReplicaSet(rs)
	return err
}

// get all RS pods
func GetAllPods() []*Pod {
	return nil
}

// filter all active pods
func FilterActivePods([]*Pod) {

}

// add or delete according to the spec and state
func (rsc *ReplicaSetController) manageReplicas(spec, state int) {
	if state < spec {
		addPods()
	}
	if state > spec {
		deletePods()
	}
}
func addPods() {
	//Client.addPod
}

func deletePods() {
	//Client.deletePod
}

func (rsc *ReplicaSetController) putReplicaSet() {
	//Client.putRS
}
