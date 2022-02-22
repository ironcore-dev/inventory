// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package config

import (
	"os"
	"path/filepath"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/logger"
	bencherr "github.com/onmetal/inventory/internal/errors"
	"github.com/onmetal/inventory/internal/provider"
	conf "github.com/onmetal/metal-api-gateway/app/handlers/benchmark"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func New(machineUUID string, args *cli.Context, prv provider.Client, l logger.Logger) conf.Scheduler {
	configFile, err := getConfig(machineUUID, args.String("config"), prv)
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

func getConfig(machineUUID, fromArgs string, prv provider.Client) ([]byte, error) {
	switch {
	case os.Getenv("CONFIG_PATH") != "":
		configFile, err := os.ReadFile(filepath.Clean(os.Getenv("CONFIG_PATH")))
		if err != nil {
			return nil, err
		}
		generated, genErr := prv.GenerateConfig(machineUUID, configFile)
		if genErr != nil {
			return configFile, nil
		}
		return generated, genErr
	case fromArgs != "":
		configFile, err := os.ReadFile(filepath.Clean(fromArgs))
		if err != nil {
			return nil, err
		}
		generated, genErr := prv.GenerateConfig(machineUUID, configFile)
		if genErr != nil {
			return configFile, nil
		}
		return generated, genErr
	case fromArgs == "":
		return prv.Get(machineUUID, "config")
	default:
		return nil, bencherr.NotFound("config")
	}
}
