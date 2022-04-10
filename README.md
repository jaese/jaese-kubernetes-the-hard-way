# jaese-kubernetes-the-hard-way (Vagrant, Ansible)

Scripts and config files I produced as I followed [kelseyhightower's Kubernetes The Hard Way](https://github.com/kelseyhightower/kubernetes-the-hard-way) tutorial.

The major differences from the original are:

* [Vagrant](https://www.vagrantup.com) is used to provision necessary VMs instead of using the Google Cloud Platform.
* Is not a highly available setup; only one control plane node ('controller-0') is provisioned so that I can skip setting up a load balancer.
* In order to minimize the amount of `ssh`ing and running commands in each node, the Kubernetes control and worker planes are set up with [Ansible Playbooks](https://www.ansible.com). Because of this, nodes can be configured and modified in a declarative manner.

Minor differences:

* The nodes are assigned IP addresses from subnet `192.168.60.0/24` instead of `10.240.0.0/24`.
* `--kubelet-preferred-address-types=InternalIP` is added to the kube-apiserver command so that the kube-apiserver would not rely on worker nodes' hostnames being resolvable.

I only tested on Mac.

## Notes on each chapter of the tutorial

### 1. Prerequisites

* Vagrant with VirtualBox
* Ansible
* [Go](https://go.dev) for running scripts. TODO: Bash is a more sensible choice for a scripting language.

### 2. Installing the Client Tools

Install the `etcdctl` command, in addition to the ones in the original tutorial.

### 3. Provisioning Compute Resources

With Vagrant and VirtualBox installed, `vagrant up` is enough to launch and set up controller and three worker VMs.

```sh
vagrant up
```

Set up SSH configs to allow SSHing into the VMs by hostname. This is necessary for running Ansible Playbooks later.

```sh
vagrant ssh-config >> ~/.ssh/config

# Example of SSHing into a node using hostname
ssh worker-0
```

### 4. Provisioning a CA and Generating TLS Certificates

Create `_local` in this directory. All generated and downloaded files will be stored in this directory.

```sh
mkdir _local
```

Run `provision-certs` script to provision all necessary certificates and keys (same as the original tutorial).

```sh
go run ./provision-certs
````

### 5. Generating Kubernetes Configuration Files for Authentication

Run `generate-kubeconfigs` to generate all necessary `kubeconfig` files to `_local`.

```sh
go run ./generate-kubeconfigs
```

### 6. Generating the Data Encryption Config and Key

Generate and save the encryption config file in `_local`.

```sh
ENCRYPTION_KEY=$(head -c 32 /dev/urandom | base64)

cat > _local/encryption-config.yaml <<EOF
kind: EncryptionConfig
apiVersion: v1
resources:
  - resources:
      - secrets
    providers:
      - aescbc:
          keys:
            - name: key1
              secret: ${ENCRYPTION_KEY}
      - identity: {}
EOF
```

### 7. Bootstrapping the etcd Cluster

Download the etcd binaries and run the playbook to set up etcd in the control node.

```sh
curl -L https://github.com/etcd-io/etcd/releases/download/v3.4.15/etcd-v3.4.15-linux-amd64.tar.gz -o _local/etcd-v3.4.15-linux-amd64.tar.gz \
  && tar xvf _local/etcd-v3.4.15-linux-amd64.tar.gz --directory _local

ansible-playbook -i inventory.yaml etcd_playbook.yaml
```

#### Verification

```
ETCDCTL_API=3 etcdctl member list \
  --endpoints=https://192.168.60.20:2379 \
  --cacert=_local/ca.pem \
  --cert=_local/kubernetes.pem \
  --key=_local/kubernetes-key.pem
```

### 8. Bootstrapping the Kubernetes Control Plane

Download the necessary binaries in `_local` and run the playbook to set up the control plane.

```sh
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-apiserver -o _local/kube-apiserver
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-controller-manager -o _local/kube-controller-manager
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-scheduler -o _local/kube-scheduler
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kubectl -o _local/kubectl

ansible-playbook -i inventory.yaml controllers_playbook.yaml
```

Verify that the control plane is running.

```sh
kubectl cluster-info --kubeconfig _local/admin.kubeconfig
```

#### RBAC for Kubelet Authorization

```sh
kubectl apply -f kubelet-authorization.yaml --kubeconfig _local/admin.kubeconfig
```

### 9. Bootstrapping the Kubernetes Worker Nodes

Download the worker binaries to `_local` and run the playbook.

```sh
curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.21.0/crictl-v1.21.0-linux-amd64.tar.gz -o _local/crictl-v1.21.0-linux-amd64.tar.gz \
  && tar xvf _local/crictl-v1.21.0-linux-amd64.tar.gz --directory _local
curl -L https://github.com/opencontainers/runc/releases/download/v1.0.0-rc93/runc.amd64 -o _local/runc
curl -L https://github.com/containernetworking/plugins/releases/download/v0.9.1/cni-plugins-linux-amd64-v0.9.1.tgz -o _local/cni-plugins-linux-amd64-v0.9.1.tgz
curl -L https://github.com/containerd/containerd/releases/download/v1.4.4/containerd-1.4.4-linux-amd64.tar.gz -o _local/containerd-1.4.4-linux-amd64.tar.gz
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-proxy -o _local/kube-proxy
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kubelet -o _local/kubelet

ansible-playbook -i inventory.yaml workers_playbook.yaml
```

#### Verification

```sh
kubectl get nodes --kubeconfig _local/admin.kubeconfig
```

### 10. Configuring kubectl for Remote Access

Same as the original except `KUBERNETES_PUBLIC_ADDRESS` is the IP address of controller-0 and the certificate and key files are in `_local`.

```sh
KUBERNETES_PUBLIC_ADDRESS=192.168.60.20

kubectl config set-cluster kubernetes-the-hard-way \
  --certificate-authority=_local/ca.pem \
  --embed-certs=true \
  --server=https://${KUBERNETES_PUBLIC_ADDRESS}:6443

kubectl config set-credentials admin \
  --client-certificate=_local/admin.pem \
  --client-key=_local/admin-key.pem

kubectl config set-context kubernetes-the-hard-way \
  --cluster=kubernetes-the-hard-way \
  --user=admin

kubectl config use-context kubernetes-the-hard-way
```

### 11. Provisioning Pod Network Routes

Nothing to do as this is already done when provisioning the VMs. (Linux routing table).

### 12. Deploying the DNS Cluster Add-on

Same as the original.

### 13. Smoke Test

Should work the same way as the original.

#### Data Encryption

```sh
kubectl create secret generic kubernetes-the-hard-way \
  --from-literal="mykey=mydata"

ETCDCTL_API=3 etcdctl get /registry/secrets/default/kubernetes-the-hard-way \
  --endpoints=https://192.168.60.20:2379 \
  --cacert=_local/ca.pem \
  --cert=_local/kubernetes.pem \
  --key=_local/kubernetes-key.pem | hexdump -C
```

### 14. Cleaning Up

The VMs can be shutdown and cleaned up with Vagrant commands:

```
vagrant halt
vagrant destroy
```
