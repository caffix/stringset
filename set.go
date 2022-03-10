// Copyright © by Jeff Foley 2021-2022. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package stringset

import (
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
)

type nothing struct{}

type Set struct {
	sync.Mutex
	elements     map[string]nothing
	memSaveState bool
	db           *badger.DB
	dbPath       string
	done         chan struct{}
}

// Deduplicate utilizes the Set type to generate a unique list of strings from the input slice.
func Deduplicate(input []string) []string {
	ss := New(input...)
	defer ss.Close()

	return ss.Slice()
}

// New returns a Set containing the values provided in the arguments.
func New(initial ...string) *Set {
	s := &Set{
		elements: make(map[string]nothing, 50),
		done:     make(chan struct{}, 2),
	}

	for _, v := range initial {
		s.Insert(v)
	}

	go s.checkMemory()
	return s
}

func (s *Set) Close() {
	s.Lock()
	defer s.Unlock()

	s.done <- struct{}{}
	s.elements = make(map[string]nothing)
	if s.memSaveState {
		s.db.Close()
		os.RemoveAll(s.dbPath)
		s.memSaveState = false
	}
}

// Has returns true if the receiver Set already contains the element string argument.
func (s *Set) Has(element string) bool {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		return s.storeHas(element)
	}
	return s.memHas(element)
}

// Insert adds the element string argument to the receiver Set.
func (s *Set) Insert(element string) {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		s.storeInsert(element)
		return
	}
	s.memInsert(element)
}

// InsertMany adds all the elements strings into the receiver Set.
func (s *Set) InsertMany(elements ...string) {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		s.storeInsertMany(elements...)
		return
	}
	s.memInsertMany(elements...)
}

// Remove will delete the element string from the receiver Set.
func (s *Set) Remove(element string) {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		s.storeRemove(element)
		return
	}
	s.memRemove(element)
}

// Slice returns a string slice that contains all the elements in the Set.
func (s *Set) Slice() []string {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		return s.storeSlice()
	}
	return s.memSlice()
}

// Union adds all the elements from the other Set argument into the receiver Set.
func (s *Set) Union(other *Set) {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		s.storeUnion(other)
		return
	}
	s.memUnion(other)
}

// Len returns the number of elements in the receiver Set.
func (s *Set) Len() int {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		return s.storeLen()
	}
	return s.memLen()
}

// Subtract removes all elements in the other Set argument from the receiver Set.
func (s *Set) Subtract(other *Set) {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		s.storeSubtract(other)
		return
	}
	s.memSubtract(other)
}

// Intersect causes the receiver Set to only contain elements also found in the
// other Set argument.
func (s *Set) Intersect(other *Set) {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		s.storeIntersect(other)
		return
	}
	s.memIntersect(other)
}

// Set implements the flag.Value interface.
func (s *Set) String() string {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		return s.storeString()
	}
	return s.memString()
}

// Set implements the flag.Value interface.
func (s *Set) Set(input string) error {
	s.Lock()
	defer s.Unlock()

	if s.memSaveState {
		return s.storeSet(input)
	}
	return s.memSet(input)
}

func (s *Set) checkMemory() {
	max := 750 * uint64(1<<20) // MB
	var m runtime.MemStats
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if state, elen := s.getStateAndLength(); !state && elen > 1000 {
				runtime.ReadMemStats(&m)
				if m.Alloc >= max {
					s.setMemSaveState()
					return
				}
			}
		case <-s.done:
			return
		}
	}
}

func (s *Set) getStateAndLength() (bool, int) {
	s.Lock()
	defer s.Unlock()

	return s.memSaveState, len(s.elements)
}
