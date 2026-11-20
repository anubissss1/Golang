package main

import (
	"fmt"
	"time"
)

func run(size int) [3]time.Duration {
	ch := make(chan int, size)
	done := make(chan struct{})
	start := time.Now()
	var timings [3]time.Duration
	go func() {
		defer close(done)
		for i := 1; i <= 3; i++ {
			ch <- i
			timings[i-1] = time.Since(start)
		}
	}()
	time.Sleep(300 * time.Millisecond)
	for i := 0; i < 3; i++ {
		<-ch
		time.Sleep(100 * time.Millisecond)
	}
	<-done
	return timings
}
func main() {
	fmt.Println("Capacity | Send1 | Send2 | Send3")
	for _, n := range []int{0, 1, 2, 3, 10} {
		t := run(n)
		fmt.Printf("%8d | %5v | %5v | %5v\n", n, t[0].Round(10*time.Millisecond), t[1].Round(10*time.Millisecond), t[2].Round(10*time.Millisecond))
	}
	ch := make(chan string, 2)
	ch <- "a"
	ch <- "b"
	fmt.Println("capacity 2:", <-ch, <-ch)
	// cap=1: second send would block before the receives; avoid hanging main.
	one := make(chan string, 1)
	one <- "a"
	select {
	case one <- "b":
		fmt.Println("unexpected second send")
	default:
		fmt.Println("capacity 1: second send would block (buffer full)")
	}
	fmt.Println("first value:", <-one)
	c := make(chan int, 3)
	c <- 1
	c <- 2
	fmt.Println("len", len(c), "cap", cap(c))
	close(c)
	for i := 0; i < 3; i++ {
		v, ok := <-c
		fmt.Println(v, ok)
	}
	fmt.Println("A one-slot result buffer allows one send even if the receiver times out; two sends need more capacity or cancellation.")
}
