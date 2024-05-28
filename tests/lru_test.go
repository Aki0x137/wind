package tests

import (
	"math/rand/v2"
	"strconv"
	"sync"
	"testing"

	"github.com/Aki0x137/wind/pkg/lru"
)

func TestLRU_Concurrent(t *testing.T) {
	lru := lru.NewLRU(50)
	var wg sync.WaitGroup
	t.Parallel()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := strconv.Itoa(i)
			value := rand.IntN(i+1)
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
			value := rand.IntN(i+1)
			lru.Put(key, value)
			val := lru.Get(key)
			if val != value {
				t.Errorf("Expected value %v, but got %v", value, val)
			}
		}(i)
	}
}