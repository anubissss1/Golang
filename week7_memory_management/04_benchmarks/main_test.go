package main

import "testing"

var sinkPoint Point
var sinkPointer *Point
var sinkInts []int

func BenchmarkNewPointValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkPoint = NewPointValue(i, i+1)
	}
}
func BenchmarkNewPointPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkPointer = NewPointPointer(i, i+1)
	}
}
func BenchmarkCollectNoCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkInts = CollectNoCap(1000)
	}
}
func BenchmarkCollectWithCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkInts = CollectWithCap(1000)
	}
}
