package host

import (
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/utils"
)

type Info struct {
	Type string
	Name string
}

type Svc struct {
	printer           *printer.Svc
	switchVersionPath string
}

func NewSvc(printer *printer.Svc, basePath string) *Svc {
	return &Svc{
		printer:           printer,
		switchVersionPath: path.Join(basePath, utils.CVersionFilePath),
	}
}

func (s *Svc) GetData() (*Info, error) {
	hostType, err := utils.GetHostType(s.switchVersionPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine host type")
	}

	info := Info{}
	name, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hostname")
	}
	info.Name = name
	info.Type = hostType
	return &info, nil
}
