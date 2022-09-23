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
