package lru

import (
	"sync"
)

// type Key interface {
// 	~string | ~int
// }

type Node struct { 
	Key string
	Value any
	Prev *Node
	Next *Node
}

func newNode(key string, val any) *Node {
	node := &Node{}
	node.Key = key
	node.Value = val

	return node
}

type LRU struct {
	mutex *sync.RWMutex

	capacity int
	head *Node
	tail *Node

	cache map[string]*Node
}


func NewLRU(capacity int) (lru *LRU) {
	lru = &LRU{
		capacity: capacity,
		cache: make(map[string]*Node),
		mutex: &sync.RWMutex{},
		head: &Node{},
		tail: &Node{},
	}

	head := lru.head
	tail := lru.tail

	head.Next = tail
	tail.Prev = head
	lru.head = head
	lru.tail = tail

	return
}


func (lru *LRU) Get(key string) any {
	lru.mutex.RLock()
	defer lru.mutex.RUnlock()

	if node, ok := lru.cache[key]; ok {
		lru.delete(node)
		lru.insert(node)

		return node.Value
	}

	return nil
}

func (lru *LRU) Put(key string, val any) {
	lru.mutex.Lock()
    defer lru.mutex.Unlock()

	if node, ok := lru.cache[key]; ok {
		lru.delete(node)
		delete(lru.cache, key)
	}

	node := newNode(key, val)
	lru.cache[key] = node
	lru.insert(node)

	if len(lru.cache) > lru.capacity {
		leastRecentNode := lru.tail.Prev
		
		lru.delete(leastRecentNode)
		delete(lru.cache, leastRecentNode.Key)
	}
}

// insert at front/right of Linked List
func (lru *LRU) insert(node *Node) {
	node.Next = lru.head.Next
	lru.head.Next.Prev = node


	lru.head.Next = node
	node.Prev = lru.head
}

// Remove from Linked List
func (lru *LRU) delete(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}
