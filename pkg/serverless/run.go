package serverless

import(
	"net/http"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

type Bootstrap interface {
	Run()
}

type Knative struct {

}
