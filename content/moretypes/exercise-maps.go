// +build OMIT

package main

import (
	"code.google.com/p/go-tour/wc"
)

func WordCount(s string) map[string]int {
	return map[string]int{"x": 1}
}

func main() {
	wc.Test(WordCount)
}
