package main

import (
	"strings"
	"tour/wc"
)

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	for _, f := range strings.Fields(s) {
		m[f]++
	}
	return m
}

func main() {
	wc.Test(WordCount)
}
