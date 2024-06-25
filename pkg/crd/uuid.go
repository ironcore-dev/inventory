// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package crd

import "github.com/google/uuid"

func getUUID(namespace string, identifier string) string {
	namespaceUUID := uuid.NewMD5(uuid.UUID{}, []byte(namespace))
	newUUID := uuid.NewMD5(namespaceUUID, []byte(identifier))
	return newUUID.String()
}
