package host

import (
	"os"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

type Info struct {
	Type string
	Name string
}

type Svc struct {
	printer  *printer.Svc
	hostType string
}

func NewSvc(printer *printer.Svc, hostType string) *Svc {
	return &Svc{
		printer:  printer,
		hostType: hostType,
	}
}

func (s *Svc) GetData() (*Info, error) {
	info := Info{}
	name, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hostname")
	}
	info.Name = name
	info.Type = s.hostType
	return &info, nil
}
