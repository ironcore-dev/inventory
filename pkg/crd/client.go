package crd

import (
	"github.com/onmetal/inventory/pkg/flags"
	"github.com/onmetal/inventory/pkg/inventory"
)

type Client interface {
	BuildAndSave(inv *inventory.Inventory) error
}

func NewClient(f *flags.Flags) (Client, error) {
	if f.HTTPClient {
		return newHttp(f.Timeout, f.Host), nil
	}
	return newKubernetes(f.Kubeconfig, f.KubeNamespace)
}

