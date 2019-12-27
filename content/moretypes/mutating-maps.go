// +build OMIT

package main

import "fmt"

func main() {
	m := make(map[string]int)

	m["Odpověď"] = 42
	fmt.Println("Hodnota:", m["Odpověď"])

	m["Odpověď"] = 48
	fmt.Println("Hodnota:", m["Odpověď"])

	delete(m, "Odpověď")
	fmt.Println("Hodnota:", m["Odpověď"])

	v, ok := m["Odpověď"]
	fmt.Println("Hodnota:", v, "Existuje?", ok)
}
