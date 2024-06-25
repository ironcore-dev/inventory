// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"

	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	bencherr "github.com/onmetal/inventory/internal/errors"
	"github.com/onmetal/inventory/internal/provider"
)

func New(machineUUID string, args *cli.Context, prv provider.Client, l logger.Logger) conf.Scheduler {
	configFile, err := getConfig(machineUUID, args.String("config"), args.Bool("from-cluster-inventory"), prv)
	if err != nil {
		l.Info("can't read config file", "error", err)
		os.Exit(1)
	}

	var s conf.Scheduler
	if err := yaml.Unmarshal(configFile, &s); err != nil {
		l.Info("can't unmarshal config from yaml to json", "error", err)
		os.Exit(1)
	}
	return s
}

func getConfig(machineUUID, fromArgs string, fromClusterInventory bool, prv provider.Client) ([]byte, error) {
	switch {
	case fromArgs != "" || os.Getenv("CONFIG_PATH") != "":
		path := getConfigPath(fromArgs)
		configFile, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
		if fromClusterInventory {
			return prv.GenerateConfig(machineUUID, configFile)
		}
		return configFile, err
	case fromArgs == "":
		return prv.Get(machineUUID, "config")
	default:
		return nil, bencherr.NotFound("config")
	}
}

func getConfigPath(fromArgs string) string {
	if fromArgs != "" {
		return fromArgs
	}
	if os.Getenv("CONFIG_PATH") != "" {
		return os.Getenv("CONFIG_PATH")
	}
	return ""
}
