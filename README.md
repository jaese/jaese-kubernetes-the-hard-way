# Kubernetes The Hard Way done with Vagrant and Ansible

differences from

minimize the need to ssh into nodes.

--kubelet-preferred-address-types=InternalIP

## 1. Prerequisites

VirtualBox

Ansible

## 3. Provisioning Compute Resources

```sh
vagrant up

vagrant ssh

vagrant ssh-config >> ~/.ssh/config
````

## 4. Provisioning a CA and Generating TLS Certificates

```sh
mkdir _local

go run ./provision-certs
````

## 5. Generating Kubernetes Configuration Files for Authentication

go run ./generate-kubeconfigs

## 6. Generating the Data Encryption Config and Key

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

## 7. Bootstrapping the etcd Cluster

```sh
curl -L https://github.com/etcd-io/etcd/releases/download/v3.4.15/etcd-v3.4.15-linux-amd64.tar.gz -o _local/etcd-v3.4.15-linux-amd64.tar.gz \
  && tar xvf _local/etcd-v3.4.15-linux-amd64.tar.gz --directory _local

ansible-playbook -i inventory.yaml etcd_playbook.yaml

ETCDCTL_API=3 etcdctl member list \
  --endpoints=https://192.168.60.20:2379 \
  --cacert=_local/ca.pem \
  --cert=_local/kubernetes.pem \
  --key=_local/kubernetes-key.pem
```

## 8. Bootstrapping the Kubernetes Control Plane

```sh
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-apiserver -o $OUTPUT_DIR/kube-apiserver
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-controller-manager -o $OUTPUT_DIR/kube-controller-manager
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-scheduler -o $OUTPUT_DIR/kube-scheduler
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kubectl -o $OUTPUT_DIR/kubectl

ansible-playbook -i inventory controllers_playbook.yaml

kubectl cluster-info --kubeconfig _local/admin.kubeconfig

kubectl apply -f kubelet-authorization.yaml --kubeconfig _local/admin.kubeconfig
```


## Notes

pod-cidr=10.200.0.0/16

mkdir _local
