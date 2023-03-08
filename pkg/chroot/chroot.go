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
