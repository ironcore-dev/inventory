package distro

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/utils"
)

type Distro struct {
	BuildVersion  string
	DebianVersion string
	KernelVersion string
	AsicType      string
	CommitId      string
	BuildDate     string
	BuildNumber   uint32
	BuildBy       string
}

type Svc struct {
	printer *printer.Svc
}

func NewSvc(printer *printer.Svc) *Svc {
	return &Svc{
		printer: printer,
	}
}

func (s *Svc) GetData(hostType string) (*Distro, error) {
	distro := Distro{}
	rawInfo := make(map[string]interface{})
	switch hostType {
	case utils.CSwitchType:
		sonicInfo, err := ioutil.ReadFile(utils.CVersionFilePath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read SONiC version file")
		}
		err = yaml.Unmarshal(sonicInfo, &rawInfo)
		if err != nil {
			return nil, errors.Wrap(err, "failed to collect SONiC version")
		}
		err = convertMapStruct(&distro, rawInfo)
		if err != nil {
			return nil, errors.Wrap(err, "failed to process SONiC version")
		}
		// todo: case utils.CMachineType:
	}
	return &distro, nil
}

func convertMapStruct(obj *Distro, m map[string]interface{}) error {
	for k, v := range m {
		m[strings.Replace(k, "_", "", 1)] = v
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}
