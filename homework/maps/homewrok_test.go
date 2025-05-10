package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type OrderedMap[K cmp.Ordered, V any] struct {
	root *Node[K, V]
	size int
}

type Node[K cmp.Ordered, V any] struct {
	Key   K
	Val   V
	Left  *Node[K, V]
	Right *Node[K, V]
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{} // need to implement
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	if m.root == nil {
		m.root = &Node[K, V]{Key: key, Val: value}
		m.size++
		return
	}

	prev := m.root
	next := m.root
	for next != nil {
		prev = next
		switch {
		case key < next.Key:
			next = next.Left
		case key > next.Key:
			next = next.Right
		default:
			next.Val = value // надо ли обновлять значение?
			return
		}
	}

	if key < prev.Key {
		prev.Left = &Node[K, V]{Key: key, Val: value}
	} else {
		prev.Right = &Node[K, V]{Key: key, Val: value}
	}

	m.size++
}

func (m *OrderedMap[K, V]) Erase(key K) {
	if m.root == nil {
		return
	}
	var deleted bool
	m.root, deleted = deleteNode(m.root, key)
	if deleted {
		m.size--
	}
}

func deleteNode[K cmp.Ordered, V any](node *Node[K, V], key K) (*Node[K, V], bool) {
	if node == nil {
		return nil, false
	}

	var deleted bool
	if key < node.Key {
		node.Left, deleted = deleteNode(node.Left, key)
	} else if key > node.Key {
		node.Right, deleted = deleteNode(node.Right, key)
	} else {
		deleted = true
		if node.Left == nil && node.Right == nil {
			return nil, deleted
		}
		if node.Left == nil {
			return node.Right, deleted
		}
		if node.Right == nil {
			return node.Left, deleted
		}
		minNode := findMin(node.Right)
		node.Key = minNode.Key
		node.Val = minNode.Val
		node.Right, _ = deleteNode(node.Right, minNode.Key)
	}
	return node, deleted
}

func findMin[K cmp.Ordered, V any](node *Node[K, V]) *Node[K, V] {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	if m.root == nil {
		return false
	}

	next := m.root
	for next != nil {
		switch {
		case key < next.Key:
			next = next.Left
		case key > next.Key:
			next = next.Right
		default:
			return true
		}
	}
	return false
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	if m.root != nil {
		m.doActionOnNode(m.root, action)
	}
}

func (m *OrderedMap[K, V]) doActionOnNode(node *Node[K, V], action func(K, V)) {
	if node.Left != nil {
		m.doActionOnNode(node.Left, action)
	}
	action(node.Key, node.Val)
	if node.Right != nil {
		m.doActionOnNode(node.Right, action)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
