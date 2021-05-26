// +build OMIT

package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	mu sync.Mutex
	v  int
}

// Inc increments the counter.
func (c *SafeCounter) Inc() {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the int c.v.
	c.v++
	c.mu.Unlock()
}

// Value returns the current value of the counter.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mu.Unlock()
	return c.v
}

func main() {
	c := SafeCounter{}
	for i := 0; i < 1000; i++ {
		go c.Inc()
	}

	time.Sleep(time.Second)
	fmt.Println(c.Value())
}

