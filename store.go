package stringset

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dgraph-io/badger"
)

func (s *Set) setMemSaveState() {
	s.Lock()
	defer s.Unlock()

	if err := s.setupStore(); err != nil {
		return
	}

	s.storeInsertMany(s.memSlice()...)
	s.elements = make(map[string]nothing)
	s.memSaveState = true
}

func (s *Set) setupStore() error {
	path, err := ioutil.TempDir("", "stringset")
	if err != nil {
		return err
	}

	opts := badger.DefaultOptions(path)
	opts.EventLogging = false
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Set) storeHas(element string) bool {
	var exists bool
	item := strings.ToLower(element)

	_ = s.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(item))

		if err == nil {
			exists = true
		}

		return err
	})

	return exists
}

func (s *Set) storeInsert(element string) {
	item := []byte(strings.ToLower(element))

	_ = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(item, []byte("y"))
	})
}

func (s *Set) storeInsertMany(elements ...string) {
	for _, i := range elements {
		s.storeInsert(i)
	}
}

func (s *Set) storeRemove(element string) {
	e := strings.ToLower(element)

	_ = s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(e))
	})
}

func (s *Set) storeSlice() []string {
	var elements []string

	_ = s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			if key := item.KeyCopy(nil); key != nil {
				elements = append(elements, string(key))
			}
		}

		return nil
	})

	return elements
}

func (s *Set) storeUnion(other *Set) {
	for _, item := range other.Slice() {
		s.storeInsert(item)
	}
}

func (s *Set) storeLen() int {
	return len(s.storeSlice())
}

func (s *Set) storeSubtract(other *Set) {
	for _, item := range other.Slice() {
		s.storeRemove(item)
	}
}

func (s *Set) storeIntersect(other *Set) {
	for _, item := range s.storeSlice() {
		e := strings.ToLower(item)

		if !other.Has(e) {
			s.storeRemove(e)
		}
	}
}

func (s *Set) storeString() string {
	return strings.Join(s.storeSlice(), ",")
}

func (s *Set) storeSet(input string) error {
	if input == "" {
		return fmt.Errorf("String parsing failed")
	}

	for _, item := range strings.Split(input, ",") {
		s.storeInsert(strings.TrimSpace(item))
	}

	return nil
}
