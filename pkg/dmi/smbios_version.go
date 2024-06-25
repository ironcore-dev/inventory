// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package dmi

type SMBIOSVersion struct {
	Major    int
	Minor    int
	Revision int
}

func NewSMBIOSVersion(major int, minor int, revision int) *SMBIOSVersion {
	return &SMBIOSVersion{
		Major:    major,
		Minor:    minor,
		Revision: revision,
	}
}

func (s *SMBIOSVersion) GreaterOrEqual(alt *SMBIOSVersion) bool {
	return s.Compare(alt) >= 0
}

func (s *SMBIOSVersion) Lesser(alt *SMBIOSVersion) bool {
	return s.Compare(alt) < 0
}

// Returns positive if left greater than right,
// zero if equal, negative otherwise
func (s *SMBIOSVersion) Compare(alt *SMBIOSVersion) int {
	if r := s.Major - alt.Major; r != 0 {
		return r
	}

	if r := s.Minor - alt.Minor; r != 0 {
		return r
	}

	return s.Revision - alt.Revision
}
