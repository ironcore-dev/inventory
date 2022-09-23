package flags

import (
	"os"
	"path/filepath"

	"k8s.io/utils/pointer"

	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

type InventoryFlags struct {
	Verbose           bool
	Root              string
	Kubeconfig        string
	KubeNamespace     string
	Gateway           string
	Timeout           string
	Patch             bool
	RedisUser         string
	RedisPassword     string
	RedisPasswordFile string
}

func NewInventoryFlags() (*InventoryFlags, error) {
	var kubeconfigDefaultPath string

	if home := homedir.HomeDir(); home != "" {
		kubeconfigDefaultPath = filepath.Join(home, ".kube", "config")
	}

	verbose := pflag.BoolP("verbose", "v", false, "verbose output")
	root := pflag.StringP("root", "r", "/", "path to root file system")
	kubeconfig := pflag.StringP("kubeconfig", "k", kubeconfigDefaultPath, "path to kubeconfig")
	kubeNamespace := pflag.StringP("namespace", "n", "default", "k8s namespace")
	gateway := pflag.StringP("gateway", "g", "", "gateway address")
	timeout := pflag.StringP("timeout", "t", "30s", "request timeout, if gateway is used")
	patch := pflag.BoolP("patch", "p", false, "patch crd object instead of creation")
	redisUser := pflag.String("redis-user", "", "redis user")
	redisPassword := pflag.String("redis-password", "", "redis password")
	redisPasswordFile := pflag.String("redis-password-file", "", "redis password file")
	pflag.Parse()

	if *redisPasswordFile != "" {
		passwordFromFile, err := os.ReadFile(*redisPasswordFile)
		if err != nil {
			return nil, err
		}

		if len(passwordFromFile) > 0 {
			redisPassword = pointer.String(string(passwordFromFile))
		}
	}

	return &InventoryFlags{
		Verbose:       *verbose,
		Root:          *root,
		Kubeconfig:    *kubeconfig,
		KubeNamespace: *kubeNamespace,
		Gateway:       *gateway,
		Timeout:       *timeout,
		Patch:         *patch,
		RedisUser:     *redisUser,
		RedisPassword: *redisPassword,
	}, nil
}
