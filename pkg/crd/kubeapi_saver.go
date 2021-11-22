package crd

import (
	"context"
	"fmt"

	apiv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	clientv1alpha1 "github.com/onmetal/k8s-inventory/clientset/v1alpha1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	CSonicNamespace = "switch.onmetal.de"
)

type KubeAPISaverSvc struct {
	client clientv1alpha1.InventoryInterface
}

func NewKubeAPISaverSvc(kubeconfig string, namespace string) (SaverSvc, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read kubeconfig from path %s", kubeconfig)
	}

	if err := apiv1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return nil, errors.Wrap(err, "unable to add registered types to client scheme")
	}

	clientset, err := clientv1alpha1.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build clientset from config")
	}

	client := clientset.Inventories(namespace)

	return &KubeAPISaverSvc{
		client: client,
	}, nil
}

func (s *KubeAPISaverSvc) Save(inv *apiv1alpha1.Inventory) error {
	_, err := s.client.Create(context.Background(), inv, metav1.CreateOptions{})
	if err == nil {
		return nil
	}
	if !apierrors.IsAlreadyExists(err) {
		return errors.Wrap(err, "unhandled error on creation")
	}

	existing, err := s.client.Get(context.Background(), inv.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to get resource")
	}

	existing.Spec = inv.Spec

	if _, err := s.client.Update(context.Background(), existing, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(err, "unhandled error on update")
	}

	return nil
}

func (s *KubeAPISaverSvc) Patch(name string, patch interface{}) error {
	patchData, err := json.Marshal(patch)
	if err != nil {
		return errors.Wrap(err, "unable to marshal inventory")
	}
	fmt.Println(string(patchData))
	_, err = s.client.Patch(context.Background(), name, types.MergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to patch inventory")
	}
	return nil
}
