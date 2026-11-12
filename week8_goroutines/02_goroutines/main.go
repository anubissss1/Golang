package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func say(msg string) {
	for i := 0; i < 3; i++ {
		fmt.Println(msg, i)
		time.Sleep(20 * time.Millisecond)
	}
}

func waitForBoth() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); say("goroutine") }()
	say("main")
	wg.Wait()
}

func atomicCounter() int64 {
	var wg sync.WaitGroup
	var counter atomic.Int64
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); counter.Add(1) }()
	}
	wg.Wait()
	return counter.Load()
}

func separateResults() int {
	const n = 1000
	var wg sync.WaitGroup
	results := make([]int, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(index int) { defer wg.Done(); results[index] = 1 }(i)
	}
	wg.Wait()
	total := 0
	for _, value := range results {
		total += value
	}
	return total
}

func recoverInsideGoroutine() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if value := recover(); value != nil {
				fmt.Println("recovered in worker:", value)
			}
		}()
		panic("example panic")
	}()
	wg.Wait()
}

func main() {
	fmt.Println("-- WaitGroup --")
	waitForBoth()
	fmt.Println("-- Atomic counter --")
	fmt.Println("counter:", atomicCounter())
	fmt.Println("-- Separate slice elements --")
	fmt.Println("total:", separateResults())
	fmt.Println("-- Recover --")
	recoverInsideGoroutine()
}
