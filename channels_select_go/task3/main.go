package main

import (
	"fmt"
	"os"
	"sync"
)

func collect(workers int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				out <- w*10 + i
			}
		}(w)
	}
	go func() { wg.Wait(); close(out) }()
	return out
}
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}
func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}
func sum(in <-chan int) int {
	total := 0
	for n := range in {
		total += n
	}
	return total
}
func wrong() { // optional: demonstrates why producers must NOT close shared output
	out := make(chan int)
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				out <- w*10 + i
			}
			close(out)
		}(w)
	}
	for v := range out {
		fmt.Println(v)
	}
	wg.Wait()
}
func main() {
	if len(os.Args) > 1 && os.Args[1] == "wrong" {
		wrong()
		return
	}
	total := 0
	for v := range collect(3) {
		total += v
	}
	fmt.Println("collect sum:", total, "(expected 189)")
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("sum of squares:", sum(square(generator(nums...))), "(expected 385)")
	fmt.Println("sum of even squares:", sum(filterEven(square(generator(nums...)))), "(expected 220)")
	fmt.Println("No change to generator/square/sum was needed to insert filterEven.")
}
