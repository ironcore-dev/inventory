// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package fake

import "github.com/onmetal/inventory/internal/provider"

type fakeCrds struct{}

func New() provider.Client {
	return &fakeCrds{}
}

func (f *fakeCrds) GenerateConfig(machineUUID string, config []byte) ([]byte, error) {
	return nil, nil
}

func (f *fakeCrds) Patch(machineUUID string, patch []byte) error {
	return nil
}

func (f *fakeCrds) Get(machineUUID, kind string) ([]byte, error) {
	return nil, nil
}
