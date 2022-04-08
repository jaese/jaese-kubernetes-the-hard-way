package main

import (
	"github.com/codeskyblue/go-sh"

	"github.com/jaese/jaese-kubernetes-the-hard-way/utils"
)

var OutPathPrefix = "_local/"

var WorkerNodeNames = []string{"worker-0", "worker-1"}

var IPAddresses = map[string]string{
	"controller-0": "192.168.60.20",
	"worker-0":     "192.168.60.30",
	"worker-1":     "192.168.60.31",
}

// Deploying only one master instance ('controller-0') in order to avoid having
// to set up a reverse proxy.
var KubernetesPublicAddress = IPAddresses["controller-0"]

func main() {
	for _, name := range WorkerNodeNames {
		generateKubeconfigFile(&kubeconfigSpec{
			User:          "system:node:" + name,
			ClientCert:    name + ".pem",
			ClientKey:     name + "-key.pem",
			OutFilePrefix: name,
		})
	}

	generateKubeconfigFile(&kubeconfigSpec{
		User:          "system:kube-proxy",
		ClientCert:    "kube-proxy.pem",
		ClientKey:     "kube-proxy-key.pem",
		OutFilePrefix: "kube-proxy",
	})

	generateKubeconfigFile(&kubeconfigSpec{
		User:          "system:kube-controller-manager",
		ClientCert:    "kube-controller-manager.pem",
		ClientKey:     "kube-controller-manager-key.pem",
		OutFilePrefix: "kube-controller-manager",
	})

	generateKubeconfigFile(&kubeconfigSpec{
		User:          "system:kube-scheduler",
		ClientCert:    "kube-scheduler.pem",
		ClientKey:     "kube-scheduler-key.pem",
		OutFilePrefix: "kube-scheduler",
	})

	generateKubeconfigFile(&kubeconfigSpec{
		User:          "admin",
		ClientCert:    "admin.pem",
		ClientKey:     "admin-key.pem",
		OutFilePrefix: "admin",
	})
}

type kubeconfigSpec struct {
	User          string
	ClientCert    string
	ClientKey     string
	OutFilePrefix string
}

func generateKubeconfigFile(spec *kubeconfigSpec) {
	utils.MustRun(sh.Command(
		"kubectl", "config", "set-cluster", "kubernetes-the-hard-way",
		"--certificate-authority="+OutPathPrefix+"ca.pem",
		"--embed-certs=true",
		"--server=https://"+KubernetesPublicAddress,
		"--kubeconfig="+OutPathPrefix+spec.OutFilePrefix+".kubeconfig"))

	utils.MustRun(sh.Command(
		"kubectl", "config", "set-credentials", spec.User,
		"--client-certificate="+OutPathPrefix+spec.ClientCert,
		"--client-key="+OutPathPrefix+spec.ClientKey,
		"--embed-certs=true",
		"--kubeconfig="+OutPathPrefix+spec.OutFilePrefix+".kubeconfig"))

	utils.MustRun(sh.Command(
		"kubectl", "config", "set-context", "default",
		"--cluster=kubernetes-the-hard-way",
		"--user="+spec.User,
		"--kubeconfig="+OutPathPrefix+spec.OutFilePrefix+".kubeconfig"))

	utils.MustRun(sh.Command(
		"kubectl", "config", "use-context", "default", "--kubeconfig="+OutPathPrefix+spec.OutFilePrefix+".kubeconfig"))
}
