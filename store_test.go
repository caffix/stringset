// Copyright © by Jeff Foley 2021-2022. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package stringset

import (
	"strings"
	"testing"
)

func TestStoreHas(t *testing.T) {
	set := New("test1", "test2", "test3", "test1", "test2")
	defer set.Close()
	set.setMemSaveState()

	if !set.storeHas("test1") {
		t.Errorf("Set missing expected value")
	}
}

func TestStoreInsert(t *testing.T) {
	expected := 3
	set := New()
	defer set.Close()
	set.setMemSaveState()

	set.storeInsert("test1")
	set.storeInsert("test2")
	set.storeInsert("test3")
	set.storeInsert("test3")
	set.storeInsert("test2")
	set.storeInsert("test1")
	if l := set.storeLen(); l != expected {
		t.Errorf("Got %d, expected %d", l, expected)
	}
}

func TestStoreInsertMany(t *testing.T) {
	expected := 3
	set := New()
	defer set.Close()
	set.setMemSaveState()
	set.storeInsertMany("test1", "test2", "test3", "test1", "test2")
	if l := set.storeLen(); l != expected {
		t.Errorf("Got %d, expected %d", l, expected)
	}
}

func TestStoreRemove(t *testing.T) {
	expected := 2
	set := New("test1", "test2", "test3", "test1", "test2")
	defer set.Close()
	set.setMemSaveState()
	set.storeRemove("test1")
	if l := set.storeLen(); l != expected {
		t.Errorf("Got %d, expected %d", l, expected)
	}
}

func TestStoreSlice(t *testing.T) {
	expected := 3
	set := New("test1", "test2", "test3", "test1", "test2")
	defer set.Close()
	set.setMemSaveState()
	slice := set.storeSlice()
	if l := len(slice); l != expected {
		t.Errorf("Got %d, expected %d", l, expected)
	}
}

func TestStoreLen(t *testing.T) {
	tests := []struct {
		Set         []string
		ExpectedLen int
	}{
		{[]string{"test1"}, 1},
		{[]string{"test1", "test2", "test3", "test1", "test2"}, 3},
		{[]string{"test1", "test1", "test1", "test1", "test1"}, 1},
		{[]string{"test1", "test2", "test3", "test4", "test5"}, 5},
	}

	for _, test := range tests {
		set := New(test.Set...)
		defer set.Close()
		set.setMemSaveState()

		if l := set.storeLen(); l != test.ExpectedLen {
			t.Errorf("Returned a set len of %d instead of %d", l, test.ExpectedLen)
		}
	}
}

func TestStoreUnion(t *testing.T) {
	expected := 6
	set1 := New("test1", "test2", "test3", "test6")
	defer set1.Close()
	set1.setMemSaveState()

	set2 := New("test1", "test2", "test3", "test4", "test5")
	defer set2.Close()
	set2.setMemSaveState()

	set1.storeUnion(set2)
	if l := set1.storeLen(); l != expected {
		t.Errorf("Got %d, expected %d", l, expected)
	}
}

func TestStoreIntersect(t *testing.T) {
	expected := 3
	set1 := New("test1", "test2", "test3", "test6")
	defer set1.Close()
	set1.setMemSaveState()

	set2 := New("test1", "test2", "test3", "test4", "test5")
	defer set2.Close()
	set2.setMemSaveState()

	set1.storeIntersect(set2)
	if l := set1.storeLen(); l != expected {
		t.Errorf("Got %d, expected %d", l, expected)
	}
}

func TestStoreSubtract(t *testing.T) {
	expected := 1
	set1 := New("test1", "test2", "test3", "test6")
	defer set1.Close()
	set1.setMemSaveState()

	set2 := New("test1", "test2", "test3", "test4", "test5")
	defer set2.Close()
	set2.setMemSaveState()

	set1.storeSubtract(set2)
	if l := set1.storeLen(); l != expected {
		t.Errorf("Got %d, expected %d", l, expected)
	}
}

func TestStoreString(t *testing.T) {
	tests := []struct {
		Set      []string
		Expected string
	}{
		{[]string{"test1"}, "test1"},
		{[]string{"test1", "test2", "test3"}, "test1,test2,test3"},
	}

	for _, test := range tests {
		set := New(test.Set...)
		defer set.Close()
		set.setMemSaveState()

		for _, e := range strings.Split(test.Expected, ",") {
			var found bool

			for _, s := range set.storeSlice() {
				if strings.EqualFold(s, e) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s was missing from the set", e)
			}
		}
	}
}

func TestStoreSet(t *testing.T) {
	tests := []struct {
		Value    string
		Expected []string
	}{
		{"", []string{}},
		{"test1", []string{"test1"}},
		{"test1,test2,test3", []string{"test1", "test2", "test3"}},
	}

	for _, test := range tests {
		set := New()
		defer set.Close()
		set.setMemSaveState()

		_ = set.storeSet(test.Value)
		for _, e := range test.Expected {
			var found bool

			for _, s := range set.storeSlice() {
				if strings.EqualFold(s, e) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s was missing from the set", e)
			}
		}
	}
}
