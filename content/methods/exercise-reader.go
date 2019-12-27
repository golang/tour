// +build no-build OMIT

package main

import "golang.org/x/tour/reader"

type MyReader struct{}

// TODO: Přidej metodu Read([]byte) (int, error) do MyReader.

func main() {
	reader.Validate(MyReader{})
}
