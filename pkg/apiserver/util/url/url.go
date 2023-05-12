package url

const (
	HttpScheme string = "http://"
	HostURL string = "localhost"
	Port string = ":8080"

	Prefix string = HttpScheme + HostURL + Port
	
	APIV1 string = "/api/v1"

	PodURL string = APIV1 + "/pod"
	PodStatusGetURL string = PodURL + "/status/get"
	PodStatusGetAllURL string = PodURL + "/status/getall"
	PodStatusDelURL string = PodURL + "/status/del"
	PodStatusPutURL string = PodURL + "/status/put"
)