package controlleroptions

// ReplicaSetControllerConfiguration contains elements describing ReplicaSetController.
type ReplicaSetControllerOptions struct {
	// concurrentRSSyncs is the number of replica sets that are  allowed to sync
	// concurrently. Larger number = more responsive replica  management, but more
	// CPU (and network) load.
	ConcurrentRSSyncs int32
}

//temporarily set for 5
func (opts *ReplicaSetControllerOptions) InitOptions() {
	opts.ConcurrentRSSyncs = 5
}
