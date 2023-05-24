package clientutil

import (
	"fmt"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"net/http"
	"minik8s.io/pkg/apis/core"
)

// return: error
// @objType: type want to apply; @obj: the obj to be applied
func HttpApply(objType string, obj any) error {
	fmt.Println("http apply")
	client := http.Client{}
	payload, _ := json.Marshal(obj)
	switch objType {
	case "Deployment":
		urlparam := "?namespace=default"
		request, err := http.NewRequest("POST", apiurl.Prefix + apiurl.DeploymentStatusApplyURL + urlparam, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("http apply deployment")
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("apply fail")
		}
	case "Pod":
		request, err := http.NewRequest("POST", apiurl.Prefix + apiurl.PodStatusApplyURL, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("http apply pod")
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("apply fail")
		}
	case "Service":
		var service core.Service
		json.Unmarshal([]byte(payload), &service)
		// fmt.Printf("httpclient: service is : %v\n", service)
		// fmt.Printf("httpclinet: service name is : %v\n", service.Name)
		postURL := apiurl.Prefix + apiurl.ServiceApplyURL + fmt.Sprintf("?namespace=default&name=%s", service.Name)
		fmt.Printf("httpclient: send request to %s\n", postURL)
		request, err := http.NewRequest("POST", postURL, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("http apply service")
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("apply fail")
		}
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response: %s\n", string(body))
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
		request, err := http.NewRequest("GET", apiurl.Prefix + apiurl.DeploymentStatusGetURL+urlparam, nil)
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
	var(
		data []byte
	)
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
		request, err := http.NewRequest("DELETE", apiurl.Prefix + apiurl.DeploymentStatusDelURL+urlparam, nil)
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
		data, err = ioutil.ReadAll(response.Body)
	case "Service":
		request, _ := http.NewRequest("DELETE", apiurl.Prefix + apiurl.ServiceDelURL + urlparam, nil)
		fmt.Printf("URL is %s\n", request.URL)
		response, err := client.Do(request)
		if err != nil {
			return []byte{}, err
		}
		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("del service fail")
		}
		data, err = ioutil.ReadAll(response.Body)
	}
	return data, nil
}