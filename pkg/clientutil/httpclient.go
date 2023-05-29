package clientutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"minik8s.io/pkg/apis/core"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"net/http"
)

// return: error
// @objType: type want to apply; @obj: the obj to be applied
func HttpApply(objType string, obj any) error {
	fmt.Println("http apply")
	client := http.Client{}
	payload, _ := json.Marshal(obj)
	switch objType {
	case "Autoscaler":
		urlparam := "?namespace=default"
		request, err := http.NewRequest("POST", apiurl.Prefix+apiurl.AutoscalerStatusApplyURL+urlparam, bytes.NewReader(payload))
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
	case "Deployment":
		urlparam := "?namespace=default"
		request, err := http.NewRequest("POST", apiurl.Prefix+apiurl.DeploymentStatusApplyURL+urlparam, bytes.NewReader(payload))
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
		request, err := http.NewRequest("POST", apiurl.Prefix+apiurl.PodStatusApplyURL, bytes.NewReader(payload))
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
	case "DNS":
		var dns core.DNS
		json.Unmarshal([]byte(payload), &dns)
		postURL := apiurl.Prefix + apiurl.DNSApplyURL + fmt.Sprintf("?namespace=default&name=%s", dns.Name)
		request, err := http.NewRequest("POST", postURL, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("http apply dns")
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("apply fail")
		}
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response: %s\n", string(body))
	case "Node":
		var node core.Node
		json.Unmarshal([]byte(payload), &node)
		postURL := apiurl.HttpScheme + node.Spec.MasterIp + apiurl.Port + apiurl.NodeRergisterUrl
		request, err := http.NewRequest("POST", postURL, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("register fail")
		}
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response: %s\n", string(body))
	case "Job":
		job := obj.(core.Job)
		// send http request to apiserver
		request, err := http.NewRequest("POST", apiurl.Prefix+apiurl.JobApplyUrl, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("http apply job name is %s\n", job.Meta.Name)
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("apply fail")
		}
	}

	return nil
}

func HttpPlus(objType string, obj any, url string) (error, string) {
	fmt.Println("http plus apply: ",url)
	client := http.Client{}
	payload, _ := json.Marshal(obj)
	var res string
	switch objType {
	case "Deployment":
		//urlparam := "?namespace=default"
		request, err := http.NewRequest("POST", url, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("http apply deployment")
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("apply fail"), ""
		}
	case "Pod":
		request, err := http.NewRequest("POST", url, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("http plus apply pod")
		response, err := client.Do(request)
		fmt.Println(response)
		if err != nil {
			fmt.Println(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("apply fail"), ""
		}
		fmt.Println("http plus apply pod success")
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response: %s\n", string(body))
		res = string(body)
	case "Service":
		var service core.Service
		json.Unmarshal([]byte(payload), &service)
		// fmt.Printf("httpclient: service is : %v\n", service)
		// fmt.Printf("httpclinet: service name is : %v\n", service.Name)
		postURL := url + fmt.Sprintf("?namespace=default&name=%s", service.Name)
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
			return errors.New("apply fail"), ""
		}
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response: %s\n", string(body))
	case "Node":
		var node core.Node
		json.Unmarshal([]byte(payload), &node)
		postURL := url
		request, err := http.NewRequest("POST", postURL, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("register fail"), ""
		}
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response: %s\n", string(body))
	case "Job":
		job := obj.(core.Job)
		fmt.Printf("job name is %s\n", job.Meta.Name)
		postURL := url
		request, err := http.NewRequest("POST", postURL, bytes.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		if response.StatusCode != http.StatusOK {
			return errors.New("register fail"), ""
		}
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response: %s\n", string(body))
	}

	return nil, res
}

//get from kubelet
func HttpGetPlus(objType string, url string) ([]byte,error) {
	client := http.Client{}
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("http get pod metrics",url)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		return nil,errors.New("get plus fail")
	}
	data, err := ioutil.ReadAll(response.Body)
	return data, nil
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
	fmt.Printf("httpclient get params: %s\n",urlparam)
	var requestUrl string
	switch objType {
	case "Autoscaler":
		requestUrl = apiurl.Prefix + apiurl.AutoscalerStatusGetURL + urlparam
	case "Pod":
		requestUrl = apiurl.Prefix + apiurl.PodStatusGetURL + urlparam
	case "Deployment":
		requestUrl = apiurl.Prefix + apiurl.DeploymentStatusGetURL + urlparam
	case "nodes":
		requestUrl = apiurl.Prefix + apiurl.NodesGetUrl + urlparam
	case "jobs":
		requestUrl = apiurl.Prefix + apiurl.JobGetUrl
	case "metrics":	//params: name, nodeip
		//requestUrl = apiurl.HttpScheme + apiurl.Vmeet1IP + apiurl.Port + apiurl.MetricsGetUrl + urlparam
		requestUrl = apiurl.Prefix + apiurl.MetricsGetUrl + urlparam
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

	//return nil, errors.New("invalid request")
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
	fmt.Printf("httpclient get params: %s\n",urlparam)
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
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println("get with prefix not ok")
		return nil, errors.New("get fail")
	}
	data, err := ioutil.ReadAll(response.Body)
	//fmt.Println(data)
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
	case "Service":
		requestUrl = apiurl.Prefix + apiurl.ServiceGetAllURL
	case "DNS":
		requestUrl = apiurl.Prefix + apiurl.DNSGetAllURL
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
func HttpDel(objType string, params map[string]string) error {
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
	case "Service":
		requestUrl = apiurl.Prefix + apiurl.ServiceDelURL + urlparam
	case "Pod":
		requestUrl = apiurl.Prefix + apiurl.PodStatusDelURL + urlparam
	case "DNS":
		requestUrl = apiurl.Prefix + apiurl.DNSDelURL + urlparam
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

