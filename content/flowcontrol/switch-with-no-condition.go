// +build OMIT

package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Dobré ráno!")
	case t.Hour() < 17:
		fmt.Println("Dobré odpoledne.")
	default:
		fmt.Println("Dobrý večer.")
	}
}
