// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package redis

type DatabaseConfig struct {
	Instances map[string]Instance `json:"INSTANCES"`
	Databases map[string]Database `json:"DATABASES"`
	Version   string              `json:"VERSION"`
}

type Database struct {
	ID        int    `json:"id"`
	Separator string `json:"separator"`
	Instance  string `json:"instance"`
}

type Instance struct {
	Hostname               string `json:"hostname"`
	Port                   int    `json:"port"`
	UnixSocketPath         string `json:"unix_socket_path"`
	PasswordPath           string `json:"password_path,omitempty"`
	ConfPath               string `json:"conf_path,omitempty"`
	PersistenceForWarmBoot string `json:"persistence_for_warm_boot,omitempty"`
}
