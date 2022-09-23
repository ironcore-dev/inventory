package frame

import (
	"os"

	"github.com/mdlayher/lldp"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/printer"
)

type Svc struct {
	printer *printer.Svc
}

func NewFrameSvc(printer *printer.Svc) *Svc {
	return &Svc{
		printer: printer,
	}
}

func (s *Svc) GetFrame(interfaceID string, thePath string) (*Frame, error) {
	contents, err := os.ReadFile(thePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read file %s", thePath)
	}

	// 1-8 bytes - ?
	// 9-22 bytes - ethernet frame part
	// 23-rest - LLDP frame part
	frame := lldp.Frame{}
	err = frame.UnmarshalBinary(contents[22:])
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal LLDP frame")
	}

	frameInfo := &Frame{
		InterfaceID: interfaceID,
		TTL:         frame.TTL,
	}

	err = frameInfo.setChassisID(frame.ChassisID)
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to set chassis ID"))
	}

	err = frameInfo.setPortID(frame.PortID)
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to unmarshal port ID"))
	}

	for _, tlv := range frame.Optional {
		err = frameInfo.setOptional(tlv)
		if err != nil {
			s.printer.VErr(errors.Wrap(err, "unable to process optional TLV"))
		}
	}

	return frameInfo, nil
}
