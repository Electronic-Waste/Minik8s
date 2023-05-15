package url

const (
	HttpScheme string = "http://"
	HostURL string = "localhost"
	Port string = ":8080"

	Prefix string = HttpScheme + HostURL + Port

	PodStatus string = "/pods/status"
	PodStatusGetURL string = PodStatus + "/get"
	PodStatusGetAllURL string = PodStatus + "/getall"
	PodStatusDelURL string = PodStatus + "/del"
	PodStatusPutURL string = PodStatus + "/put"
)