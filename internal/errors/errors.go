// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package errors

import (
	"errors"
	"fmt"
)

type Reason struct {
	Message string
	Status  ErrorReason
}

type ErrorReason string

const (
	cGroupIssue      ErrorReason = "cGroup binding is failed"
	exec             ErrorReason = "exec failed"
	notFound         ErrorReason = "object not found"
	notSupportedOS   ErrorReason = "current OS not supported"
	notSupported     ErrorReason = "current object not supported"
	disksInfoRead    ErrorReason = "disk read failed"
	jsonMarshal      ErrorReason = "json marshal failed"
	providerGet      ErrorReason = "get failed"
	providerGenerate ErrorReason = "generate config process failed"
	update           ErrorReason = "update failed"
	unknown          ErrorReason = "unknown"
)

type SchedulerError struct {
	RequestStatus Reason
}

type SchedulerStatus interface {
	Status() Reason
}

func (e *SchedulerError) Error() string { return e.RequestStatus.Message }

func (e *SchedulerError) Status() Reason { return e.RequestStatus }

func ReasonForError(err error) ErrorReason {
	if reason := SchedulerStatus(nil); errors.As(err, &reason) {
		return reason.Status().Status
	}
	return unknown
}

func Unknown(name string) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("unknown: %s", name),
			Status:  unknown,
		},
	}
}

func NotFound(name string) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("%s, not found", name),
			Status:  notFound,
		},
	}
}

func NotSupportedOS(os string) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("Your os is: %s, currently not supported", os),
			Status:  notSupportedOS,
		},
	}
}

func NotSupported(object string) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("%s, not supported", object),
			Status:  notSupported,
		},
	}
}

func CGroupError(err error) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("can't create or reuse cGroup: %s", err),
			Status:  cGroupIssue,
		},
	}
}

func ExecError(err string) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("program exec failed: %s", err),
			Status:  exec,
		},
	}
}

func DiskReadError(err error) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("program exec failed, can't get list of disks: %s", err),
			Status:  disksInfoRead,
		},
	}
}

func JSONMarshalError(err error) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("json marhsal failed: %s", err),
			Status:  jsonMarshal,
		},
	}
}

func UpdateError(code int, reason string) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("update process failed with status code: %d, reason: %s", code, reason),
			Status:  update,
		},
	}
}

func GetConfigError(reason string) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("can't retrieve data from the cluster: %s", reason),
			Status:  providerGet,
		},
	}
}

func GenerateConfigError(reason string) *SchedulerError {
	return &SchedulerError{
		RequestStatus: Reason{
			Message: fmt.Sprintf("generate process failed: %s", reason),
			Status:  providerGenerate,
		},
	}
}
