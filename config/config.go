package config

import (
	"time"
)

type DLFUConfig struct {
	Weight        float64
	Capacity      int
	TrimInterval  time.Duration
	ExpiryEnabled bool
}
