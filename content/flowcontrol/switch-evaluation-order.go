// +build OMIT

package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Kiedy jest sobota?")
	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		fmt.Println("Dzisiaj.")
	case today + 1:
		fmt.Println("Jutro.")
	case today + 2:
		fmt.Println("Za dwa dni.")
	default:
		fmt.Println("Zbyt daleko!")
	}
}
