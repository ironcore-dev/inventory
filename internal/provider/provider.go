// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
)

const (
	HTTP = "http"
)

type Client interface {
	Get(name, kind string) ([]byte, error)
	GenerateConfig(name string, config []byte) ([]byte, error)
	Patch(name string, body []byte) error
}

func New(ctx context.Context, l logger.Logger, cliCtx *cli.Context) (Client, error) {
	prv := getFrom(os.Getenv("PROVIDER"), cliCtx.String("provider"))
	gateway := cliCtx.String("gateway")
	namespace := cliCtx.String("namespace")
	switch prv {
	case HTTP:
		return newHTTP(ctx, l, gateway, namespace)
	default:
		l.Info("provider not found. default http returned", "name", prv)
		return newHTTP(ctx, l, gateway, namespace)
	}
}

func getFrom(fromEnv, fromArgs string) string {
	switch {
	case fromEnv != "":
		return fromEnv
	case fromArgs != "":
		return fromArgs
	default:
		return HTTP
	}
}
