package printer

import (
	"fmt"
	"os"
)

type Svc struct {
	verbose bool
}

func NewSvc(verbose bool) *Svc {
	return &Svc{
		verbose: verbose,
	}
}

func (s *Svc) VOut(msg string) {
	if s.verbose {
		fmt.Fprintln(os.Stdout, msg)
	}
}

func (s *Svc) Out(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}

func (s *Svc) VErr(err error) {
	if s.verbose {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (s *Svc) Err(err error) {
	fmt.Fprintln(os.Stderr, err)
}
