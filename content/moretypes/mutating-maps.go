// +build OMIT

package main

import "fmt"

func main() {
	m := make(map[string]int)

	m["Odpowiedź"] = 42
	fmt.Println("Wartość:", m["Odpowiedź"])

	m["Odpowiedź"] = 48
	fmt.Println("Wartość:", m["Odpowiedź"])

	delete(m, "Odpowiedź")
	fmt.Println("Wartość:", m["Odpowiedź"])

	v, ok := m["Odpowiedź"]
	fmt.Println("Wartość:", v, "Obecna?", ok)
}
