package main

import (
	"strings"
	"text/template"

	"github.com/codeskyblue/go-sh"

	"github.com/jaese/jaese-kubernetes-the-hard-way/utils"
)

var OutPathPrefix = "_local/"

var WorkerNodeNames = []string{"worker-0", "worker-1", "worker-2"}

var IPAddresses = map[string]string{
	"controller-0": "192.168.60.20",
	"worker-0":     "192.168.60.30",
	"worker-1":     "192.168.60.31",
	"worker-2":     "192.168.60.32",
}

// Deploying only one master instance ('controller-0') in order to avoid having
// to set up a reverse proxy.
var KubernetesPublicAddress = IPAddresses["controller-0"]

func main() {
	provisionCACert()
	provisionClientCerts()
}

func provisionCACert() {
	templated, err := utils.TextTemplateExecuteString(
		csrJSONTemplate,
		map[string]string{
			"CN": "Kubernetes",
			"O":  "Kubernetes",
		})
	if err != nil {
		panic(err)
	}

	utils.MustRunWithStringInput(
		templated,
		sh.Command(
			"cfssl",
			"gencert",
			"-initca",
			"-").
			Command(
				"cfssljson",
				"-bare",
				OutPathPrefix+"ca"))
}

func provisionClientCerts() {
	// The Admin Client Certificate
	provisionClientCert(&clientCertSpec{
		CN:            "admin",
		O:             "system:masters",
		OutFilePrefix: "admin",
	})

	// The Kubelet Client Certificates
	for _, name := range WorkerNodeNames {
		provisionClientCert(&clientCertSpec{
			CN:            "system:node:" + name,
			O:             "system:nodes",
			Hostname:      name + "," + IPAddresses[name],
			OutFilePrefix: name,
		})
	}

	// The Controller Manager Client Certificate
	provisionClientCert(&clientCertSpec{
		CN:            "system:kube-controller-manager",
		O:             "system:kube-controller-manager",
		OutFilePrefix: "kube-controller-manager",
	})

	// The Kube Proxy Client Certificate
	provisionClientCert(&clientCertSpec{
		CN:            "system:kube-proxy",
		O:             "system:node-proxier",
		OutFilePrefix: "kube-proxy",
	})

	// The Scheduler Client Certificate
	provisionClientCert(&clientCertSpec{
		CN:            "system:kube-scheduler",
		O:             "system:kube-scheduler",
		OutFilePrefix: "kube-scheduler",
	})

	// The Kubernetes API Server Certificate
	apiserverHostnames := []string{
		"10.32.0.1", // The service cluster IP address of the Kubernetes API server
		IPAddresses["controller-0"],
		KubernetesPublicAddress,
		"127.0.0.1",
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		"kubernetes.default.svc.cluster",
		"kubernetes.default.svc.cluster.local",
	}
	provisionClientCert(&clientCertSpec{
		CN:            "kubernetes",
		O:             "Kubernetes",
		Hostname:      strings.Join(apiserverHostnames, ","),
		OutFilePrefix: "kubernetes",
	})

	// The Service Account Key Pair
	provisionClientCert(&clientCertSpec{
		CN:            "service-accounts",
		O:             "Kubernetes",
		OutFilePrefix: "service-account",
	})
}

type clientCertSpec struct {
	// 'CN' is interpreted by the API server as user name.
	// See
	// https://kubernetes.io/docs/reference/access-authn-authz/authentication/#x509-client-certs
	CN string

	// 'O' is interpreted as the user's group memberships.
	O string

	Hostname      string
	OutFilePrefix string
}

func provisionClientCert(spec *clientCertSpec) {
	templated, err := utils.TextTemplateExecuteString(
		csrJSONTemplate,
		map[string]string{
			"CN": spec.CN,
			"O":  spec.O,
		})
	if err != nil {
		panic(err)
	}

	utils.MustRunWithStringInput(
		templated,
		sh.Command(
			"cfssl",
			"gencert",
			"-ca="+OutPathPrefix+"ca.pem",
			"-ca-key="+OutPathPrefix+"ca-key.pem",
			"-config=ca-config.json",
			"-hostname="+spec.Hostname,
			"-profile=kubernetes",
			"-").
			Command(
				"cfssljson",
				"-bare",
				OutPathPrefix+spec.OutFilePrefix))
}

var csrJSONTemplate = template.Must(template.New("").Parse(`{
  "CN": "{{.CN}}",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "US",
      "L": "Portland",
      "O": "{{.O}}",
      "OU": "CA",
      "ST": "Oregon"
    }
  ]
}
`))
