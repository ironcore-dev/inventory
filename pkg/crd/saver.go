package crd

import apiv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"

type SaverSvc interface {
	Save(inv *apiv1alpha1.Inventory) error
	Patch(name string, patch interface{}) error
}
