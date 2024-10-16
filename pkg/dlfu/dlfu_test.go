package dlfu

import (
	"context"
	"testing"
	"time"

	"github.com/Aki0x137/wind/config"
)

func TestNewDLFUCache(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.DLFUConfig{
		Capacity:      10,
		Weight:        0.5,
		ExpiryEnabled: true,
		TrimInterval:  1 * time.Second,
	}

	cache := NewDLFUCache[int](ctx, config)

	if cache == nil {
		t.Errorf("Expected cache to be non-nil")
	}
	if cache.capacity != config.Capacity {
		t.Errorf("Expected capacity %d, got %d", config.Capacity, cache.capacity)
	}
	if cache.weight != config.Weight {
		t.Errorf("Expected weight %f, got %f", config.Weight, cache.weight)
	}
	if cache.expiryEnabled != config.ExpiryEnabled {
		t.Errorf("Expected expiryEnabled %v, got %v", config.ExpiryEnabled, cache.expiryEnabled)
	}
	if cache.cache == nil {
		t.Errorf("Expected cache to be non-nil")
	}
	if cache.increment.Load() != 1.0 {
		t.Errorf("Expected increment to be 1.0, got %f", cache.increment.Load())
	}
	if config.Weight == 0.0 {
		if cache.decay != 1.0 {
			t.Errorf("Expected decay to be 1.0, got %f", cache.decay)
		}
	} else {
		p := float64(config.Capacity) * config.Weight
		expectedDecay := (p + 1.0) / p
		if cache.decay != expectedDecay {
			t.Errorf("Expected decay %f, got %f", expectedDecay, cache.decay)
		}
	}
}
