package app

import (
	"bytes"
	"encoding/json"

	apiv1alpha1 "github.com/onmetal/metal-api/apis/inventory/v1alpha1"
	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/crd"
	"github.com/onmetal/inventory/pkg/flags"
	"github.com/onmetal/inventory/pkg/gatherer"
	"github.com/onmetal/inventory/pkg/inventory"
	"github.com/onmetal/inventory/pkg/mlc"
	"github.com/onmetal/inventory/pkg/printer"
)

type BenchmarkApp struct {
	printer *printer.Svc

	crdBuilderSvc *crd.BuilderSvc
	crdSaverSvc   crd.SaverSvc

	gathererSvc *gatherer.Svc
}

func NewBenchmarkApp() (*BenchmarkApp, int) {
	f := flags.NewBenchmarkFlags()

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

	mlcPerfSvc := mlc.NewPerfSvc(p, f.Root)

	opts := []gatherer.Option{
		gatherer.WithMLCPerf(mlcPerfSvc),
	}

	gathererSvc := gatherer.NewSvc(p, opts...)

	return &BenchmarkApp{
		printer:       p,
		crdBuilderSvc: crdBuilderSvc,
		crdSaverSvc:   crdSaverSvc,
		gathererSvc:   gathererSvc,
	}, COKRetCode
}

func (s *BenchmarkApp) Run() int {
	gatherSetters := []func(inventory *inventory.Inventory) error{
		s.gathererSvc.SetMlcPerf,
	}

	inv := s.gathererSvc.GatherInOrder(gatherSetters)

	jsonBytes, err := json.Marshal(inv)
	if err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to marshal result to json"))
	}

	var prettifiedJsonBuf bytes.Buffer
	if err := json.Indent(&prettifiedJsonBuf, jsonBytes, "", "\t"); err != nil {
		s.printer.VErr(errors.Wrap(err, "unable to indent json"))
	}

	s.printer.VOut("Gathered data:")
	s.printer.VOut(prettifiedJsonBuf.String())

	buildSetters := []func(*apiv1alpha1.Inventory, *inventory.Inventory){
		s.crdBuilderSvc.SetMLCPerf,
	}

	cr, err := s.crdBuilderSvc.BuildInOrder(inv, buildSetters)
	if err != nil {
		s.printer.Err(errors.Wrap(err, "unable to build inventory resource"))
		return CErrRetCode
	}

	patch := struct {
		Spec struct {
			// TODO define real spec field when CRD will get perf fields
			// MLCPerf []apiv1alpha1.MLCPerfSpec `json:"mlcPerf"`
		} `json:"spec"`
	}{}
	// TODO set data to patch when CRD will get perf fields
	// patch.Spec.MLCPerf = cr.Spec.MLCPerf

	err = s.crdSaverSvc.Patch(cr.Name, patch)
	if err != nil {
		s.printer.Err(errors.Wrap(err, "unable to patch"))
		return CErrRetCode
	}

	return COKRetCode
}
