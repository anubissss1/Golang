package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func randomChoice() {
	a, b := make(chan string, 1000), make(chan string, 1000)
	for i := 0; i < 1000; i++ {
		a <- "a"
		b <- "b"
	}
	counts := map[string]int{}
	for i := 0; i < 1000; i++ {
		select {
		case v := <-a:
			counts[v]++
		case v := <-b:
			counts[v]++
		}
	}
	fmt.Println("random:", counts)
}
func priority() {
	a, b := make(chan string, 1000), make(chan string, 1000)
	for i := 0; i < 1000; i++ {
		a <- "a"
		b <- "b"
	}
	counts := map[string]int{}
	for i := 0; i < 1000; i++ {
		select {
		case v := <-a:
			counts[v]++
			continue
		default:
		}
		select {
		case v := <-a:
			counts[v]++
		case v := <-b:
			counts[v]++
		}
	}
	fmt.Println("prefer a:", counts)
}
func fetchWithTimeout(timeout time.Duration) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	result := make(chan int, 1)
	go func() {
		select {
		case <-time.After(200 * time.Millisecond):
			select {
			case result <- 42:
			case <-ctx.Done():
			}
		case <-ctx.Done():
			return
		}
	}()
	select {
	case v := <-result:
		return v, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
func merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}
func source(nums ...int) <-chan int {
	ch := make(chan int, len(nums))
	for _, v := range nums {
		ch <- v
	}
	close(ch)
	return ch
}
func nilRemovedDemo() int { // bounded observation of useless closed-channel receives
	a := make(chan int)
	close(a)
	b := make(chan int)
	done := make(chan struct{})
	count := 0
	go func() { time.Sleep(10 * time.Millisecond); b <- 9; close(done) }()
	for {
		select {
		case <-a:
			count++ // BAD: ready forever; usually thousands of useless receives
		case <-b:
			<-done
			return count
		}
	}
}
func ball() {
	done := make(chan struct{})
	ball := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case v := <-ball:
					select {
					case ball <- v + 1:
					case <-done:
						return
					}
				case <-done:
					return
				}
			}
		}()
	}
	// The starter must be a goroutine, since unbuffered sends require a receiver.
	go func() {
		select {
		case ball <- 0:
		case <-done:
		}
	}()
	time.Sleep(100 * time.Millisecond)
	// A separate receiver competes with players for the ball; stop immediately after receive.
	select {
	case v := <-ball:
		fmt.Println("ball hits:", v)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("ball: no handoff before deadline")
	}
	close(done)
	wg.Wait()
	fmt.Println("remaining goroutines:", runtime.NumGoroutine())
}
func main() {
	for i := 0; i < 3; i++ {
		randomChoice()
	}
	priority()
	v, e := fetchWithTimeout(100 * time.Millisecond)
	fmt.Println("timeout fetch:", v, e)
	v, e = fetchWithTimeout(300 * time.Millisecond)
	fmt.Println("successful fetch:", v, e)
	fmt.Print("merge: ")
	for v := range merge(source(1, 2, 3), source(100, 101, 102)) {
		fmt.Print(v, " ")
	}
	fmt.Println()
	fmt.Println("closed-case useless selects without nil:", nilRemovedDemo())
	q := make(chan int, 2)
	for i := 1; i <= 4; i++ {
		select {
		case q <- i:
			fmt.Println("queued", i)
		default:
			fmt.Println("dropped", i)
		}
	}
	ball()
}
