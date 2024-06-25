// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func ToString(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read file %s", path)
	}

	contentsString := string(contents)
	trimmedString := strings.TrimSpace(contentsString)

	return trimmedString, nil
}

func ToInt(path string) (int, error) {
	fileString, err := ToString(path)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to read string from file %s ", path)
	}

	num, err := strconv.Atoi(fileString)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to convert %s file to int", fileString)
	}

	return num, nil
}

func ToUint64(path string) (uint64, error) {
	fileString, err := ToString(path)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to read string from file %s ", path)
	}

	num, err := strconv.ParseUint(fileString, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to convert %s file to int", fileString)
	}

	return num, nil
}

func ToUint32(path string) (uint32, error) {
	fileString, err := ToString(path)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to read string from file %s ", path)
	}

	num, err := strconv.ParseUint(fileString, 10, 32)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to convert %s file to int", fileString)
	}

	return uint32(num), nil
}

func ToUint16(path string) (uint16, error) {
	fileString, err := ToString(path)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to read string from file %s ", path)
	}

	num, err := strconv.ParseUint(fileString, 10, 16)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to convert %s file to int", fileString)
	}

	return uint16(num), nil
}

func ToUint8(path string) (uint8, error) {
	fileString, err := ToString(path)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to read string from file %s ", path)
	}

	num, err := strconv.ParseUint(fileString, 10, 8)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to convert %s file to int", fileString)
	}

	return uint8(num), nil
}

func ToBool(path string) (bool, error) {
	num, err := ToInt(path)
	if err != nil {
		return false, errors.Wrapf(err, "unable to read int from file %s ", path)
	}

	return num == 1, nil
}
