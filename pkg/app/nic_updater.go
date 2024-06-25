// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package app

import (
	metalv1alpha4 "github.com/ironcore-dev/metal/apis/metal/v1alpha4"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/crd"
	"github.com/onmetal/inventory/pkg/dmi"
	"github.com/onmetal/inventory/pkg/flags"
	"github.com/onmetal/inventory/pkg/gatherer"
	"github.com/onmetal/inventory/pkg/host"
	"github.com/onmetal/inventory/pkg/inventory"
	"github.com/onmetal/inventory/pkg/lldp"
	"github.com/onmetal/inventory/pkg/lldp/frame"
	"github.com/onmetal/inventory/pkg/netlink"
	"github.com/onmetal/inventory/pkg/nic"
	"github.com/onmetal/inventory/pkg/printer"
	"github.com/onmetal/inventory/pkg/redis"
)

type NICUpdaterApp struct {
	printer       *printer.Svc
	gathererSvc   *gatherer.Svc
	crdBuilderSvc *crd.BuilderSvc
	crdSaverSvc   crd.SaverSvc
}

func NewNICUpdaterApp() (*NICUpdaterApp, int) {
	f := flags.NewNICUpdaterFlags()
	p := printer.NewSvc(f.Verbose)

	crdBuilderSvc := crd.NewBuilderSvc(p)

	var crdSvcConstructor func() (crd.SaverSvc, error)
	if f.Gateway != "" {
		crdSvcConstructor = func() (crd.SaverSvc, error) {
			return crd.NewGatewaySaverSvc(f.Gateway, f.KubeNamespace, f.Timeout)
		}
	} else {
		crdSvcConstructor = func() (crd.SaverSvc, error) {
			return crd.NewKubeAPISaverSvc(f.Kubeconfig, f.KubeNamespace)
		}
	}

	crdSaverSvc, err := crdSvcConstructor()
	if err != nil {
		p.Err(errors.Wrapf(err, "unable to create k8s resorce saver svc"))
		return nil, CErrRetCode
	}

	rawDmiSvc := dmi.NewRawSvc(f.Root)
	dmiSvc := dmi.NewSvc(p, rawDmiSvc)

	hostSvc := host.NewSvc(p, f.Root)

	redisSvc, err := redis.NewRedisSvc(f.Root)
	if err != nil {
		p.Err(errors.Wrapf(err, "unable to init redis client"))
		return nil, CErrRetCode
	}

	lldpFrameInfoSvc := frame.NewFrameSvc(p)
	lldpSvc := lldp.NewSvc(p, lldpFrameInfoSvc, hostSvc, redisSvc, f.Root)

	nlSvc := netlink.NewSvc(p, f.Root)

	nicDevSvc := nic.NewDeviceSvc(p)
	nicSvc := nic.NewSvc(p, nicDevSvc, hostSvc, redisSvc, f.Root)

	opts := []gatherer.Option{
		gatherer.WithDMI(dmiSvc),
		gatherer.WithLLDP(lldpSvc),
		gatherer.WithNIC(nicSvc),
		gatherer.WithNetlink(nlSvc),
		gatherer.WithHost(hostSvc),
	}

	gathererSvc := gatherer.NewSvc(p, opts...)

	return &NICUpdaterApp{
		printer:       p,
		gathererSvc:   gathererSvc,
		crdBuilderSvc: crdBuilderSvc,
		crdSaverSvc:   crdSaverSvc,
	}, COKRetCode
}

func (s *NICUpdaterApp) Run() int {
	gatherSetters := []func(inventory *inventory.Inventory) error{
		s.gathererSvc.SetDMI,
		s.gathererSvc.SetNICs,
		s.gathererSvc.SetLLDPFrames,
		s.gathererSvc.SetNDPFrames,
		s.gathererSvc.SetHost,
	}

	inv := s.gathererSvc.GatherInOrder(gatherSetters)

	buildSetters := []func(*metalv1alpha4.Inventory, *inventory.Inventory){
		s.crdBuilderSvc.SetSystem,
		s.crdBuilderSvc.SetNICs,
	}

	cr, err := s.crdBuilderSvc.BuildInOrder(inv, buildSetters)
	if err != nil {
		s.printer.Err(errors.Wrap(err, "unable to build inventory resource"))
		return CErrRetCode
	}

	patch := struct {
		Spec struct {
			Nics []metalv1alpha4.NICSpec `json:"nics"`
		} `json:"spec"`
	}{}
	patch.Spec.Nics = cr.Spec.NICs

	err = s.crdSaverSvc.Patch(cr.Name, patch)
	if err != nil {
		s.printer.Err(errors.Wrap(err, "unable to patch"))
		return CErrRetCode
	}

	return COKRetCode
}
