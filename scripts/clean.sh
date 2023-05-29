#!/bin/bash
echo "clean etcd"

/opt/etcd-v3.4.26/etcdctl del --prefix /node
/opt/etcd-v3.4.26/etcdctl del --prefix /job