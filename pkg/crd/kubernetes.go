package crd

import (
	"context"
	apiv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	clientv1alpha1 "github.com/onmetal/k8s-inventory/clientset/v1alpha1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/onmetal/inventory/pkg/inventory"
)


type Kubernetes struct {
	client clientv1alpha1.InventoryInterface
}

func newKubernetes(kubeconfig string, namespace string) (Client, error) {
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

	return &Kubernetes{
		client: client,
	}, nil
}

func (s *Kubernetes) BuildAndSave(inv *inventory.Inventory) error {
	cr, err := Build(inv)
	if err != nil {
		return errors.Wrap(err, "unable to build inventory resource manifest")
	}

	if err := s.Save(cr); err != nil {
		return errors.Wrap(err, "unable to save inventory resource")
	}

	return nil
}

func (s *Kubernetes) Save(inv *apiv1alpha1.Inventory) error {
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

	if existing.Labels == nil {
		existing.Labels = inv.Labels
	} else {
		for k, v := range inv.Labels {
			existing.Labels[k] = v
		}
	}

	if _, err := s.client.Update(context.Background(), existing, metav1.UpdateOptions{}); err != nil {
		return errors.Wrap(err, "unhandled error on update")
	}

	return nil
}