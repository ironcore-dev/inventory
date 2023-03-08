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

package pci

type DeviceType struct {
	ID   string
	Name string
}

type DeviceVendor struct {
	ID   string
	Name string
}

type DeviceSubtype struct {
	ID   string
	Name string
}

type DeviceClass struct {
	ID   string
	Name string
}

type DeviceSubclass struct {
	ID   string
	Name string
}

type DeviceProgrammingInterface struct {
	ID   string
	Name string
}

type Device struct {
	Address              string
	Vendor               *DeviceVendor
	Type                 *DeviceType
	Subvendor            *DeviceVendor
	Subtype              *DeviceSubtype
	Class                *DeviceClass
	Subclass             *DeviceSubclass
	ProgrammingInterface *DeviceProgrammingInterface
}
