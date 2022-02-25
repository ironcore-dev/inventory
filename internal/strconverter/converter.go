// /*
// Copyright (c) 2021 T-Systems International GmbH, SAP SE or an SAP affiliate company. All right reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package strconverter

import (
	"crypto/rand"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	milliCoresMultiplier float64 = 100
	coresMultiplier      float64 = 100000
)

const floatBitSize = 32

const defaultQuota int64 = 100000

var ServerCPUFullCapacity = int64(runtime.NumCPU()) * defaultQuota

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func QuotaToInt(quota string) int64 {
	if quota == "all" || quota == "*" {
		return ServerCPUFullCapacity
	}
	v := strings.Split(quota, "m")
	if len(v) == 0 || v[0] == "" {
		log.Printf("can't parse or not provided cpu quota: `%s`. default value: `1 CPU` is returned\n", quota)
		return defaultQuota
	}
	// we need float because `0.1` is also possible.
	f, err := strconv.ParseFloat(v[0], floatBitSize)
	if err == nil {
		if strings.Contains(quota, "m") {
			return int64(f * milliCoresMultiplier)
		}
		return int64(f * coresMultiplier)
	}
	log.Printf("can't parse or not provided cpu quota: `%s`. default value: `1 CPU` is returned\n", quota)
	return defaultQuota
}

func RandomString(length int) string {
	if r := uuid.New().String(); r != "" {
		return r[:length]
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}
