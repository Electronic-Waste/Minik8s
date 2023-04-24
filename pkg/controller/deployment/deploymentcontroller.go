package deployment

import (
	"context"
	"k8s/utils/concurrentmap"
	"k8s/utils/queue"
	"time"
)

// wait for deployment to be finished
// since client and listwatch are not completed yet, this file simply shows the workflow of controller
// this controller relate to listwatch less but to etcd storage more, so the previous non-usable only-for-test version is deleted

const (
	// maxRetries is the number of times a deployment will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the times
	// a deployment is going to be requeued:
	maxRetries = 15
)

// TODO: add client and informet/listwatch
type DeploymentController struct {
	//Client
	//listwatch

	// work queue
	queue         queue.ConcurrentQueue
	deploymentMap *concurrentmap.ConcurrentMapTrait[string, object.VersionedDeployment]
	replicasetMap *concurrentmap.ConcurrentMapTrait[string, object.ReplicaSet]
}

// this is decided by etcd
type Deployment struct {
}

// TODO: add client and informet/listwatch
func NewDeploymentController(ctx context.Context) (*DeploymentController, error) {
	dc := &DeploymentController{
		queue: queue.ConcurrentQueue(),
	}
	return dc, nil
}

func (dc *DeploymentController) Run(ctx context.Context) {
	go dc.worker(ctx)
}

func (dc *DeploymentController) worker(ctx context.Context) {
	for {
		if !dc.queue.Empty() {
			key := dc.queue.Front()
			dc.queue.Dequeue()
			dc.syncDeployment(ctx, key.(string))
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (dc *DeploymentController) syncDeployment(ctx context.Context, key string) {
	return nil
}

// add deployment to etcd
func (dc *DeploymentController) addDeployment(deployment *Deployment) {
	/*
		rsUidNew := uuid.New().String()
			rsNameNew := deployment.Metadata.Name + rsUidNew
			rsKeyNew := path.Join(config.RSConfigPrefix, rsNameNew)
			rs := object.ReplicaSet{
				ObjectMeta: object.ObjectMeta{
					Name:   rsNameNew,
					Labels: deployment.Metadata.Labels,
					UID:    rsUidNew,
					OwnerReferences: []object.OwnerReference{{
						Kind:       "Deployment",
						Name:       deployment.Metadata.Name,
						UID:        deployment.Metadata.UID,
						Controller: false,
					}},
				},
				Spec: object.ReplicaSetSpec{
					Replicas: deployment.Spec.Replicas,
					Template: deployment.Spec.Template,
				},
			}
			dc.dm2rs.Put(res.Key, rsKeyNew)

			err = client.Put(dc.apiServerBase+rsKeyNew, rs)
			if err != nil {
				klog.Errorf("Error send new rs to etcd\n")
			}
			dc.replicasetMap.Put(rsKeyNew, rs)
	*/
}

// delete deployment
func (dc *DeploymentController) deleteDeployment() {

}
