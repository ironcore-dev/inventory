package host

import (
	"os"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/utils"
)

type Info struct {
	Type string
	Name string
}

type Svc struct {
	printer  *printer.Svc
	hostType string
}

func NewSvc(printer *printer.Svc) *Svc {
	return &Svc{
		printer: printer,
	}
}

func (s *Svc) GetData() (*Info, error) {
	hostType, err := utils.GetHostType()
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
