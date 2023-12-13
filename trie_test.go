package trie_test

import (
	"testing"

	"moehl.dev/trie"
)

func TestStringSimple(t *testing.T) {
	tr := trie.New[string]("/")

	var (
		key   = "foo/bar"
		value = "baz"
	)
	tr.Put(key, value)

	gotValue, ok := tr.Get(key)
	if !ok {
		t.Errorf("expected value to be '%v' but got no value at all", value)
	}
	if value != gotValue {
		t.Errorf("expected value to be '%v' but got '%v'", value, gotValue)
	}

	tr.Delete(key)

	gotValue, ok = tr.Get(key)
	if ok {
		t.Errorf("expected to be deleted but got '%v'", gotValue)
	}
}

func TestSliceSimple(t *testing.T) {
	tr := trie.NewSlice[int, string]()

	var (
		key   = []int{1, 2, 3}
		value = "baz"
	)
	tr.Put(key, value)

	gotValue, ok := tr.Get(key)
	if !ok {
		t.Errorf("expected value to be '%v' but got no value at all", value)
	}
	if value != gotValue {
		t.Errorf("expected value to be '%v' but got '%v'", value, gotValue)
	}

	tr.Delete(key)

	gotValue, ok = tr.Get(key)
	if ok {
		t.Errorf("expected to be deleted but got '%v'", gotValue)
	}
}

func TestTrie_EmptyFinalKey(t *testing.T) {
	tr := trie.New[string]("/")

	var (
		key   = "foo/"
		value = "baz"
	)
	tr.Put(key, value)

	gotValue, ok := tr.Get(key)
	if !ok {
		t.Errorf("expected value to be '%v' but got no value at all", value)
	}
	if value != gotValue {
		t.Errorf("expected value to be '%v' but got '%v'", value, gotValue)
	}

	tr.Delete(key)

	gotValue, ok = tr.Get(key)
	if ok {
		t.Errorf("expected to be deleted but got '%v'", gotValue)
	}
}

func BenchmarkSlicePut(b *testing.B) {
	b.ReportAllocs()

	tr := trie.NewSlice[int, string]()

	key := []int{1, 2, 3}
	value := "foobar"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tr.Put(key, value)
	}
}

func BenchmarkSlicePutLongKey(b *testing.B) {
	b.ReportAllocs()

	tr := trie.NewSlice[int, string]()

	n := 100
	key := make([]int, 0, n)
	for i := 0; i < n; i++ {
		key = append(key, i)
	}
	value := "foobar"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tr.Put(key, value)
	}
}
