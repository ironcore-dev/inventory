package crd

import (
	"github.com/google/uuid"
	apiv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
)

type Capabilities []apiv1alpha1.LLDPCapabilities

func (c Capabilities) Len() int {
	return len(c)
}

func (c Capabilities) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Capabilities) Less(i, j int) bool {
	return c[i] < c[j]
}

func getUUID(namespace string, identifier string) string {
	namespaceUUID := uuid.NewMD5(uuid.UUID{}, []byte(namespace))
	newUUID := uuid.NewMD5(namespaceUUID, []byte(identifier))
	return newUUID.String()
}
