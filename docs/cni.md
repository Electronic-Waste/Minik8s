## CNI Note

-   install etcd cluster
    详见知乎--ETCD集群安装配置
-   install flannel plugin
    ```
    从GitHub拉取最新的flannel发行版
    mkdir -p /opt/flannel
    tar xzf flannel-vx.y.z-linux-amd64.tar.gz -C /opt/flannel
    cd /opt/flannel && ls 
    ```
-   启动flannel(在每个node节点都要跑)
    ```
    flanneld  --ip-masq --kube-subnet-mgr=false --etcd-endpoints=http://127.0.0.1:2379 
    ```
-   配置新的flannel的cni配置
    ```
    # vim /etc/cni/net.d/10-flannel.conflist
    {
      "name": "flannel",
      "cniVersion": "0.3.1",
      "plugins": [
        {
          "type": "flannel",
          "delegate": {
            "isDefaultGateway": true
          }
        },
        {
          "type": "portmap",
          "capabilities": {
            "portMappings": true
          }
        }
      ]
    }
    ```
-   此时查看containerd网络：
    ```
    nerdctl network ls 
    ```
    运行一个由flannel管理的容器：
    ```
    nerdctl run -d --net flannel --name flannel nginx:latest
    ```
    这样就可以使用curl ip:80来访问nginx页面