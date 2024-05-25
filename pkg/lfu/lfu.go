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
	Freq int

	Prev *Node
	Next *Node
}

func NewNode(key string, val any) *Node {
	node := &Node{}
	node.Key = key
	node.Value = val
	node.Freq = 1

	return node
}

type FreqList struct {
	head *Node
	tail *Node
}

func NewFreqList() (fl * FreqList) {
	head := &Node{}
	tail := &Node{}
	head.Next = tail
	tail.Prev = head

	fl = &FreqList{
		head: head,
		tail:tail,
	}

	return
}

// insert at front/right of Linked List
func (fl *FreqList) insert(node *Node) {
	first := fl.head.Next

	node.Next = first
	first.Prev = node

	fl.head.Next = node
	node.Prev = fl.head
}

// Remove from Linked List
func (fl *FreqList) delete(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

func (fl *FreqList) hasNodes() bool {
	if fl.head.Next != fl.tail {
		return true
	} else {
		return false
	}
}


type LFU struct {
	mutex *sync.RWMutex

	capacity int
	leastFreq int
	freqMap map[int]*FreqList

	cache map[string]*Node
}

func NewLFU(capacity int) (lfu *LFU) {
	lfu = &LFU{
		capacity: capacity,
		leastFreq: 1,
		cache: make(map[string]*Node),
		freqMap: make(map[int]*FreqList),
		mutex: &sync.RWMutex{},
	}

	return
}

func (lfu *LFU) Get(key string) any {
	lfu.mutex.RLock()
	defer lfu.mutex.RUnlock()

	if node, ok := lfu.cache[key]; ok {
		node.Freq++
		lfu.rebalance(node)
	}

	return nil
}

func (lfu *LFU) Put(key string, val any) {
	lfu.mutex.Lock()
    defer lfu.mutex.Unlock()

	if node, ok := lfu.cache[key]; ok {
		node.Value = val
		node.Freq++
		lfu.rebalance(node)
	} else {
		if len(lfu.cache) == lfu.capacity {
			fl := lfu.freqMap[lfu.leastFreq]
			leastRecentNode := fl.tail.Prev
			delete(lfu.cache, leastRecentNode.Key)
			fl.delete(leastRecentNode)
			if !fl.hasNodes() {
				delete(lfu.freqMap, lfu.leastFreq)
			}
		}
	
		node := &Node{Key: key, Value: val}
		lfu.rebalance(node)

	}	
}

func (lfu *LFU) rebalance(node *Node) {
	freq := node.Freq

	if fl, ok := lfu.freqMap[freq]; ok {
		fl.insert(node)
	} else {
		fl = NewFreqList()
		fl.insert(node)
	}

	if prevFL, ok := lfu.freqMap[freq - 1]; ok {
		prevFL.delete(node)

		if !prevFL.hasNodes() {
			delete(lfu.freqMap, freq - 1)
	
			if lfu.leastFreq == freq - 1 {
				lfu.leastFreq++
			}
		}
	}
}
