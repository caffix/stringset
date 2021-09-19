package stringset

import (
	"strings"
	"testing"
)

func TestMemHas(t *testing.T) {
	set := New("test1", "test2", "test3", "test1", "test2")
	if !set.memHas("test1") {
		t.Errorf("Set missing expected value")
	}
}

func TestMemInsert(t *testing.T) {
	expected := 3
	set := New()
	defer set.Close()

	set.memInsert("test1")
	set.memInsert("test2")
	set.memInsert("test3")
	set.memInsert("test3")
	set.memInsert("test2")
	set.memInsert("test1")
	if len(set.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set.elements), expected)
	}
}

func TestMemInsertMany(t *testing.T) {
	expected := 3
	set := New()
	defer set.Close()

	set.memInsertMany("test1", "test2", "test3", "test1", "test2")
	if len(set.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set.elements), expected)
	}
}

func TestMemRemove(t *testing.T) {
	expected := 2
	set := New("test1", "test2", "test3", "test1", "test2")
	defer set.Close()

	set.memRemove("test1")
	if len(set.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set.elements), expected)
	}
}

func TestMemSlice(t *testing.T) {
	expected := 3
	set := New("test1", "test2", "test3", "test1", "test2")
	defer set.Close()

	slice := set.memSlice()
	if len(slice) != expected {
		t.Errorf("Got %d, expected %d", len(set.elements), expected)
	}
}

func TestMemLen(t *testing.T) {
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

		if l := set.memLen(); l != test.ExpectedLen {
			t.Errorf("Returned a set len of %d instead of %d", l, test.ExpectedLen)
		}
	}
}

func TestMemUnion(t *testing.T) {
	expected := 6
	set1 := New("test1", "test2", "test3", "test6")
	defer set1.Close()

	set2 := New("test1", "test2", "test3", "test4", "test5")
	defer set2.Close()

	set1.memUnion(set2)
	if len(set1.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set1.elements), expected)
	}
}

func TestMemIntersect(t *testing.T) {
	expected := 3
	set1 := New("test1", "test2", "test3", "test6")
	defer set1.Close()

	set2 := New("test1", "test2", "test3", "test4", "test5")
	defer set2.Close()

	set1.memIntersect(set2)
	if len(set1.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set1.elements), expected)
	}
}

func TestMemSubtract(t *testing.T) {
	expected := 1
	set1 := New("test1", "test2", "test3", "test6")
	defer set1.Close()

	set2 := New("test1", "test2", "test3", "test4", "test5")
	defer set2.Close()

	set1.memSubtract(set2)
	if len(set1.elements) != expected {
		t.Errorf("Got %d, expected %d", len(set1.elements), expected)
	}
}

func TestMemString(t *testing.T) {
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

			for _, s := range set.memSlice() {
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

func TestMemSet(t *testing.T) {
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

		set.memSet(test.Value)
		for _, e := range test.Expected {
			var found bool

			for _, s := range set.memSlice() {
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
