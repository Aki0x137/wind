package lfu

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

func NewNode(key string, val any) *Node {
	node := &Node{}
	node.Key = key
	node.Value = val

	return node
}

type LFU struct {
	mutex *sync.RWMutex

	capacity int
	head *Node
	tail *Node

	cache map[string]*Node
}


func (lfu *LFU) init(capaity int) {
	lfu.capacity = capaity
	lfu.cache = make(map[string]*Node)
	lfu.mutex = &sync.RWMutex{}

	head := &Node{}
	tail := &Node{}


	head.Next = tail
	tail.Prev = head

	lfu.head = head
	lfu.tail = tail
}


func (lfu *LFU) Get(key string) any {
	lfu.mutex.RLock()
	defer lfu.mutex.RUnlock()

	if lfu.cache[key] != nil {
		node := lfu.cache[key]
		lfu.delete(node)
		lfu.insert(node)

		return node.Value
	}

	return nil
}

func (lfu *LFU) Put(key string, val any) {
	if lfu.cache[key] != nil {
		node := lfu.cache[key]
		node.Value = val

		lfu.delete(node)
		lfu.insert(node)
	} else {
		node := NewNode(key, val)

		lfu.cache[key] = node
		lfu.insert(node)
	}

	if len(lfu.cache) > lfu.capacity {
		leastRecentNode := lfu.tail.Prev
		
		lfu.delete(leastRecentNode)
		delete(lfu.cache, key)
	}
}

// insert at front/right of Linked List
func (lfu *LFU) insert(node *Node) {
	node.Next = lfu.head.Next
	lfu.head.Next.Prev = node


	lfu.head.Next = node
	node.Prev = lfu.head
}

// Remove from Linked List
func (lfu *LFU) delete(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}
