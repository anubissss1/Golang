package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func slowQuery() int { time.Sleep(50 * time.Millisecond); return 42 }

// Intentionally leaks: after timeout, nobody receives from the unbuffered channel.
func fetchLeaky(timeout time.Duration) (int, error) {
	ch := make(chan int)
	go func() { ch <- slowQuery() }()
	select {
	case value := <-ch:
		return value, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}

// Fixed: a one-item buffer lets the worker send even if the caller timed out.
func fetchFixed(timeout time.Duration) (int, error) {
	ch := make(chan int, 1)
	go func() { ch <- slowQuery() }()
	select {
	case value := <-ch:
		return value, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}

func startTickerLeaky() {
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
		}
	}()
}

func startTicker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):
			}
		}
	}()
}

func printGoroutineProfile() { _ = pprof.Lookup("goroutine").WriteTo(os.Stdout, 1) }

func main() {
	fmt.Println("start:", runtime.NumGoroutine())
	for i := 0; i < 5; i++ {
		_, _ = fetchLeaky(time.Millisecond)
	}
	time.Sleep(60 * time.Millisecond)
	fmt.Println("after leaky fetches:", runtime.NumGoroutine())
	for i := 0; i < 5; i++ {
		_, _ = fetchFixed(time.Millisecond)
	}
	time.Sleep(60 * time.Millisecond)
	fmt.Println("after fixed fetches:", runtime.NumGoroutine())
	ctx, cancel := context.WithCancel(context.Background())
	startTicker(ctx)
	fmt.Println("ticker started:", runtime.NumGoroutine())
	cancel()
	time.Sleep(20 * time.Millisecond)
	fmt.Println("after cancellation:", runtime.NumGoroutine())
	fmt.Println("goroutine profile:")
	printGoroutineProfile()
}
