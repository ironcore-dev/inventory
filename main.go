// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/printer"
)

func main() {
	p := printer.NewSvc(true)

	rawDmiSvc := dmi.NewRawSvc("/")
	sm := dmi.NewSvc(p, rawDmiSvc)
	data, err := sm.GetData()
	if err != nil {
		p.Err(err)
		os.Exit(1)
	}
	fmt.Printf("uuid: %s", data.SystemInformation.UUID)
}
