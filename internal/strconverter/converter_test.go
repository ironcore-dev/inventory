// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
