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