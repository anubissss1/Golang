package main

import (
	"sync"
	"testing"
)

var benchmarkSink int

func tinyWork() int { return 42 }

func BenchmarkGoroutineStartStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() { defer wg.Done(); benchmarkSink = tinyWork() }()
		wg.Wait()
	}
}

func BenchmarkFuncCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkSink = tinyWork()
	}
}
