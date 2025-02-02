// Copyright © by Jeff Foley 2017-2025. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package stringset

import (
	"strings"
	"testing"
)

func TestDeduplicate(t *testing.T) {
	tests := []struct {
		Set      []string
		Expected []string
	}{
		{[]string{"dup", "dup", "dup", "test1", "test2", "test3"}, []string{"dup", "test1", "test2", "test3"}},
		{[]string{"test1", "test2", "test3"}, []string{"test1", "test2", "test3"}},
	}

	for _, test := range tests {
		set := Deduplicate(test.Set)

		if l := len(set); l != len(test.Expected) {
			t.Errorf("Returned %d elements instead of %d", l, len(test.Expected))
			continue
		}

		for _, e := range test.Expected {
			var found bool

			for _, s := range set {
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

func TestHas(t *testing.T) {
	set := New("test1", "test2", "test3", "test1", "test2")
	if !set.Has("test1") {
		t.Errorf("Set missing expected value")
	}
}

func TestInsert(t *testing.T) {
	expected := 3
	set := New()
	defer set.Close()

	set.Insert("test1")
	set.Insert("test2")
	set.Insert("test3")
	set.Insert("test3")
	set.Insert("test2")
	set.Insert("test1")
	if len(set.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set.elements), expected)
	}
}

func TestInsertMany(t *testing.T) {
	expected := 3
	set := New()
	defer set.Close()

	set.InsertMany("test1", "test2", "test3", "test1", "test2")
	if len(set.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set.elements), expected)
	}
}

func TestRemove(t *testing.T) {
	expected := 2
	set := New("test1", "test2", "test3", "test1", "test2")
	defer set.Close()

	set.Remove("test1")
	if len(set.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set.elements), expected)
	}
}

func TestSlice(t *testing.T) {
	expected := 3
	set := New("test1", "test2", "test3", "test1", "test2")
	defer set.Close()

	slice := set.Slice()
	if len(slice) != expected {
		t.Errorf("Got %d, expected %d", len(set.elements), expected)
	}
}

func TestLen(t *testing.T) {
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

		if l := set.Len(); l != test.ExpectedLen {
			t.Errorf("Returned a set len of %d instead of %d", l, test.ExpectedLen)
		}
	}
}

func TestUnion(t *testing.T) {
	expected := 6
	set1 := New("test1", "test2", "test3", "test6")
	defer set1.Close()

	set2 := New("test1", "test2", "test3", "test4", "test5")
	defer set2.Close()

	set1.Union(set2)
	if len(set1.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set1.elements), expected)
	}
}

func TestIntersect(t *testing.T) {
	expected := 3
	set1 := New("test1", "test2", "test3", "test6")
	defer set1.Close()

	set2 := New("test1", "test2", "test3", "test4", "test5")
	defer set2.Close()

	set1.Intersect(set2)
	if len(set1.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set1.elements), expected)
	}
}

func TestSubtract(t *testing.T) {
	expected := 1
	set1 := New("test1", "test2", "test3", "test6")
	defer set1.Close()

	set2 := New("test1", "test2", "test3", "test4", "test5")
	defer set2.Close()

	set1.Subtract(set2)
	if len(set1.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set1.elements), expected)
	}
}

func TestString(t *testing.T) {
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

		for _, e := range strings.Split(test.Expected, ",") {
			var found bool

			for _, s := range set.Slice() {
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

func TestSet(t *testing.T) {
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

		_ = set.Set(test.Value)
		for _, e := range test.Expected {
			var found bool

			for _, s := range set.Slice() {
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
