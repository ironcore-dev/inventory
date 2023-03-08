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

package cpu

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CCPUInfoPath = "/proc/cpuinfo"

	CCPUInfoLinePattern = "^(\\w+\\s?\\w+?)\\s*:\\s*(.*)$"
)

var CCPUInfoLineRegexp = regexp.MustCompile(CCPUInfoLinePattern)

type InfoSvc struct {
	printer     *printer.Svc
	cpuInfoPath string
}

func NewInfoSvc(printer *printer.Svc, basePath string) *InfoSvc {
	return &InfoSvc{
		printer:     printer,
		cpuInfoPath: path.Join(basePath, CCPUInfoPath),
	}
}

func (s *InfoSvc) GetInfo() ([]Info, error) {
	cpuInfoData, err := os.ReadFile(s.cpuInfoPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read cpuinfo from %s", s.cpuInfoPath)
	}

	cpus := make([]Info, 0)
	cpu := Info{}

	bufReader := bytes.NewReader(cpuInfoData)
	scanner := bufio.NewScanner(bufReader)
	for scanner.Scan() {
		line := scanner.Text()

		// cpu records are separated with empty line
		if strings.TrimSpace(line) == "" {
			cpus = append(cpus, cpu)
			cpu = Info{}
		}

		groups := CCPUInfoLineRegexp.FindStringSubmatch(line)

		// should contain 3 groups according to regexp
		// [0] self; [1] key; [2] value
		if len(groups) < 3 {
			continue
		}

		key := groups[1]
		val := groups[2]

		err = cpu.setField(key, val)

		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to set field %s with value %s", key, val))
		}
	}

	return cpus, nil
}
