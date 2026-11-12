package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type Report struct{ Data []byte }

func buildReport() *Report {
	time.Sleep(20 * time.Millisecond)
	return &Report{Data: make([]byte, 64*1024)}
}

// Leaky: if handle times out, the worker blocks sending to an unbuffered channel.
func handleLeaky(timeout time.Duration) error {
	ch := make(chan *Report)
	go func() { ch <- buildReport() }()
	select {
	case report := <-ch:
		_ = report
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	}
}

// Fix A: buffered result channel prevents the send from blocking after timeout.
func handleBuffered(timeout time.Duration) error {
	ch := make(chan *Report, 1)
	go func() { ch <- buildReport() }()
	select {
	case report := <-ch:
		_ = report
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	}
}

// Fix B: cancellation also stops the work, not merely the blocked send.
func buildReportContext(ctx context.Context) (*Report, error) {
	select {
	case <-time.After(20 * time.Millisecond):
		return &Report{Data: make([]byte, 64*1024)}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func handleContext(parent context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	ch := make(chan *Report, 1)
	go func() {
		report, err := buildReportContext(ctx)
		if err == nil {
			ch <- report
		}
	}()
	select {
	case report := <-ch:
		_ = report
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func printStats(label string) {
	runtime.GC()
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	fmt.Printf("%s: goroutines=%d heap_alloc=%.2f MB\n", label, runtime.NumGoroutine(), float64(stats.HeapAlloc)/1024/1024)
}

func runRounds(label string, handler func() error) {
	for round := 1; round <= 3; round++ {
		for i := 0; i < 1000; i++ {
			_ = handler()
		}
		time.Sleep(30 * time.Millisecond)
		printStats(fmt.Sprintf("%s round %d", label, round))
	}
}

func main() {
	printStats("start")
	// Uncomment one run at a time to compare: the leaky version intentionally accumulates goroutines.
	// runRounds("leaky",func()error{return handleLeaky(5*time.Millisecond)})
	// runRounds("buffered",func()error{return handleBuffered(5*time.Millisecond)})
	runRounds("context", func() error { return handleContext(context.Background(), 5*time.Millisecond) })
}
