#!/bin/bash
echo "clean etcd"

<<<<<<< HEAD
#/opt/etcd-v3.4.26/etcdctl del --prefix /node
/opt/etcd-v3.4.26/etcdctl del --prefix /deployment
/opt/etcd-v3.4.26/etcdctl del --prefix /service
/opt/etcd-v3.4.26/etcdctl del --prefix /pods

nerdctl stop $(nerdctl ps -a)
nerdctl rm $(nerdctl ps -a)
=======
/opt/etcd-v3.4.26/etcdctl del --prefix /node
<<<<<<< HEAD
/opt/etcd-v3.4.26/etcdctl del --prefix /deployment_test
=======
/opt/etcd-v3.4.26/etcdctl del --prefix /job
>>>>>>> 246bdcbc8103eb3a17c48085310ffc5a23577d09
/opt/etcd-v3.4.26/etcdctl del --prefix /service
/opt/etcd-v3.4.26/etcdctl del --prefix /pods
>>>>>>> develop
