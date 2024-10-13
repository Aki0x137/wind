package types

import (
	"math"
	"sync/atomic"
)

type Float64 struct {
	v uint64
}

// NewFloat64 creates a new Float64 with the given initial value
func NewFloat64(val float64) *Float64 {
	return &Float64{v: math.Float64bits(val)}
}

// Load atomically loads and returns the value
func (f *Float64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64(&f.v))
}

// Store atomically stores the given value
func (f *Float64) Store(val float64) {
	atomic.StoreUint64(&f.v, math.Float64bits(val))
}

// Add atomically adds delta to the value and returns the new value
func (f *Float64) Add(delta float64) float64 {
	for {
		oldBits := atomic.LoadUint64(&f.v)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + delta)
		if atomic.CompareAndSwapUint64(&f.v, oldBits, newBits) {
			return math.Float64frombits(newBits)
		}
	}
}

// CompareAndSwap executes the compare-and-swap operation for a float64 value
func (f *Float64) CompareAndSwap(old, new float64) bool {
	return atomic.CompareAndSwapUint64(&f.v,
		math.Float64bits(old), math.Float64bits(new))
}
