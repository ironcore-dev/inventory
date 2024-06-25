// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"

	"github.com/onmetal/inventory/cmd/benchmark-scheduler/command"
)

var VERSION = "dev"

func main() {
	app := command.NewRoot(VERSION)
	if err := app.Run(os.Args); err != nil {
		log.Println("application exited not normally.", "error:", err)
		os.Exit(1)
	}
}
