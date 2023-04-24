package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type RestClient struct {
	// base is the root URL for all invocations of the client
	base *url.URL
	// versionedAPIPath is a path segment connecting the base URL to the resource root
	versionedAPIPath string
	// Set specific behavior of the client.  If not set http.DefaultClient will be used.
	Client *http.Client
}

func NewRestClient(baseURL *url.URL, client *http.Client) (*RestClient, error) {
	return &RestClient{
		base:   baseURL,
		Client: client,
	}, nil
}

// basic client options

func Get(url *url.URL) ([]byte, error) {
	request, err := http.NewRequest("GET", url.String(), nil)
	response, err := http.DefaultClient.Do(request)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("StatusCode not 200")
	}
	reader := response.Body
	data, err := io.ReadAll(reader)
	defer response.Body.Close()
	return data, nil
}

func Put(url *url.URL, obj any) error {
	data, err := json.Marshal(obj)
	reader := bytes.NewReader(data)
	request, err := http.NewRequest("PUT", url.String(), reader)
	request.Header.Add("Content-Type", "application/json")
}

func Delete(url *url.URL) error {
	request, err := http.NewRequest("DELETE", url.String(), nil)
	response, err := http.DefaultClient.Do(request)
	if response.StatusCode != http.StatusOK {
		return errors.New("StatusCode not 200")
	}
	return nil
}
