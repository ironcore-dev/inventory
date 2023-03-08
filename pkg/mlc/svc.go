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

package mlc

import (
	"bufio"
	"bytes"
	"os/exec"
	"path"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CMlcPath = "/usr/sbin/mlc"

	// meminfo measures in kibibytes
	CMlcValueMultiplier = 1024

	/*
	   Example output:
	   Intel(R) Memory Latency Checker - v3.9
	   Command line parameters: --bandwidth_matrix

	   Using buffer size of 100.000MiB/thread for reads and an additional 100.000MiB/thread for writes
	   Measuring Memory Bandwidths between nodes within system
	   Bandwidths are in MB/sec (1 MB/sec = 1,000,000 Bytes/sec)
	   Using all the threads from each core if Hyper-threading is enabled
	   Using Read-only traffic type
	                   Numa node
	   Numa node            0       1
	          0        43179.2 39972.6
	          1        39917.1 43281.2
	*/
	/*
	   Intel(R) Memory Latency Checker - v3.9
	   Command line parameters: --latency_matrix

	   Using buffer size of 2000.000MiB
	   Measuring idle latencies (in ns)...
	                   Numa node
	   Numa node            0       1
	          0          64.2   120.3
	          1         120.3    64.2
	*/
	// Capture the digits for the numa nodes, not the whitespaces in front of it
	CMlcNumaNodePattern    = "^Numa node(?:\\s+(\\d+))+\\s*$"
	CMlcBWLinePattern      = "^\\s*\\d(?:\\s*((?:[0-9]*[.])?[0-9]+)+)\\s*$"
	CMlcLatencyLinePattern = "^\\s*\\d(?:\\s*((?:[0-9]*[.])?[0-9]+)+)\\s*$"
)

var CMlcNumaNodeRegexp = regexp.MustCompile(CMlcNumaNodePattern)
var CMlcBWLineRegexp = regexp.MustCompile(CMlcBWLinePattern)
var CMlcLatencyLineRegexp = regexp.MustCompile(CMlcLatencyLinePattern)

type PerfSvc struct {
	printer *printer.Svc
	mlcPath string
}

func NewPerfSvc(printer *printer.Svc, basePath string) *PerfSvc {
	return &PerfSvc{
		printer: printer,
		mlcPath: path.Join(basePath, CMlcPath),
	}
}

func (s *PerfSvc) GetInfo() (*Perf, error) {
	return s.GetInfoFromMlc(s.mlcPath)
}

func (s *PerfSvc) GetInfoFromMlc(mlcPath string) (*Perf, error) {

	_, err := exec.Command("cpupower", "frequency-set", "-g", "performance").Output()
	if err != nil {
		return nil, errors.Wrapf(err, "MLC: unable to execute cpower frequency-set -g performance, not root or not installed?")
	}

	// Exec mlc here and then read from stdout instead
	mlcLatencyOutput, err := exec.Command(mlcPath, "--latency_matrix").Output()
	if err != nil {
		return nil, errors.Wrapf(err, "MLC: unable to read mlc latency output using %s", mlcPath)
	}

	mlcBWOutput, err := exec.Command(mlcPath, "--bandwidth_matrix").Output()
	if err != nil {
		return nil, errors.Wrapf(err, "MLC: unable to read mlc BW output using %s", mlcPath)
	}

	mlcPerf := &Perf{}

	bufReaderLatency := bytes.NewReader(mlcLatencyOutput)
	scannerLatency := bufio.NewScanner(bufReaderLatency)

	var nodes = 0

	for scannerLatency.Scan() {
		line := scannerLatency.Text()

		groups := CMlcNumaNodeRegexp.FindStringSubmatch(line)

		// [0] self; [1..X] NUMA node numbers;
		if len(groups) == 0 {
			continue
		}

		nodes = len(groups) - 1
		break
	}

	if nodes == 0 {
		return nil, errors.Wrapf(err, "Cannot determine Numa nodes in mlc BW output")
	}
	for scannerLatency.Scan() {
		line := scannerLatency.Text()
		groups := CMlcLatencyLineRegexp.FindStringSubmatch(line)

		if len(groups) == 0 {
			continue
		}

		// groups[0] is the entire capture, groups[1] the local, groups[2+] if present remote
		// more remotes should be captured in case of not fully meshed interconnects but requires numa distances to determine
		if len(groups) > 1 {
			localLatencyString := groups[1]
			localLatencyVal, err := strconv.ParseFloat(localLatencyString, 64)
			if err != nil {
				s.printer.VErr(errors.Wrapf(err, "unable to parse %s:%s into uint64", CMemPerfLocalMemLatencyKey, localLatencyString))
			}

			err = mlcPerf.setField(CMemPerfLocalMemLatencyKey, localLatencyVal)
			if err != nil {
				s.printer.VErr(errors.Wrapf(err, "unable to set %s:%f", CMemPerfLocalMemLatencyKey, localLatencyVal))
			}
		}
		if len(groups) > 2 {
			remoteLatencyString := groups[2]
			remoteLatencyVal, err := strconv.ParseFloat(remoteLatencyString, 64)
			if err != nil {
				s.printer.VErr(errors.Wrapf(err, "unable to parse %s:%s into uint64", CMemPerfRemoteMemLatencyKey, remoteLatencyString))
			}

			err = mlcPerf.setField(CMemPerfRemoteMemLatencyKey, remoteLatencyVal)
			if err != nil {
				s.printer.VErr(errors.Wrapf(err, "unable to set %s:%f", CMemPerfRemoteMemLatencyKey, remoteLatencyVal))
			}
		}
	}

	bufReaderBW := bytes.NewReader(mlcBWOutput)
	scannerBW := bufio.NewScanner(bufReaderBW)

	nodes = 0

	for scannerBW.Scan() {
		line := scannerBW.Text()

		groups := CMlcNumaNodeRegexp.FindStringSubmatch(line)

		// [0] self; [1..X] NUMA node numbers;
		if len(groups) == 0 {
			continue
		}

		nodes = len(groups) - 1
		break
	}

	if nodes == 0 {
		return nil, errors.Wrapf(err, "Cannot determine Numa nodes in mlc BW output")
	}
	for scannerBW.Scan() {
		line := scannerBW.Text()
		groups := CMlcBWLineRegexp.FindStringSubmatch(line)

		if len(groups) == 0 {
			continue
		}

		if len(groups) > 1 {
			localBWString := groups[1]
			localBWVal, err := strconv.ParseFloat(localBWString, 64)
			if err != nil {
				s.printer.VErr(errors.Wrapf(err, "unable to parse %s:%s into uint64", CMemPerfLocalMemBWKey, localBWString))
			}

			err = mlcPerf.setField(CMemPerfLocalMemBWKey, localBWVal)
			if err != nil {
				s.printer.VErr(errors.Wrapf(err, "unable to set %s:%f", CMemPerfLocalMemBWKey, localBWVal))
			}
		}
		if len(groups) > 2 {
			remoteBWString := groups[2]
			remoteBWVal, err := strconv.ParseFloat(remoteBWString, 64)
			if err != nil {
				s.printer.VErr(errors.Wrapf(err, "unable to parse %s:%s into uint64", CMemPerfRemoteMemBWKey, remoteBWString))
			}

			err = mlcPerf.setField(CMemPerfRemoteMemBWKey, remoteBWVal)
			if err != nil {
				s.printer.VErr(errors.Wrapf(err, "unable to set %s:%f", CMemPerfRemoteMemBWKey, remoteBWVal))
			}
		}
	}

	return mlcPerf, nil
}
