## Simple usage

-   cd minik8s
-   make
-   /bin/kubeadm version
-   /bin/vctl [command flag]
-   /bin/nervctl [command flag]

## Build a Simple Pod using nervctl
-   ./bin/nervctl runp test
-   add two container to the pod
    -   ./bin/nervctl run nginx:latest go1 8000:8000 /root/test_vo:/mnt container:test bash
    -   ./bin/nervctl run golang:latest go2 8000:8000 /root/test_vo:/go/src container:test bash
-   after the above step, we can find that this two container shared vloume and network (can use localhost to communicate)