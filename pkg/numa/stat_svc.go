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

package numa

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CNodeStat = "/numastat"

	CNodeStatLinePattern = "^(\\w+)\\s+(\\d+)$"
)

var CNodeStatLineRegexp = regexp.MustCompile(CNodeStatLinePattern)

type StatSvc struct {
	printer *printer.Svc
}

func NewStatSvc(printer *printer.Svc) *StatSvc {
	return &StatSvc{
		printer: printer,
	}
}

func (s *StatSvc) GetStat(thePath string) (*Stat, error) {
	statPath := path.Join(thePath, CNodeStat)
	statData, err := os.ReadFile(statPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read stat file from path %s", statPath)
	}

	stat := &Stat{}

	bufReader := bytes.NewReader(statData)
	scanner := bufio.NewScanner(bufReader)
	for scanner.Scan() {
		line := scanner.Text()

		groups := CNodeStatLineRegexp.FindStringSubmatch(line)

		// groups [0] string [1] key [2] value
		if len(groups) < 3 {
			continue
		}

		key := groups[1]
		valString := groups[2]

		val, err := strconv.ParseUint(valString, 10, 64)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to parse %s:%s into uint64", key, valString))
			continue
		}

		err = stat.setField(key, val)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to set %s:%d", key, val))
		}
	}

	return stat, nil
}
