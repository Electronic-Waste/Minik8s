package url

const (
	HttpScheme    string = "http://"
	HostURL       string = "localhost"
	Port          string = ":8080"
	SchedulerPort string = "1234"

	Prefix string = HttpScheme + HostURL + Port

	PodStatus                 string = "/pods/status"
	PodStatusGetURL           string = PodStatus + "/get"
	PodStatusGetAllURL        string = PodStatus + "/getall"
	PodStatusDelURL           string = PodStatus + "/del"
	PodStatusApplyURL         string = PodStatus + "/apply"
	PodStatusUpdateURL        string = PodStatus + "/update"
	PodStatusGetWithPrefixURL string = PodStatus + "/getwithprefix"

	Service          string = "/service"
	ServiceGetURL    string = Service + "/get"
	ServiceGetAllURL string = Service + "/getall"
	ServiceDelURL    string = Service + "/del"
	ServiceApplyURL  string = Service + "/apply"
	ServiceUpdateURL string = Service + "/update"

	DNS          string = "/dns"
	DNSGetURL    string = DNS + "/get"
	DNSGetAllURL string = DNS + "/getall"
	DNSDelURL    string = DNS + "/del"
	DNSApplyURL  string = DNS + "/apply"
	DNSUpdateURL string = DNS + "/update"

	DeploymentStatus          string = "/deployment/status"
	DeploymentStatusGetURL    string = DeploymentStatus + "/get"
	DeploymentStatusGetAllURL string = DeploymentStatus + "/getall"
	DeploymentStatusDelURL    string = DeploymentStatus + "/del"
	DeploymentStatusApplyURL  string = DeploymentStatus + "/apply"
	DeploymentStatusUpdateURL string = DeploymentStatus + "/update"

	AutoscalerStatus          string = "/autoscaler/status"
	AutoscalerStatusGetURL    string = AutoscalerStatus + "/get"
	AutoscalerStatusGetAllURL string = AutoscalerStatus + "/getall"
	AutoscalerStatusDelURL    string = AutoscalerStatus + "/del"
	AutoscalerStatusApplyURL  string = AutoscalerStatus + "/apply"
	AutoscalerStatusUpdateURL string = AutoscalerStatus + "/update"

	Node             string = "/node"
	NodeRergisterUrl string = Node + "/register"
	NodesGetUrl      string = Node + "/getall"

	Sched         string = "/sched"
	SchedApplyURL        = Sched + "/apply"

	Metrics 		string = "/metrics"
	MetricsGetUrl 	string = Metrics + "/get"

	Vmeet1IP string = "192.168.1.5"
	Vmeet2IP string = "192.168.1.6"
	Vmeet3IP string = "192.168.1.7"
)
