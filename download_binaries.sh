#!/bin/sh

set -xe

OUTPUT_DIR=_local

curl -L https://github.com/etcd-io/etcd/releases/download/v3.4.15/etcd-v3.4.15-linux-amd64.tar.gz -o _local/etcd-v3.4.15-linux-amd64.tar.gz \
  && tar xvf _local/etcd-v3.4.15-linux-amd64.tar.gz --directory $OUTPUT_DIR

curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-apiserver -o $OUTPUT_DIR/kube-apiserver
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-controller-manager -o $OUTPUT_DIR/kube-controller-manager
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-scheduler -o $OUTPUT_DIR/kube-scheduler
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kubectl -o $OUTPUT_DIR/kubectl

curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.21.0/crictl-v1.21.0-linux-amd64.tar.gz -o _local/crictl-v1.21.0-linux-amd64.tar.gz \
  && tar xvf _local/crictl-v1.21.0-linux-amd64.tar.gz --directory _local
curl -L https://github.com/opencontainers/runc/releases/download/v1.0.0-rc93/runc.amd64 -o _local/runc
curl -L https://github.com/containernetworking/plugins/releases/download/v0.9.1/cni-plugins-linux-amd64-v0.9.1.tgz -o _local/cni-plugins-linux-amd64-v0.9.1.tgz
curl -L https://github.com/containerd/containerd/releases/download/v1.4.4/containerd-1.4.4-linux-amd64.tar.gz -o _local/containerd-1.4.4-linux-amd64.tar.gz
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-proxy -o _local/kube-proxy
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kubelet -o _local/kubelet
