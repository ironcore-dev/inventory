// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/onmetal/inventory/pkg/app"
)

func main() {
	appInstance, ret := app.NewInventoryApp()
	if ret != 0 {
		os.Exit(ret)
	}
	ret = appInstance.Run()
	os.Exit(ret)
}
