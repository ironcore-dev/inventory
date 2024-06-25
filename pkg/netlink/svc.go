// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package netlink

import (
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	"github.com/onmetal/inventory/pkg/chroot"
	"github.com/onmetal/inventory/pkg/printer"
)

type Svc struct {
	printer  *printer.Svc
	rootPath string
}

func NewSvc(printer *printer.Svc, basePath string) *Svc {
	return &Svc{
		printer:  printer,
		rootPath: basePath,
	}
}

func (s *Svc) GetIPv6NeighbourData() ([]IPv6Neighbour, error) {
	chr, err := chroot.New(s.rootPath)
	if err != nil {
		s.printer.VErr(errors.Errorf("got error on chroot to %s, will try to collect data without it", s.rootPath))
	}
	defer func() {
		// Not sure if it is best to test for err != nil or chr == nil
		if chr == nil {
			s.printer.VErr(errors.Wrapf(err, "unable to create chroot"))
			return
		}
		if err := chr.Close(); err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to exit chroot"))
		}
	}()

	ll, err := netlink.LinkList()
	if err != nil {
		return nil, errors.Wrap(err, "unable to obtain device list")
	}

	neighbours := make([]IPv6Neighbour, 0)
	for _, l := range ll {
		iIdx := l.Attrs().Index
		iName := l.Attrs().Name
		nl, err := netlink.NeighList(iIdx, unix.AF_INET6)
		if err != nil {
			s.printer.VErr(errors.Wrapf(err, "unable to get neighbours for %s", iName))
			continue
		}

		for _, n := range nl {
			neighbour := NewIPv6Neighbour(iIdx, iName, &n)
			neighbours = append(neighbours, *neighbour)
		}
	}

	return neighbours, nil
}
