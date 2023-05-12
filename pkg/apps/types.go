package apps

type MetaData struct {
	Name      string
	Label     string
	Namespace string
	//……
}
type PodTemplate struct {
	//pod implement
}
type Deployment struct {
	Metadata MetaData
	Spec     DeploymentSpec
	Status   DeploymentStatus
}

type DeploymentSpec struct {
	Replicas int
	Template PodTemplate
	Selector string //must match .spec.template.metadata.labels
	//strategy	DeploymentStrategy
}

type DeploymentStatus struct {
	ObservedGeneration int
	AvailableReplicas  int
	//for later use
	UpdatedReplicas int
	ReadyReplicas   int
}
