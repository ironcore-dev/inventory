// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package chroot

import (
	"io"
	"os"
	"syscall"

	"github.com/pkg/errors"
)

const (
	CRootPath = "/"
)

type Chroot struct {
	initRoot    *os.File
	initWorkDir string
}

func (c *Chroot) Close() error {
	defer c.initRoot.Close()
	if err := c.initRoot.Chdir(); err != nil {
		return errors.Wrapf(err, "unable to chdir to initial root dir")
	}

	if err := syscall.Chroot("."); err != nil {
		return errors.Wrapf(err, "unable to chroot to initial root dir")
	}

	if err := syscall.Chdir(c.initWorkDir); err != nil {
		return errors.Wrapf(err, "unable to chdir to initial working dir")
	}

	return nil
}

func New(thePath string) (io.Closer, error) {
	wd, err := os.Getwd()
	// If working directory is not set,
	// setting current root as a fallback
	if err != nil {
		wd = CRootPath
	}

	root, err := os.Open(CRootPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open %s", CRootPath)
	}

	if err := syscall.Chroot(thePath); err != nil {
		root.Close()
		return nil, errors.Wrapf(err, "unable to chroot to %s", thePath)
	}

	return &Chroot{
		initRoot:    root,
		initWorkDir: wd,
	}, nil
}
