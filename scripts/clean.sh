#!/bin/bash
echo "clean etcd"

#/opt/etcd-v3.4.26/etcdctl del --prefix /node
/opt/etcd-v3.4.26/etcdctl del --prefix /deployment_test
/opt/etcd-v3.4.26/etcdctl del --prefix /service
/opt/etcd-v3.4.26/etcdctl del --prefix /pods

nerdctl stop $(nerdctl ps -a)
nerdctl rm $(nerdctl ps -a)