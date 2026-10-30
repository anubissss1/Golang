package main

import (
	"fmt"
	"testing"
)

var sinkString string
var tenIDs = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func TestJoinIDs(t *testing.T) {
	if JoinIDsSlow(tenIDs) != JoinIDsFast(tenIDs) {
		t.Fatal("results differ")
	}
	slow := testing.AllocsPerRun(1000, func() { sinkString = JoinIDsSlow(tenIDs) })
	fast := testing.AllocsPerRun(1000, func() { sinkString = JoinIDsFast(tenIDs) })
	fmt.Printf("JoinIDsSlow: %.2f allocs/run; JoinIDsFast: %.2f allocs/run\n", slow, fast)
	if fast > 1 {
		t.Errorf("fast: %v allocations, want at most 1", fast)
	}
}
func BenchmarkJoinIDsSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkString = JoinIDsSlow(tenIDs)
	}
}
func BenchmarkJoinIDsFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkString = JoinIDsFast(tenIDs)
	}
}
