// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
