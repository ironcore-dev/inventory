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
