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
