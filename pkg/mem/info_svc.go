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

package mem

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CProcMemInfoPath = "/proc/meminfo"

	// meminfo measures in kibibytes
	CMemInfoValueMultiplier = 1024

	// ^ - begin line
	// (\w+\s\d+\s)? - optional group to parse NUMA prefix "Node 0 "
	// ([\w\(\)]+) - property key group
	// : - colon that separates key from value
	// \s* - whitespace between colon and value
	// (\d+) - numerical value
	// (\s\w*)? - optional measurement unit identifier "kB"
	// $ - end line
	CMemInfoLinePattern = "^(\\w+\\s\\d+\\s)?([\\w\\(\\)]+):\\s*(\\d+)(\\s\\w*)?$"
)

var CMemInfoLineRegexp = regexp.MustCompile(CMemInfoLinePattern)

type InfoSvc struct {
	printer     *printer.Svc
	memInfoPath string
}

func NewInfoSvc(printer *printer.Svc, basePath string) *InfoSvc {
	return &InfoSvc{
		printer:     printer,
		memInfoPath: path.Join(basePath, CProcMemInfoPath),
	}
}

func (s *InfoSvc) GetInfo() (*Info, error) {
	return s.GetInfoFromFile(s.memInfoPath)
}

func (s *InfoSvc) GetInfoFromFile(thePath string) (*Info, error) {
	memInfoData, err := os.ReadFile(thePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read meminfo from %s", thePath)
	}

	mem := &Info{}

	bufReader := bytes.NewReader(memInfoData)
	scanner := bufio.NewScanner(bufReader)
	for scanner.Scan() {
		line := scanner.Text()

		groups := CMemInfoLineRegexp.FindStringSubmatch(line)

		// should contain 5 groups according to regexp
		// [0] self; [1] NUMA prefix; [2] key; [3] value; [4] measurement unit
		if len(groups) < 5 {
			continue
		}

		key := groups[2]
		valString := groups[3]

		val, err := strconv.ParseUint(valString, 10, 64)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to parse %s:%s into uint64", key, valString))
			continue
		}

		// check if measurement unit is applied to the value
		// if applied, multiply to get bytes
		if strings.TrimSpace(groups[4]) != "" {
			val = val * CMemInfoValueMultiplier
		}

		err = mem.setField(key, val)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to set %s:%d", key, val))
		}
	}

	return mem, nil
}
