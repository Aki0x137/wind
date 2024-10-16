package tests

import (
	"sync"
	"testing"

	"github.com/Aki0x137/wind/types"
)

func TestFloat64_Add(t *testing.T) {
	f := types.NewFloat64(0.0)

	tests := []struct {
		delta    float64
		expected float64
	}{
		{1.0, 1.0},
		{2.5, 3.5},
		{-1.5, 2.0},
		{0.0, 2.0},
	}

	for _, test := range tests {
		result := f.Add(test.delta)
		if result != test.expected {
			t.Errorf("Add(%v) = %v; want %v", test.delta, result, test.expected)
		}
	}
}

func NewFloat64(f float64) {
	panic("unimplemented")
}

func TestFloat64_Add_Concurrent(t *testing.T) {
	f := types.NewFloat64(0.0)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f.Add(1.0)
		}()
	}

	wg.Wait()

	if result := f.Load(); result != 1000.0 {
		t.Errorf("Concurrent Add() = %v; want %v", result, 1000.0)
	}
}
