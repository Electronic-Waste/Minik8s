package controlleroptions

// DeploymentControllerConfiguration contains elements describing DeploymentController.
type DeploymentControllerOptions struct {
	// concurrentDeploymentSyncs is the number of deployment objects that are
	// allowed to sync concurrently. Larger number = more responsive deployments,
	// but more CPU (and network) load.
	ConcurrentDeploymentSyncs int32
}

//temporarily set for 5
func (opts *DeploymentControllerOptions) InitOptions() {
	opts.ConcurrentDeploymentSyncs = 5
}
