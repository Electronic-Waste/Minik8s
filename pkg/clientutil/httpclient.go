package clientutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"net/http"
)

// return: error
// @objType: type want to apply; @obj: the obj to be applied
func HttpApply(objType string, obj any) error {
	client := http.Client{}
	payload, _ := json.Marshal(obj)
	switch objType {
	case "Deployment":
		request, err := http.NewRequest("POST", apiurl.DeploymentStatusApplyURL, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
			//return errors.New("")
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("put fail")
		}
	}
	return nil
}

// return: obj queried
// @objType: type want to get; @params: query params
func HttpGet(objType string, params map[string]string) ([]byte, error) {
	client := http.Client{}
	urlparam := ""
	//if there are params, construct get url
	if len(params) != 0 {
		urlparam += "?"
		i := 0
		for k, v := range params {
			urlparam += k + "=" + v
			if i != len(params)-1 {
				urlparam += "&"
			}
			i++
		}
	}
	switch objType {
	case "Deployment":
		request, err := http.NewRequest("GET", apiurl.DeploymentStatusGetURL+urlparam, nil)
		if err != nil {
			log.Fatal(err)
		}
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
			//return errors.New("")
		}
		if response.StatusCode != http.StatusOK {
			return nil, errors.New("get fail")
		}
		data, err := ioutil.ReadAll(response.Body)
		return data, nil
	}
	return nil, errors.New("invalid request")
}

// return: error
// @objType: type want to get; @params: query params
func HttpDel(objType string, params map[string]string) ([]byte, error) {
	client := http.Client{}
	urlparam := ""
	//if there are params, construct get url
	if len(params) != 0 {
		urlparam += "?"
		i := 0
		for k, v := range params {
			urlparam += k + "=" + v
			if i != len(params)-1 {
				urlparam += "&"
			}
			i++
		}
	}
	switch objType {
	case "Deployment":
		request, err := http.NewRequest("DELETE", apiurl.DeploymentStatusDelURL+urlparam, nil)
		if err != nil {
			log.Fatal(err)
		}
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
			//return errors.New("")
		}
		if response.StatusCode != http.StatusOK {
			return nil, errors.New("del fail")
		}
		data, err := ioutil.ReadAll(response.Body)
		return data, nil
	}
}
