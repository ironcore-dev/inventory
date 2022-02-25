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

package command

import "github.com/urfave/cli/v2"

func checkFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "provider,p",
			Usage:       "Specify provider for cluster interaction. Example [bench-scheduler run -p kubernetes]",
			Aliases:     []string{"p"},
			EnvVars:     []string{"PROVIDER"},
			DefaultText: "http",
			Value:       "http",
		},
		&cli.StringFlag{
			Name:        "gateway,g",
			Usage:       "Specify http url for benchmark-gateway. Example [bench-scheduler run -g http://localhost:8080]",
			Aliases:     []string{"g"},
			EnvVars:     []string{"GATEWAY"},
			DefaultText: "http://localhost:8080",
			Value:       "http://localhost:8080",
		},
		&cli.StringFlag{
			Name:    "config,c",
			Usage:   "Specify config file with benchmarks. Example [bench-scheduler run -c examples/config.yaml]",
			Aliases: []string{"c"},
		},
		&cli.BoolFlag{
			Name:        "from-cluster-inventory",
			Usage:       "Set up when inventory data override is needed",
			DefaultText: "false",
			Value:       false,
		},
		&cli.StringFlag{
			Name:        "kubeconfig,k",
			Usage:       "Specify provider config file. Example for kubernetes [bench-scheduler run -k ~/.kube/config]",
			Aliases:     []string{"k"},
			EnvVars:     []string{"KUBECONFIG"},
			DefaultText: "~/.kube/config",
			Value:       "~/.kube/config",
		},
		&cli.StringFlag{
			Name:        "namespace,n",
			Usage:       "Specify provider config file. Example for kubernetes [bench-scheduler run -n your_wonderful_namespace]",
			Aliases:     []string{"n"},
			DefaultText: "onmetal",
			Value:       "onmetal",
		},
	}
}
