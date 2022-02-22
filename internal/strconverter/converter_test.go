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

import "testing"

func TestQuotaToIntWithMilli(t *testing.T) {
	quota := "100m"

	if QuotaToInt(quota) != 10000 {
		t.Log("incorrect output value. Should be equal to 10000000, but got ", QuotaToInt(quota))
		t.Fail()
	}
}

func TestQuotaToInt(t *testing.T) {
	quota := "1"
	if QuotaToInt(quota) != defaultQuota {
		t.Log("incorrect output value. Should be equal to 100000, got", QuotaToInt(quota))
		t.Fail()
	}
}

func TestFailQuotaToInt(t *testing.T) {
	quota := "m"
	if QuotaToInt(quota) != defaultQuota {
		t.Log("incorrect output value. Should be equal to 100000, got", QuotaToInt(quota))
		t.Fail()
	}

	quota2 := "0.1v"
	if QuotaToInt(quota2) != defaultQuota {
		t.Log("incorrect output value. Should be equal to 100000, got", QuotaToInt(quota2))
		t.Fail()
	}
}

func TestRandomString(t *testing.T) {
	t1 := RandomString(10)
	if t1 == "" {
		t.Log("incorrect output value. should be non empty and has 10 characters", t1)
		t.Fail()
	}

	t2 := RandomString(10)
	if t2 == "" {
		t.Log("incorrect output value. should be non empty and has 10 characters", t2)
		t.Fail()
	}

	if t1 == t2 {
		t.Log("should not be the same")
		t.Fail()
	}
}
