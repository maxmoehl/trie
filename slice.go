package trie

import (
	"sync"
)

type Slice[K comparable, V any] interface {
	// Put a new key into the trie.
	Put(path []K, value V)
	// Get the value at a path. `found` indicates whether the node exists in
	// the trie which does not necessarily mean that the value is meaningful.
	// If you access a node that was created as part of a longer path the value
	// might be the default value of type V as it was not explicitly set.
	Get(path []K) (value V, found bool)
	// Delete the node at the given path (including all of its children). If
	// the node does not exist, delete does not modify the trie. Delete does
	// not check if the intermediate nodes can be garbage collected as it
	// cannot reliably determine if a value has been set or not.
	// TODO: Would it be desirable to track which nodes have values assigned
	//  and which haven't to be able to garbage collect?
	Delete(path []K)
}

type sliceTrie[K comparable, V any] struct {
	lock     *sync.RWMutex
	children map[K]*sliceTrie[K, V]

	value V
}

func NewSlice[K comparable, V any]() Slice[K, V] {
	return newSliceTrie[K, V]()
}

func newSliceTrie[K comparable, V any]() *sliceTrie[K, V] {
	return &sliceTrie[K, V]{
		lock:     new(sync.RWMutex),
		children: make(map[K]*sliceTrie[K, V]),
	}
}

func (t *sliceTrie[K, V]) Put(path []K, value V) {
	if len(path) == 0 {
		t.value = value
		return
	}

	t.lock.Lock()
	child, ok := t.children[path[0]]
	if !ok {
		child = newSliceTrie[K, V]()
		t.children[path[0]] = child
	}
	t.lock.Unlock()

	child.Put(path[1:], value)
}

func (t *sliceTrie[K, V]) Get(path []K) (value V, found bool) {
	if len(path) == 0 {
		return t.value, true
	}

	t.lock.RLock()
	child, ok := t.children[path[0]]
	t.lock.RUnlock()
	if !ok {
		return value, false
	}

	return child.Get(path[1:])
}

func (t *sliceTrie[K, V]) Delete(path []K) {
	if len(path) == 0 {
		panic("trie: cannot delete self")
	}

	if len(path) == 1 {
		t.lock.Lock()
		delete(t.children, path[0])
		t.lock.Unlock()

		return
	}

	t.lock.RLock()
	child, ok := t.children[path[0]]
	t.lock.RUnlock()

	if !ok {
		return
	}

	child.Delete(path[1:])
}
