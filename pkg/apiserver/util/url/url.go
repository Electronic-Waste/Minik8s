package url

const (
	HttpScheme    string = "http://"
	HostURL       string = "localhost"
	Port          string = ":8080"
	SchedulerPort string = "1234"

	Prefix string = HttpScheme + HostURL + Port

	PodStatus          string = "/pods/status"
	PodStatusGetURL    string = PodStatus + "/get"
	PodStatusGetAllURL string = PodStatus + "/getall"
	PodStatusDelURL    string = PodStatus + "/del"
	PodStatusApplyURL  string = PodStatus + "/apply"
	PodStatusUpdateURL string = PodStatus + "/update"

	Service          string = "/service"
	ServiceGetURL    string = Service + "/get"
	ServiceGetAllURL string = Service + "/getall"
	ServiceDelURL    string = Service + "/del"
	ServiceApplyURL  string = Service + "/apply"
	ServiceUpdateURL string = Service + "/update"

	DeploymentStatus          string = "/deployment/status"
	DeploymentStatusGetURL    string = DeploymentStatus + "/get"
	DeploymentStatusGetAllURL string = DeploymentStatus + "/getall"
	DeploymentStatusDelURL    string = DeploymentStatus + "/del"
	DeploymentStatusApplyURL  string = DeploymentStatus + "/apply"
	DeploymentStatusUpdateURL string = DeploymentStatus + "/update"

	Node             string = "/node"
	NodeRergisterUrl string = Node + "/register"

	Sched         string = "/sched"
	SchedApplyURL        = Sched + "/apply"
)
