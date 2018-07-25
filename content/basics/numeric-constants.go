// +build OMIT

package main

import "fmt"

const (
	// Vytvoř obrovské číslo posouváním 1 bitu doleva o 100 míst.
	// Jinými slovy,, binární číslo 1, následované 100 nulami.
	Big = 1 << 100
	// Posuň ho zase zpátky o 99 míst, takže skončímeme s 1<<1, nebo 2.
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
