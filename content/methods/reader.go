// +build OMIT

package main

import (
	"fmt"
	"io"
	"strings"
	"errors"
)

func main() {
	r := strings.NewReader("Hello, Reader!")

	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if errors.Is(err, io.EOF) {
			break
		}
	}
}
