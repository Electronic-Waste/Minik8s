package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/apiserver/etcd"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"net/http"
	"path"
)

func HandleNodeRegister(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	node := core.Node{}
	json.Unmarshal(body, &node)
	fmt.Println(node)
	if node.MetaData.Name == "" {
		fmt.Println("need to provide a name")
		resp.WriteHeader(http.StatusInternalServerError)
		err := errors.New("need to provide node name")
		resp.Write([]byte(err.Error()))
		return
	}
	etcdURL := path.Join(apiurl.Node, node.MetaData.Name)
	err := etcd.Put(etcdURL, string(body))
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
}

func HandleGetNodes(resp http.ResponseWriter, req *http.Request) {
	// get all node first
	fmt.Println("handle get nodes")
	NodeStrs, err := etcd.GetWithPrefix(apiurl.Node)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(err.Error()))
		return
	}
	nodeList := []core.Node{}
	fmt.Printf("the length of nodeStr is %d\n", len(NodeStrs))
	for _, str := range NodeStrs {
		node := core.Node{}
		json.Unmarshal([]byte(str), &node)
		nodeList = append(nodeList, node)
	}
	nodeArray := core.NodeList{
		NodeArray: nodeList,
	}
	jsonVal, err := json.Marshal(nodeArray)
	// Error occur in json parsing: return error to client
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonVal)
}
