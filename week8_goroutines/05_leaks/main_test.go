package main

import (
	"runtime"
	"testing"
	"time"
)

func TestFetchFixedDoesNotLeaveGoroutine(t *testing.T) {
	before := runtime.NumGoroutine()
	for i := 0; i < 10; i++ {
		_, _ = fetchFixed(time.Millisecond)
	}
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
	if got := runtime.NumGoroutine(); got > before {
		t.Fatalf("possible goroutine leak: before=%d after=%d", before, got)
	}
}
