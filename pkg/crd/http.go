package crd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/onmetal/inventory/pkg/inventory"
	apiv1alpha1 "github.com/onmetal/k8s-inventory/api/v1alpha1"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultHTTPTimeoutSecond = 30 * time.Second

type HttpClient struct {
	*http.Client

	host string
}

func newHttp(timeout, host string) Client {
	t, err := time.ParseDuration(timeout)
	if err != nil {
		log.Printf("can't parse given timeout: %s, setting default to 30s", err)
		t = defaultHTTPTimeoutSecond
	}

	if host == "" && os.Getenv("GATEWAY_HOST") != "" {
		host = os.Getenv("GATEWAY_HOST")
	}
	c := &http.Client{Timeout: t}
	return &HttpClient{c, host}
}

func (c *HttpClient) BuildAndSave(inv *inventory.Inventory) error {
	cr, err := Build(inv)
	if err != nil {
		return errors.Wrap(err, "unable to build inventory resource manifest")
	}

	if err := c.Save(cr); err != nil {
		return errors.Wrap(err, "unable to save inventory resource")
	}

	return nil
}

func (c *HttpClient) Save(inv *apiv1alpha1.Inventory) error {
	url := fmt.Sprintf("%s/api/v1/inventory", c.host)

	body, err := json.Marshal(inv)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	log.Println(string(respBody))

	return nil
}
