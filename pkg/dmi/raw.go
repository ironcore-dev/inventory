// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package dmi

import (
	"io"
	_ "unsafe"

	"github.com/digitalocean/go-smbios/smbios"
)

type Raw struct {
	Stream     io.ReadCloser
	EntryPoint smbios.EntryPoint
}
