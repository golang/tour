// +build OMIT

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()); // ignore this detail for now
	fmt.Println("My favorite number is", rand.Intn(10))
}

