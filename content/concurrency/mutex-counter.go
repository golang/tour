// +build OMIT

package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter może być użyty współbierznie.
type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

// Inc zwiększa wartość licznika danego "key".
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	// Lock użyty by tylko jedna gorutyna mogła mieć dostęp do mapy c.v. w danym czasie
	c.v[key]++
	c.mu.Unlock()
}

// Value zwraca wartość aktualną licznika dla danego "key".
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	// Lock użyty by tylko jedna gorutyna mogła mieć dostęp do mapy c.v. w danym czasie
	defer c.mu.Unlock()
	return c.v[key]
}

func main() {
	c := SafeCounter{v: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c.Inc("somekey")
	}

	time.Sleep(time.Second)
	fmt.Println(c.Value("somekey"))
}
