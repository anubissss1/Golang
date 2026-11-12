package main

import "fmt"

func main() {
	fmt.Println("Primes below 200000:", CountSeq(200_000))
	for _, parts := range []int{2, 8, 1000, 200_000} {
		fmt.Printf("parts=%d: %d\n", parts, CountPar(200_000, parts))
	}
}
