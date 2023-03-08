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
