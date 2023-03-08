// Copyright 2023 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flags

import (
	"path/filepath"

	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

type InventoryFlags struct {
	Verbose       bool
	Root          string
	Kubeconfig    string
	KubeNamespace string
	Gateway       string
	Timeout       string
	Patch         bool
}

func NewInventoryFlags() *InventoryFlags {
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
	pflag.Parse()

	return &InventoryFlags{
		Verbose:       *verbose,
		Root:          *root,
		Kubeconfig:    *kubeconfig,
		KubeNamespace: *kubeNamespace,
		Gateway:       *gateway,
		Timeout:       *timeout,
		Patch:         *patch,
	}
}
