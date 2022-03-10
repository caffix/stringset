// Copyright © by Jeff Foley 2021-2022. All rights reserved.
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
