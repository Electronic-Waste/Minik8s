package url

const (
	HttpScheme string = "http://"
	HostURL    string = "localhost"
	Port       string = ":8080"

	Prefix string = HttpScheme + HostURL + Port

	APIV1 string = "/api/v1"

	PodStatus          string = "/pods/status"
	PodStatusGetURL    string = APIV1 + PodStatus + "/get"
	PodStatusGetAllURL string = APIV1 + PodStatus + "/getall"
	PodStatusDelURL    string = APIV1 + PodStatus + "/del"
	PodStatusPutURL    string = APIV1 + PodStatus + "/put"

	DeploymentStatus          string = "/deployment/status"
	DeploymentStatusGetURL    string = APIV1 + DeploymentStatus + "/get"
	DeploymentStatusGetAllURL string = APIV1 + DeploymentStatus + "/getall"
	DeploymentStatusDelURL    string = APIV1 + DeploymentStatus + "/del"
	DeploymentStatusPurURL    string = APIV1 + DeploymentStatus + "/put"
)
