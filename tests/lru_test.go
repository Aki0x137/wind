package tests

import (
	"math/rand/v2"
	"strconv"
	"sync"
	"testing"

	"github.com/Aki0x137/wind/pkg/lru"
)

func TestLRU(t *testing.T) {
	lru := lru.NewLRU(2)

	lru.Put("1", "one")
	if val := lru.Get("1"); val != "one" {
		t.Errorf("Expected value one, but got %v", val)
	}

	lru.Put("2", "two")
	if val := lru.Get("2"); val != "two" {
		t.Errorf("Expected value two, but got %v", val)
	}

	_ = lru.Get("1")

	lru.Put("3", "three")

	if val := lru.Get("1"); val != "one" {
		t.Errorf("Expected value 'one', but got %v", val)
	}

	if val := lru.Get("2"); val != nil {
		t.Errorf("Expected value 'nil' for key '2', but got %v", val)
	}

	if val := lru.Get("3"); val != "three" {
		t.Errorf("Expected value 'three' for key 'two', but got %v", val)
	}
}

func TestLRU_Concurrent(t *testing.T) {
	lru := lru.NewLRU(50)
	var wg sync.WaitGroup
	t.Parallel()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := strconv.Itoa(i)
			value := rand.IntN(i + 1)
			lru.Put(key, value)
			val := lru.Get(key)
			if val != value {
				t.Errorf("Expected value %v, but got %v", value, val)
			}
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := strconv.Itoa(i)
			value := rand.IntN(i + 1)
			lru.Put(key, value)
			val := lru.Get(key)
			if val != value {
				t.Errorf("Expected value %v, but got %v", value, val)
			}
		}(i)
	}

	wg.Wait()
}
