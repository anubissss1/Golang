package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func startWaitingGoroutines(n int) (time.Duration, uint64, uint64) {
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	stop := make(chan struct{})
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-stop }()
	}
	elapsed := time.Since(start)
	runtime.ReadMemStats(&after)
	fmt.Println("goroutines alive:", runtime.NumGoroutine())
	close(stop)
	wg.Wait()
	return elapsed, after.StackSys - before.StackSys, after.Sys - before.Sys
}

func main() {
	elapsed, stackBytes, sysBytes := startWaitingGoroutines(100_000)
	fmt.Printf("start time: %v\nstack memory increase: %.2f MB\n", elapsed, float64(stackBytes)/1024/1024)
	fmt.Printf("total runtime memory increase: %.2f MB\n", float64(sysBytes)/1024/1024)
}
