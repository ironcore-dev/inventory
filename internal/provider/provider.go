package provider

import (
	"context"
	"os"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	"github.com/urfave/cli/v2"
)

const (
	HTTP = "http"
)

type Client interface {
	Get(name, kind string) ([]byte, error)
	GenerateConfig(name string, config []byte) ([]byte, error)
	Patch(name, namespace string, body []byte) error
}

func New(ctx context.Context, l logger.Logger, cliCtx *cli.Context) (Client, error) {
	prv := getFrom(os.Getenv("PROVIDER"), cliCtx.String("provider"))
	switch prv {
	case HTTP:
		return newHTTP(ctx, l, cliCtx.String("gateway"))
	default:
		l.Info("provider not found. default http returned", "name", prv)
		return newHTTP(ctx, l, cliCtx.String("gateway"))
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
