package url

const (
	HttpScheme string = "http://"
	HostURL    string = "localhost"
	Port       string = ":8080"

	Prefix string = HttpScheme + HostURL + Port

	APIV1 string = "/api/v1"

	PodStatus          string = "/pods/status"
	PodStatusGetURL    string = PodStatus + "/get"
	PodStatusGetAllURL string = PodStatus + "/getall"
	PodStatusDelURL    string = PodStatus + "/del"
	PodStatusApplyURL  string = PodStatus + "/apply"
	PodStatusUpdateURL string = PodStatus + "/update"

	DeploymentStatus          string = "/deployment/status"
	DeploymentStatusGetURL    string = DeploymentStatus + "/get"
	DeploymentStatusGetAllURL string = DeploymentStatus + "/getall"
	DeploymentStatusDelURL    string = DeploymentStatus + "/del"
	DeploymentStatusApplyURL  string = DeploymentStatus + "/apply"
	DeploymentStatusUpdateURL string = DeploymentStatus + "/update"
)
