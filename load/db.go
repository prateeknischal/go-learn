package util

import (
	"math/rand"
	"sync/atomic"
)

var (
	count int32
)

// Add - adds 1 to the counter
func Add() {
	atomic.AddInt32(&count, 1)
}

// Del - subtracts 1 from the counter
func Del() {
	atomic.AddInt32(&count, -1)
}

// Get - returns the value of counter
func Get() int32 {
	// Generate a random number in [20, 100]
	v := 20 + rand.Int31n(80)
	return count + v
}
