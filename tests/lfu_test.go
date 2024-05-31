package tests

import (
	"math/rand/v2"
	"strconv"
	"sync"
	"testing"

	"github.com/Aki0x137/wind/pkg/lfu"
)

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

	_ = lfu.Get("1")

	lfu.Put("3", "three")

	if val := lfu.Get("1"); val != "one" {
		t.Errorf("Expected value 'one', but got %v", val)
	}

	if val := lfu.Get("2"); val != nil {
		t.Errorf("Expected value nil for key '2', but got %v", val)
	}

	if val := lfu.Get("3"); val != "three" {
		t.Errorf("Expected value 'three' for key 'two', but got %v", val)
	}
}

// Tests for checking race cinditions
// run with go test -race
func TestLFU_Concurrent(t *testing.T) {
	t.Parallel()
	lru := lfu.NewLFU(100)
	var wg sync.WaitGroup

	values := make([]int, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := strconv.Itoa(i)
			value := rand.IntN(i + 1)
			values[i] = value
			lru.Put(key, value)
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := strconv.Itoa(i)
			lru.Get(key)
		}(i)
	}

	wg.Wait()
}
