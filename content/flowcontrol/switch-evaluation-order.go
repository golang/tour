// +build OMIT

package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Kdy je sobota ?")
	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		fmt.Println("Dneska.")
	case today + 1:
		fmt.Println("Zítra.")
	case today + 2:
		fmt.Println("Za dva dny.")
	default:
		fmt.Println("Příliš daleko.")
	}
}
