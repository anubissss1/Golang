package main

import (
	"fmt"
	"runtime"
	"time"
)

type Record struct {
	ID   int
	Data [256]byte
}

var live []*Record
var sum int
var lastConsumed *Record

//go:noinline
func consume(r *Record) int { lastConsumed = r; return r.ID + int(r.Data[0]) }
func main() {
	start := time.Now()
	live = make([]*Record, 0, 80000)
	for i := 0; i < 80000; i++ {
		live = append(live, &Record{ID: i})
	}
	for i := 0; i < 10000000; i++ {
		r := &Record{ID: i}
		sum += consume(r)
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("sum=%d GC cycles=%d HeapAlloc=%.1f MiB HeapSys=%.1f MiB TotalAlloc=%.1f MiB elapsed=%s\n", sum, m.NumGC, float64(m.HeapAlloc)/(1<<20), float64(m.HeapSys)/(1<<20), float64(m.TotalAlloc)/(1<<20), time.Since(start))
	runtime.KeepAlive(live)
}
