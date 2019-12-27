// +build OMIT

package main

import "fmt"

func main() {
	i, j := 42, 2701

	p := &i         // Ukazuje na i
	fmt.Println(*p) // načítá i přes ukazatel
	*p = 21         // nastavuje i přes ukazatel
	fmt.Println(i)  // Zobrazí novou hodnotu i

	p = &j         // ukazuje na j
	*p = *p / 37   // vydělí j přes ukazatel
	fmt.Println(j) // zobrazí novou hodnotu j
}
