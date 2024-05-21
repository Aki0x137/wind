package lfu

import "golang.org/x/tools/go/analysis/passes/nilfunc"


type Key interface {
	~string | ~int
}

type Node[T Key] struct {
	Key T
	Value any
	Prev *Node[T]
	Next *Node[T]
}

type LFU[T Key] struct {
	capacity int
	Cache map[T]any
	Head *Node[T]
	Tail *Node[T]
}


func (lfu *LFU[T]) init(capaity int) {
	lfu.capacity = capaity
	lfu.Cache = make(map[T]any)

	lfu.Head = &Node[T]{}
	lfu.Tail = &Node[T]{}

	lfu.Head.Prev = nil
	lfu.Tail.Next = nil

	lfu.Head.Next = lfu.Tail
	lfu.Tail.Prev = lfu.Head
}


