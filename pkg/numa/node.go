// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package numa

import (
	"github.com/onmetal/inventory/pkg/mem"
)

type Node struct {
	ID        int
	CPUs      []int
	Distances []int
	Memory    *mem.Info
	Stat      *Stat
}
