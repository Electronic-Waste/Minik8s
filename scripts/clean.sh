#!/bin/bash
echo "clean etcd"

#/opt/etcd-v3.4.26/etcdctl del --prefix /node
/opt/etcd-v3.4.26/etcdctl del --prefix /deployment
#/opt/etcd-v3.4.26/etcdctl del --prefix /service
/opt/etcd-v3.4.26/etcdctl del --prefix /pods
#/opt/etcd-v3.4.26/etcdctl del --prefix /job
/opt/etcd-v3.4.26/etcdctl del --prefix /autoscaler
/opt/etcd-v3.4.26/etcdctl del --prefix /func

nerdctl stop $(nerdctl ps -a)
nerdctl rm $(nerdctl ps -a)
