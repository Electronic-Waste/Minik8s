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
)

// return: error
// @objType: type want to apply; @obj: the obj to be applied
func HttpApply(objType string, obj any) error {
	fmt.Println("http apply")
	client := http.Client{}
	payload, _ := json.Marshal(obj)
	urlparam := "?namespace=default"
	var requestUrl string
	switch objType {
	case "Autoscaler":
		requestUrl = apiurl.Prefix + apiurl.AutoscalerStatusApplyURL + urlparam
		fmt.Println("http apply autoscaler")
	case "Deployment":
		requestUrl = apiurl.Prefix + apiurl.DeploymentStatusApplyURL + urlparam
		fmt.Println("http apply deployment")
	case "Pod":
		requestUrl = apiurl.Prefix + apiurl.PodStatusApplyURL + urlparam
		fmt.Println("http apply pod")
	}
	request, err := http.NewRequest("POST", requestUrl, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}
	
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("apply fail")
	}
	return nil
}

// return: error
// @objType: type want to update; @obj: the obj to be updated
func HttpUpdate(objType string, obj any) error {
	fmt.Println("http update")
	client := http.Client{}
	payload, _ := json.Marshal(obj)
	urlparam := "?namespace=default"
	var requestUrl string
	switch objType {
	case "Autoscaler":
		requestUrl = apiurl.Prefix + apiurl.AutoscalerStatusUpdateURL + urlparam
		fmt.Println("http update autoscaler")
	case "Deployment":
		requestUrl = apiurl.Prefix + apiurl.DeploymentStatusUpdateURL + urlparam
		fmt.Println("http update deployment")
	case "Pod":
		requestUrl = apiurl.Prefix + apiurl.PodStatusUpdateURL + urlparam
		fmt.Println("http update pod")
	}
	request, err := http.NewRequest("POST", requestUrl, bytes.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}
	
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("update fail")
	}
	return nil
}

// return: obj queried
// @objType: type want to get; @params: query params: namespace & name
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
	//fmt.Printf("httpclient get params: %s\n",urlparam)
	var requestUrl string
	switch objType {
	case "Autoscaler":
		requestUrl = apiurl.Prefix + apiurl.AutoscalerStatusGetURL + urlparam
	case "Deployment":
		requestUrl = apiurl.Prefix + apiurl.DeploymentStatusGetURL + urlparam
	case "Pod":
		requestUrl = apiurl.Prefix + apiurl.PodStatusGetURL + urlparam
	}
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		//return nil,err
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("get fail")
	}
	data, err := ioutil.ReadAll(response.Body)
	return data, nil
}

// return: obj queried
// @objType: type want to get; @params: query params: namespace & prefix
func HttpGetWithPrefix(objType string, params map[string]string) ([]byte, error) {
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
	//fmt.Printf("httpclient get params: %s\n",urlparam)
	var requestUrl string
	switch objType {
	case "Pod":
		requestUrl = apiurl.Prefix + apiurl.PodStatusGetWithPrefixURL + urlparam
	}
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
		//return nil,err
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("get fail")
	}
	data, err := ioutil.ReadAll(response.Body)
	return data, nil
}

// return: obj queried
// @objType: type want to get
func HttpGetAll(objType string) ([]byte, error) {
	client := http.Client{}
	var requestUrl string
	switch objType {
	case "Autoscaler":
		requestUrl = apiurl.Prefix + apiurl.AutoscalerStatusGetAllURL
	case "Deployment":
		requestUrl = apiurl.Prefix + apiurl.DeploymentStatusGetAllURL
	case "Pod":
		requestUrl = apiurl.Prefix + apiurl.PodStatusGetAllURL
	}
	request, err := http.NewRequest("GET", requestUrl, nil)
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

// return: error
// @objType: type want to get; @params: query params: namespace & name
func HttpDel(objType string, params map[string]string)  error {
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
	var requestUrl string
	switch objType {
	case "Autoscaler":
		requestUrl = apiurl.Prefix + apiurl.AutoscalerStatusDelURL + urlparam
	case "Deployment":
		requestUrl = apiurl.Prefix + apiurl.DeploymentStatusDelURL + urlparam
	case "Pod":
		requestUrl = apiurl.Prefix + apiurl.PodStatusDelURL + urlparam
	}
	request, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		//return errors.New("")
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("del fail")
	}
	return nil
}