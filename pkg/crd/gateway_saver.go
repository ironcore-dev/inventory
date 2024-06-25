// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package crd

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/json"
)

const (
	CSaveRequestURLTemplate  = "%s/apis/v1alpha1/inventory/%s"
	CPatchRequestURLTemplate = "%s/apis/v1alpha1/inventory/%s/%s"

	CContentTypeHeader  = "Content-Type"
	CRequestContentType = "application/json"
)

type GatewaySaverSvc struct {
	httpClient      *http.Client
	host, namespace string
}

func NewGatewaySaverSvc(host, namespace string, timeout string) (SaverSvc, error) {
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse string %s as duration", timeout)
	}

	return GatewaySaverSvc{
		httpClient: &http.Client{
			Timeout: timeoutDuration,
		},
		host:      host,
		namespace: namespace,
	}, nil
}

func (g GatewaySaverSvc) Save(inv *metalv1alpha4.Inventory) error {
	url := fmt.Sprintf(CSaveRequestURLTemplate, g.host, g.namespace)

	body, err := json.Marshal(inv)
	if err != nil {
		return errors.Wrap(err, "unable to serialize inventory to json")
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "unable to form post request")
	}
	request.Header.Set(CContentTypeHeader, CRequestContentType)

	response, err := g.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "unable to post creation request")
	}
	defer response.Body.Close()

	return nil
}

func (g GatewaySaverSvc) Patch(name string, patch interface{}) error {
	url := fmt.Sprintf(CPatchRequestURLTemplate, g.host, g.namespace, name)

	body, err := json.Marshal(patch)
	if err != nil {
		return errors.Wrap(err, "unable to serialize inventory patch to json")
	}

	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "unable to form patch request")
	}
	request.Header.Set(CContentTypeHeader, CRequestContentType)

	response, err := g.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "unable to execute patch request")
	}
	defer response.Body.Close()

	return nil
}
