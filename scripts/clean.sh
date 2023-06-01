#!/bin/bash
echo "clean etcd"

#/opt/etcd-v3.4.26/etcdctl del --prefix /node
/opt/etcd-v3.4.26/etcdctl del --prefix /deployment
#/opt/etcd-v3.4.26/etcdctl del --prefix /service
<<<<<<< HEAD
/opt/etcd-v3.4.26/etcdctl del --prefix /pods
#/opt/etcd-v3.4.26/etcdctl del --prefix /job
#/opt/etcd-v3.4.26/etcdctl del --prefix /autoscaler
/opt/etcd-v3.4.26/etcdctl del --prefix /func
=======
#/opt/etcd-v3.4.26/etcdctl del --prefix /pods
#/opt/etcd-v3.4.26/etcdctl del --prefix /job
#/opt/etcd-v3.4.26/etcdctl del --prefix /autoscaler
#/opt/etcd-v3.4.26/etcdctl del --prefix /func
>>>>>>> fix(serverless): fix wrong pod delete

nerdctl stop $(nerdctl ps -a)
nerdctl rm $(nerdctl ps -a)
