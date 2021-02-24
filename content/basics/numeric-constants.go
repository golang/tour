// +build OMIT

package main

import "fmt"

const (
	// Stwórz dużą liczbę poprzez przeniesienie 1 bit w lewo o 100 miejsc.
	// Innymi słowy, liczba binarna zaczynająca się od 1 która dalej posiada 100 zer.
	Big = 1 << 100
	// Przenieś bit o 99 miejsc w prawo, otrzymamy 1<<1, lub 2. 
	Small = Big >> 99
)

func needInt(x int) int { return x*10 + 1 }
func needFloat(x float64) float64 {
	return x * 0.1
}

func main() {
	fmt.Println(needInt(Small))
	fmt.Println(needFloat(Small))
	fmt.Println(needFloat(Big))
}
