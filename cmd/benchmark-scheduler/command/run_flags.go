// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
