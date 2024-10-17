// Copyright © by Jeff Foley 2017-2024. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package stringset

import (
	"fmt"
	"strings"
	"sync"
)

type nothing struct{}

type Set struct {
	sync.Mutex
	elements map[string]nothing
}

// Deduplicate utilizes the Set type to generate a unique list of strings from the input slice.
func Deduplicate(input []string) []string {
	ss := New(input...)
	defer ss.Close()

	return ss.Slice()
}

// New returns a Set containing the values provided in the arguments.
func New(initial ...string) *Set {
	s := &Set{elements: make(map[string]nothing, 50)}

	for _, v := range initial {
		s.Insert(v)
	}
	return s
}

func (s *Set) Close() {
	s.Lock()
	defer s.Unlock()

	s.elements = make(map[string]nothing)
}

// Has returns true if the receiver Set already contains the element string argument.
func (s *Set) Has(element string) bool {
	s.Lock()
	defer s.Unlock()

	_, exists := s.elements[strings.ToLower(element)]
	return exists
}

// Insert adds the element string argument to the receiver Set.
func (s *Set) Insert(element string) {
	s.Lock()
	defer s.Unlock()

	s.elements[strings.ToLower(element)] = nothing{}
}

// InsertMany adds all the elements strings into the receiver Set.
func (s *Set) InsertMany(elements ...string) {
	for _, i := range elements {
		s.Insert(i)
	}
}

// Remove will delete the element string from the receiver Set.
func (s *Set) Remove(element string) {
	s.Lock()
	defer s.Unlock()

	e := strings.ToLower(element)
	delete(s.elements, e)
}

// Slice returns a string slice that contains all the elements in the Set.
func (s *Set) Slice() []string {
	s.Lock()
	defer s.Unlock()

	var i uint64
	k := make([]string, len(s.elements))
	for key := range s.elements {
		k[i] = key
		i++
	}
	return k
}

// Union adds all the elements from the other Set argument into the receiver Set.
func (s *Set) Union(other *Set) {
	for _, item := range other.Slice() {
		s.Insert(item)
	}
}

// Len returns the number of elements in the receiver Set.
func (s *Set) Len() int {
	s.Lock()
	defer s.Unlock()

	return len(s.elements)
}

// Subtract removes all elements in the other Set argument from the receiver Set.
func (s *Set) Subtract(other *Set) {
	for _, item := range other.Slice() {
		s.Remove(item)
	}
}

// Intersect causes the receiver Set to only contain elements also found in the
// other Set argument.
func (s *Set) Intersect(other *Set) {
	for _, item := range s.Slice() {
		e := strings.ToLower(item)

		if !other.Has(e) {
			s.Remove(e)
		}
	}
}

// Set implements the flag.Value interface.
func (s *Set) String() string {
	return strings.Join(s.Slice(), ",")
}

// Set implements the flag.Value interface.
func (s *Set) Set(input string) error {
	if input == "" {
		return fmt.Errorf("String parsing failed")
	}

	for _, item := range strings.Split(input, ",") {
		s.Insert(strings.TrimSpace(item))
	}
	return nil
}
