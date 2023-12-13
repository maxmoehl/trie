package trie

import (
	"strings"
	"sync"
)

// String is a trie based on string paths delimited by a given delimiter. It is
// safe for concurrent reads and writes, although it does not guarantee that
// they are executed in a deterministic order. This can result in lost writes
// if a path is concurrently written and deleted.
type String[V any] interface {
	// Put a new key into the trie. The path is split at the delimiter.
	Put(path string, value V)
	// Get the value at a path. `found` indicates whether the node exists in
	// the trie which does not necessarily mean that the value is meaningful.
	// If you access a node that was created as part of a longer path the value
	// might be the default value of type V as it was not explicitly set.
	Get(path string) (value V, found bool)
	// Delete the node at the given path (including all of its children). If
	// the node does not exist, delete does not modify the trie. Delete does
	// not check if the intermediate nodes can be garbage collected as it
	// cannot reliably determine if a value has been set or not.
	// TODO: Would it be desirable to track which nodes have values assigned
	//  and which haven't to be able to garbage collect?
	Delete(path string)
	// Delimiter that has been specified on creation of the trie.
	Delimiter() string
}

// stringTrie is the underlying implementation of a simple string-based trie.
//
// The locks are only acquired while the children map is being read or written.
type stringTrie[V any] struct {
	lock     *sync.RWMutex
	children map[string]*stringTrie[V]

	delimiter string
	value     V
}

func New[V any](delimiter string) String[V] {
	return newStringTrie[V](delimiter)
}

func newStringTrie[V any](delimiter string) *stringTrie[V] {
	return &stringTrie[V]{
		lock:      new(sync.RWMutex),
		children:  make(map[string]*stringTrie[V]),
		delimiter: delimiter,
	}
}

func (t *stringTrie[V]) Delimiter() string {
	return t.delimiter
}

func (t *stringTrie[V]) Put(path string, value V) {
	if path == "" {
		t.value = value
		return
	}

	key, path, _ := strings.Cut(path, t.delimiter)

	t.lock.Lock()
	child, ok := t.children[key]
	if !ok {
		child = newStringTrie[V](t.delimiter)
		t.children[key] = child
	}
	t.lock.Unlock()

	child.Put(path, value)
}

func (t *stringTrie[V]) Get(path string) (value V, found bool) {
	if path == "" {
		return t.value, true
	}

	key, path, _ := strings.Cut(path, t.delimiter)

	t.lock.RLock()
	child, ok := t.children[key]
	t.lock.RUnlock()
	if !ok {
		return value, false
	}

	return child.Get(path)
}

func (t *stringTrie[V]) Delete(path string) {
	key, path, _ := strings.Cut(path, t.delimiter)

	if path == "" {
		t.lock.Lock()
		delete(t.children, key)
		t.lock.Unlock()

		return
	}

	t.lock.RLock()
	child, ok := t.children[key]
	t.lock.RUnlock()

	if !ok {
		return
	}

	child.Delete(path)
}
