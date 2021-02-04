package ipmi

import (
	"io/ioutil"
	"path"
	"regexp"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

const (
	CDevPath        = "/dev"
	CIPMIDevPattern = "ipmi\\d+"
)

var CIPMIDevRegexp = regexp.MustCompile(CIPMIDevPattern)

type Svc struct {
	printer     *printer.Svc
	ipmiInfoSvc *DeviceSvc
	devPath     string
}

func NewSvc(printer *printer.Svc, ipmiDevInfoSvc *DeviceSvc, basePath string) *Svc {
	return &Svc{
		printer:     printer,
		ipmiInfoSvc: ipmiDevInfoSvc,
		devPath:     path.Join(basePath, CDevPath),
	}
}

func (s *Svc) GetData() ([]Device, error) {
	devFolderContents, err := ioutil.ReadDir(s.devPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read contents of %s", s.devPath)
	}

	infos := make([]Device, 0)
	for _, dev := range devFolderContents {
		devName := dev.Name()

		matches := CIPMIDevRegexp.MatchString(devName)

		if !matches {
			continue
		}

		thePath := path.Join(s.devPath, devName)
		info, err := s.ipmiInfoSvc.GetDevice(thePath)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unabale to obtain IPMI device info"))
			continue
		}

		infos = append(infos, *info)
	}

	return infos, nil
}
