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
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CNodeDevicePath = "/sys/devices/system/node"

	CNumericNodeDeviceDirNamePattern = "node([0-9]+)"
)

var CNumericNodeDeviceDirNameRegexp = regexp.MustCompile(CNumericNodeDeviceDirNamePattern)

type Svc struct {
	printer        *printer.Svc
	nodeSvc        *NodeSvc
	nodeDevicePath string
}

func NewSvc(printer *printer.Svc, nodeSvc *NodeSvc, basePath string) *Svc {
	return &Svc{
		printer:        printer,
		nodeSvc:        nodeSvc,
		nodeDevicePath: path.Join(basePath, CNodeDevicePath),
	}
}

func (s *Svc) GetData() ([]Node, error) {
	numaFolders, err := os.ReadDir(s.nodeDevicePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of numa node devices")
	}

	numaNodes := make([]Node, 0)
	for _, numaFolder := range numaFolders {
		name := numaFolder.Name()

		if !numaFolder.IsDir() {
			continue
		}

		groups := CNumericNodeDeviceDirNameRegexp.FindStringSubmatch(name)

		// String itself is always a first element in results
		// so we need at least 2 to get number from our group
		if len(groups) < 2 {
			continue
		}

		nodeNumberString := groups[1]
		nodeNumber, err := strconv.Atoi(nodeNumberString)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to convert node number string %s to int", nodeNumberString))
			continue
		}

		nodePath := path.Join(s.nodeDevicePath, name)
		node, err := s.nodeSvc.GetNode(nodePath, nodeNumber)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to collect  %s", nodePath))
			continue
		}

		numaNodes = append(numaNodes, *node)
	}

	return numaNodes, nil
}
