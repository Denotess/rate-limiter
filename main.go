package main

import (
	"fmt"
	"sync"
	"time"
)

type ConfigError struct {
	ErrorField string
}

func (e ConfigError) Error() string {
	return fmt.Sprint("Numbers cannot be negative or 0: ", e.ErrorField)
}

type LeakyBucket struct {
	capacity   float64
	leakRate   float64 // requests leaking per second
	waterLevel float64
	lastTime   time.Time
	mu         sync.Mutex
}

func NewLeakyBucket(capacity float64, leakRate float64) (*LeakyBucket, error) {
	if capacity <= 0 || leakRate <= 0 {

		return nil, ConfigError{ErrorField: "capacity"}
	}
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		lastTime: time.Now(),
	}, nil
}

func (c *LeakyBucket) leak() {
	currentTime := time.Now()
	elapsedTime := currentTime.Sub(c.lastTime).Seconds()
	c.waterLevel = max(0, c.waterLevel-(elapsedTime*c.leakRate))
	c.lastTime = currentTime
}

func (c *LeakyBucket) AllowRequest() bool {
	c.mu.Lock()
	defer c.mu.Unlock() // unlock when function exits
	c.leak()
	if c.waterLevel+1 < c.capacity {
		c.waterLevel++
		return true
	}
	return false
}
