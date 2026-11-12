package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const total = 6_000_000

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func countPrimes(from, to int) int {
	count := 0
	for n := from; n < to; n++ {
		if isPrime(n) {
			count++
		}
	}
	return count
}

func cpuWork(workers int) (int, time.Duration) {
	start := time.Now()
	results := make([]int, workers)
	var wg sync.WaitGroup
	chunk := total / workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			from := worker * chunk
			to := from + chunk
			if worker == workers-1 {
				to = total
			}
			results[worker] = countPrimes(from, to)
		}(w)
	}
	wg.Wait()
	sum := 0
	for _, n := range results {
		sum += n
	}
	return sum, time.Since(start)
}

func ioWork(workers int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); time.Sleep(100 * time.Millisecond) }()
	}
	wg.Wait()
	return time.Since(start)
}

func main() {
	fmt.Println("CPU cores:", runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	for _, workers := range []int{1, 2, 4} {
		count, elapsed := cpuWork(workers)
		fmt.Printf("CPU work, %d goroutine(s): %d primes in %v\n", workers, count, elapsed)
	}
	for _, workers := range []int{1, 10, 1000} {
		fmt.Printf("I/O work, %4d goroutine(s): %v\n", workers, ioWork(workers))
	}
}
