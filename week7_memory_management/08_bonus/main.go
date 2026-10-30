package main

import "fmt"

//go:noinline
func sum1000() int {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i
	}
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

//go:noinline
func sum10000() int {
	nums := make([]int, 10000)
	for i := range nums {
		nums[i] = i
	}
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
func main() { fmt.Println(sum1000(), sum10000()) }
