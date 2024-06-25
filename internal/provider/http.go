// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	nethttp "net/http"
	"os"
	"strconv"
	"time"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	bencherr "github.com/onmetal/inventory/internal/errors"
)

const defaultTimeoutSecond = 60

const (
	basePatchURL      = "apis/v1alpha3/benchmark"
	getURL            = "apis/v1alpha3/benchmark"
	getConfigURL      = "apis/v1alpha3/benchmark/config"
	generateConfigURL = "apis/v1alpha3/benchmark/config/generate"
)

type httpClient struct {
	*nethttp.Client

	ctx                context.Context
	log                logger.Logger
	gateway, namespace string
}

func newHTTP(ctx context.Context, l logger.Logger, gateway, namespace string) (Client, error) {
	if namespace == "" {
		namespace = "default"
	}
	c := &nethttp.Client{Timeout: getClientTimeoutSecond()}
	return &httpClient{
		Client:    c,
		gateway:   gateway,
		namespace: namespace,
		ctx:       ctx,
		log:       l,
	}, nil
}

func getClientTimeoutSecond() time.Duration {
	timeout := defaultTimeoutSecond
	if os.Getenv("HTTP_CLIENT_TIMEOUT_SECOND") != "" {
		t, err := strconv.Atoi(os.Getenv("HTTP_CLIENT_TIMEOUT"))
		if err != nil {
			log.Printf("can't convert client timeout: %s", err)
		} else {
			timeout = t
		}
	}
	return time.Duration(timeout) * time.Second
}

func (s *httpClient) Patch(machineUUID string, body []byte) error {
	req, err := s.prepareRequestWithBody(machineUUID, basePatchURL, nethttp.MethodPatch, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.Do(req) //nolint:bodyclose //reason: false alarm. it's closed in defer
	if err != nil {
		return err
	}
	defer func(body io.ReadCloser) {
		if bodyErr := body.Close(); bodyErr != nil {
			s.log.Info("failed to close body", "error", err)
		}
	}(resp.Body)
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != nethttp.StatusOK {
		return bencherr.UpdateError(resp.StatusCode, string(respbody))
	}
	s.log.Info("patched successfully", "body", string(respbody), "status code", resp.StatusCode)
	return nil
}

func (s *httpClient) Get(uuid, kind string) ([]byte, error) {
	switch kind {
	case "benchmark":
		req, err := s.prepareRequestWithoutBody(uuid, getURL, nethttp.MethodGet)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.Do(req) //nolint:bodyclose //reason: false alarm. it's closed
		if err != nil {
			return nil, err
		}
		defer func(body io.ReadCloser) {
			if bodyErr := body.Close(); bodyErr != nil {
				s.log.Info("failed to close body", "error", err)
			}
		}(resp.Body)
		if resp.StatusCode != nethttp.StatusOK {
			return nil, bencherr.GetConfigError(strconv.Itoa(resp.StatusCode))
		}
		return io.ReadAll(resp.Body)
	case "config":
		req, err := s.prepareRequestWithoutBody(uuid, getConfigURL, nethttp.MethodPost)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.Do(req) //nolint:bodyclose //reason: false alarm. it's closed
		if err != nil {
			return nil, err
		}
		defer func(body io.ReadCloser) {
			if bodyErr := body.Close(); bodyErr != nil {
				s.log.Info("failed to close body", "error", err)
			}
		}(resp.Body)
		if resp.StatusCode != nethttp.StatusOK {
			return nil, bencherr.GetConfigError(strconv.Itoa(resp.StatusCode))
		}
		return io.ReadAll(resp.Body)
	default:
		return nil, bencherr.NotSupported(kind)
	}
}

func (s *httpClient) GenerateConfig(machineUUID string, config []byte) ([]byte, error) {
	req, err := s.prepareRequestWithBody(machineUUID, generateConfigURL, nethttp.MethodPost, config)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.Do(req) //nolint:bodyclose //reason: false alarm. it's closed
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		if bodyErr := body.Close(); bodyErr != nil {
			s.log.Info("failed to close body", "error", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != nethttp.StatusOK {
		return nil, bencherr.GenerateConfigError(string(body))
	}
	return body, err
}

func (s *httpClient) prepareRequestWithBody(machineUUID, url, method string, body []byte) (*nethttp.Request, error) {
	uri := fmt.Sprintf("%s/%s/%s/%s", s.gateway, url, s.namespace, machineUUID)
	return nethttp.NewRequestWithContext(s.ctx, method, uri, bytes.NewBuffer(body))
}

func (s *httpClient) prepareRequestWithoutBody(machineUUID, url, method string) (*nethttp.Request, error) {
	uri := fmt.Sprintf("%s/%s/%s/%s", s.gateway, url, s.namespace, machineUUID)
	return nethttp.NewRequestWithContext(s.ctx, method, uri, nethttp.NoBody)
}
