// +build OMIT

package main

import "fmt"

func main() {
	s := []int{2, 3, 5, 7, 11, 13}
	fmt.Println("s ==", s)
	fmt.Println("s[1:4] ==", s[1:4])

	// missing low index implies 0
	fmt.Println("s[:3] ==", s[:3])

	// missing high index implies len(s)
	fmt.Println("s[4:] ==", s[4:])
}
