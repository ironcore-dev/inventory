package flags

import (
	"path/filepath"

	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

type Flags struct {
	Verbose       bool
	Root          string
	Kubeconfig    string
	KubeNamespace string
	HTTPClient    bool
	Timeout       string
	Host          string
}

func NewFlags() *Flags {
	var kubeconfigDefaultPath string

	if home := homedir.HomeDir(); home != "" {
		kubeconfigDefaultPath = filepath.Join(home, ".kube", "config")
	}

	verbose := pflag.BoolP("verbose", "v", false, "verbose output")
	root := pflag.StringP("root", "r", "/", "path to root file system")
	kubeconfig := pflag.StringP("kubeconfig", "k", kubeconfigDefaultPath, "path to kubeconfig")
	kubeNamespace := pflag.StringP("namespace", "n", "default", "k8s namespace")
	httpClient := pflag.BoolP("gateway", "g", false, "use rest gateway for inventory creation")
	timeout := pflag.StringP("timeout", "t", "5s", "put timeout for client")
	host := pflag.StringP("host", "h", "http://localhost:8080", "inventory gateway")

	pflag.Parse()

	return &Flags{
		Verbose:       *verbose,
		Root:          *root,
		Kubeconfig:    *kubeconfig,
		KubeNamespace: *kubeNamespace,
		HTTPClient:    *httpClient,
		Timeout:       *timeout,
		Host:          *host,
	}
}
