// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package flags

import (
	"path/filepath"

	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

type NICUpdaterFlags struct {
	Verbose       bool
	Root          string
	Kubeconfig    string
	KubeNamespace string
	Gateway       string
	Timeout       string
}

func NewNICUpdaterFlags() *NICUpdaterFlags {
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
	pflag.Parse()

	return &NICUpdaterFlags{
		Verbose:       *verbose,
		Root:          *root,
		Kubeconfig:    *kubeconfig,
		KubeNamespace: *kubeNamespace,
		Gateway:       *gateway,
		Timeout:       *timeout,
	}
}
