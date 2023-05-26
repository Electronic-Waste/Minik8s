package dns

import (
	"net/http"
	"encoding/json"
	"path"
	"io/ioutil"
	"fmt"
	// "github.com/go-redis/redis/v8"

	"minik8s.io/pkg/util/listwatch"
	"minik8s.io/pkg/apiserver/etcd"
	"minik8s.io/pkg/apiserver/util/url"
	"minik8s.io/pkg/apis/core"
)

// Return certain dns info
// uri: /dns/get?namespace=...&name=...
// @namespace: namespace requested; @name: dns name
func HandleGetDNS(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	dnsName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || dnsName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	etcdKey := path.Join(url.DNS, namespace, dnsName)
	dns, err := etcd.Get(etcdKey)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(dns))
	return
}

// Return all dns infos
// uri: /dns/getall
func HandleGetAllDNS(resp http.ResponseWriter, req *http.Request) {
	etcdPrefix := url.DNS
	var dnsArr []string
	dnsArr, err := etcd.GetWithPrefix(etcdPrefix)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	var jsonVal []byte
	jsonVal, err = json.Marshal(dnsArr)
	// Error occur in json parsing: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonVal)
	// return
}


// Apply a dns rule
// uri: /dns/apply?namespace=...&name=...
// @namespace: namespace requested; @name: dns name
// body: core.DNS in JSON form
func HandleApplyDNS(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	dnsName := vars.Get("name")
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("HandleApplyDNS receive request: %v\n", string(body))

	dns := core.DNS{}
	json.Unmarshal(body, &dns)
	// Param miss: return error to client
	if namespace == "" || dnsName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	// 1. Query etcd to get clutserIP related to serviceName
	for i, subpath := range dns.Spec.Subpaths {
		var service core.Service
		serviceURL := path.Join(url.Service, namespace, subpath.Service)
		serviceString, err := etcd.Get(serviceURL)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			resp.Write([]byte(err.Error()))
			return
		}
		json.Unmarshal([]byte(serviceString), &service)
		dns.Spec.Subpaths[i].ClusterIP = service.Spec.ClusterIP
	}
	// 2. Persist to etcd
	etcdURL := path.Join(url.DNS, namespace, dnsName)
	jsonVal, _ := json.Marshal(dns)
	err := etcd.Put(etcdURL, string(jsonVal))
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// 3. Publish to redis & Send response to client
	// - topic: /dns/apply
	// - payload: <core.DNS>
	listwatch.Publish(url.DNSApplyURL, string(jsonVal))	
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprintf("Apply DNS rule %s successfully!", dnsName)))
}

// Delete a dns rule
// uri: /dns/del?namespace=...&name=...
// @namespace: namespace requested; @name: dns name
func HandleDelDNS(resp http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	namespace := vars.Get("namespace")
	dnsName := vars.Get("name")
	// Param miss: return error to client
	if namespace == "" || dnsName == "" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Name is missing"))
		return
	}
	// 1. Get hostName mapped by dnsName & Delete k-v pair in etcd
	var dns core.DNS
	etcdURL := path.Join(url.DNS, namespace, dnsName)
	output, err := etcd.Get(etcdURL)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	json.Unmarshal([]byte(output), &dns)
	hostName := dns.Spec.Host
	err = etcd.Del(etcdURL)
	// Error occur in etcd: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	// 2. Publish to redis & Send response to client
	// - topic: /dns/del
	// - payload: <hostName>
	listwatch.Publish(url.DNSDelURL, hostName)	
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprintf("Delete DNS rule %s successfully", dnsName)))
}