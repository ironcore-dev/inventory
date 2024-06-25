// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package output

import (
	"github.com/google/uuid"

	bencherr "github.com/onmetal/inventory/internal/errors"
)

type JobResultsQueue chan []Result

type JobResultQueue chan Result

type Result struct {
	Name, BenchmarkName, OutputSelector string
	Error                               *bencherr.SchedulerError
	Message                             []byte
	UUID                                uuid.UUID
}
