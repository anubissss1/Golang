package main

import "testing"

var benchmarkSink float64

func BenchmarkReportFirstValue(b *testing.B) {
	r := Report{Title: "large report"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSink = r.FirstValue()
	}
}

func BenchmarkReportFirstPointer(b *testing.B) {
	r := Report{Title: "large report"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSink = r.FirstPointer()
	}
}

func BenchmarkPointFirstValue(b *testing.B) {
	p := Point{X: 3, Y: 4}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSink = p.FirstValue()
	}
}

func BenchmarkPointFirstPointer(b *testing.B) {
	p := Point{X: 3, Y: 4}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkSink = p.FirstPointer()
	}
}
