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