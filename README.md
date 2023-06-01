## Simple usage

-   cd minik8s
-   make
-   /bin/kubeadm version
-   /bin/vctl [command flag]
-   /bin/nervctl [command flag]

## Build a Simple Pod using nervctl
-   ./bin/nervctl runp test
-   add two container to the pod
    -   ./bin/nervctl run golang:latest go1 8000:8000 /root/test_vo:/mnt container:test bash
    -   ./bin/nervctl run golang:latest go2 8000:8000 /root/test_vo:/go/src container:test bash
-   after the above step, we can find that this two container shared vloume and network (can use localhost to communicate)

## CNI TEST
-   First, run a pod using nervctl and check its ipaddress
    ```
    nerdctl inspect -f '{{.NetworkSettings.IPAddress}}' test
    ```
    ![本地路径](./docs/cni-ip.png "相对路径演示")

-   Second, run http server in the pod

    ![本地路径](./docs/cni-run.png "相对路径演示")

-   Third, using curl in other node to test the network

    ![本地路径](./docs/cni-test.png "相对路径演示")


## Deployment Controller 流程

### 启动

`./bin/apiserver`启动apiserver  
`./bin/kubelet`启动kubelet  
`./bin/scheduler`启动scheduler  
`./bin/kubeadm join --config=./testcases/vmeet1.yaml`加入集群  
`./bin/kube-controller-manager`启动所有controller  

### 使用Deployment controller创建deployment实例  

`./bin/kubectl apply <filename>`(expmple: `./bin/kubectl apply ./cmd/kubectl/app/src/test_deployment.yaml`)创建deployment实例  
`nerdctl ps`可以看到启动了`replicas`数量的pod和container  

## Autoscaler Controller 流程

先启动Deployment Controller，然后执行`./bin/kubectl apply <filename>`(expmple: `./bin/kubectl apply ./cmd/kubectl/app/src/test_autoscaler.yaml`)创建autoscaler实例  

## Service功能

> Service中的selector字段会去匹配pod的label字段，进而实现Servcie端口的映射(目前仅支持app标签匹配)

### 启动
`./bin/apiserver`启动apiserver  
`./bin/kubelet`启动kubelet  
`./bin/scheduler`启动scheduler  
`./bin/kubeadm join --config=./testcases/vmeet2.yaml`加入集群  

### 使用

`./bin/kubectl apply <service.yaml>` 创建Service服务
`./bin/kubectl get service` 查看创建的Service的状态、ClusterIP以及Port等信息
`./bin/kubectl delete <serviceName>` 删除name为serviceName的Service服务

### 使用效果
> Apply Service后，可以通过虚拟IP访问服务，同时IPtables规则也会修改，具体效果如下所示

![service-result](./docs/img/service-result.png)


## GPU Usage
### Start Up
> ./bin/kubelet
-   effect
![gpu-up](./docs/img/GPU-up.png)
### Apply Job using kubectl
> ./bin/kubectl apply ./testcases/job.yaml
-   effect
![pod-result](./docs/img/pod-run1.png)
-   after job finish
![pod-result](./docs/img/pod-run2.png)
> run `cat /root/minik8s/minik8s/scripts/data/result.out` can get the job's result

## Function Abstraction
### Setup
- First, we need to pull python image to build up our environment
```
$ nerdctl pull python:3.8.10
```

- Then, we need to run other supporting components for serverless function
```
$ ./bin/apiserver
$ ./bin/scheduler
$ ./bin/kubelet
$ ./bin/kube-controller-manager
$ ./bin/knative
```

### Run Serverless Function
- Register Function
> We need register a function to Knative first
```
$ ./bin/kubectl register <path-to-python-file>
```

- Trigger Function
> After registration, we can invoke function 

```
$ ./bin/kubectl trigger <funcName> <data>
```
> - "funcName" in the prefix of python file name. For example, if we have a python file named as "func.py", then its "funcName" is "func".
> 
> - "data" is a JSON string, and its value should be determined by the function. A valid data form is like '{"x":1,"y":2}' (for "Add.py" in ./testcases). By the way, if the function do not need params, we should set "data" to ''(empty string).

