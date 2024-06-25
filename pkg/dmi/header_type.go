// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package dmi

// Starting enumeration intentionally from zero
// since header type min value is zero
const (
	CBIOSInformationHeaderType = iota
	CSystemInformationHeaderType
	CBoardInformationHeaderType
)
