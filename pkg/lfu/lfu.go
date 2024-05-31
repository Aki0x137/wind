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

func newNode(key string, val any) *Node {
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
		leastFreq: 0,
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
		lfu.rebalance(node)
		return node.Value
	}

	return nil
}

func (lfu *LFU) Put(key string, val any) {
	if lfu.capacity == 0 {
		return;
	}
	lfu.mutex.Lock()
    defer lfu.mutex.Unlock()

	if node, ok := lfu.cache[key]; ok {
		node.Value = val
		lfu.rebalance(node)
	} else {
		node := newNode(key, val)

		if lfu.capacity == len(lfu.cache) {
			leastFrequentFL := lfu.freqMap[lfu.leastFreq]
			leastRecentNode := leastFrequentFL.tail.Prev
			leastFrequentFL.delete(leastRecentNode)
			delete(lfu.cache, leastRecentNode.Key)
		}

		lfu.leastFreq = 1
		lfu.updateFrequencyList(node)
		lfu.cache[key] = node
	}
}

func (lfu *LFU) rebalance(node *Node) {
	currentFL := lfu.freqMap[node.Freq]
	currentFL.delete(node)
	if lfu.leastFreq == node.Freq && !currentFL.hasNodes() {
		lfu.leastFreq++
		delete(lfu.freqMap, node.Freq)
	}

	node.Freq++
	lfu.updateFrequencyList(node)
}

func (lfu *LFU) updateFrequencyList(node *Node) {
	if fl, ok := lfu.freqMap[node.Freq]; ok {
		fl.insert(node)
	} else {
		fl = NewFreqList()
		fl.insert(node)
		lfu.freqMap[node.Freq] = fl
	}
}
