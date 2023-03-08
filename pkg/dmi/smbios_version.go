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
