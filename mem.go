// Copyright © by Jeff Foley 2021-2022. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package stringset

import (
	"fmt"
	"strings"
)

func (s *Set) memHas(element string) bool {
	_, exists := s.elements[strings.ToLower(element)]
	return exists
}

func (s *Set) memInsert(element string) {
	s.elements[strings.ToLower(element)] = nothing{}
}

func (s *Set) memInsertMany(elements ...string) {
	for _, i := range elements {
		s.memInsert(i)
	}
}

func (s *Set) memRemove(element string) {
	e := strings.ToLower(element)

	delete(s.elements, e)
}

func (s *Set) memSlice() []string {
	var i uint64

	k := make([]string, len(s.elements))
	for key := range s.elements {
		k[i] = key
		i++
	}
	return k
}

func (s *Set) memUnion(other *Set) {
	for _, item := range other.memSlice() {
		s.memInsert(item)
	}
}

func (s *Set) memLen() int {
	return len(s.elements)
}

func (s *Set) memSubtract(other *Set) {
	for _, item := range other.memSlice() {
		s.memRemove(item)
	}
}

func (s *Set) memIntersect(other *Set) {
	for _, item := range s.memSlice() {
		e := strings.ToLower(item)

		if !other.Has(e) {
			s.memRemove(e)
		}
	}
}

func (s *Set) memString() string {
	return strings.Join(s.memSlice(), ",")
}

func (s *Set) memSet(input string) error {
	if input == "" {
		return fmt.Errorf("String parsing failed")
	}

	for _, item := range strings.Split(input, ",") {
		s.memInsert(strings.TrimSpace(item))
	}
	return nil
}
