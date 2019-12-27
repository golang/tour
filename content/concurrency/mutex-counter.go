// +build OMIT

package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter jde bezpečně použít konkurentě
type SafeCounter struct {
	v   map[string]int
	mux sync.Mutex
}

// Inc zvýší počítadlo pro daný klíč
func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	// Zamkne pomocí Lock, takže jenom jedna gorutina může přistoupit k mapě c.v v jednu chvíli.
	c.v[key]++
	c.mux.Unlock()
}

// Value vrátí aktuální hodnotu počítadla pro daný klíč
func (c *SafeCounter) Value(key string) int {
	c.mux.Lock()
	// Zamkne, takže jen jedna gorutina může v jeden čas přistoupit k mapě c.v.
	defer c.mux.Unlock()
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
