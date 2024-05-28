package tests

import (
	"github.com/Aki0x137/wind/pkg/lfu"
	"testing"
)

func TestNewNode(t *testing.T) {
	key := "testKey"
	val := "testVal"
	node := lfu.NewNode(key, val)

	if node.Key != key {
		t.Errorf("Expected key %v, but got %v", key, node.Key)
	}
	if node.Value != val {
		t.Errorf("Expected key %v, but got %v", key, node.Value)
	}

	// mode data assertion stuff
}

func TestLFU(t *testing.T) {
	lfu := lfu.NewLFU(2)

	lfu.Put("1", "one")
	if val := lfu.Get("1"); val != "one" {
		t.Errorf("Expected value one, but got %v", val)
	}

	lfu.Put("2", "two")
	if val := lfu.Get("2"); val != "two" {
		t.Errorf("Expected value two, but got %v", val)
	}

	lfu.Put("3", "three")
	if val := lfu.Get("1"); val != nil {
		t.Errorf("Expected value nil, but got %v", val)
	}
}

