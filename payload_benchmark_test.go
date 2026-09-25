package main

import "testing"

var benchmarkResult string

func BenchmarkProcessPayloadOriginal(b *testing.B) {
	data := []byte("  Hello World\r\nHello Go\r\n ")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkResult = processPayloadOriginal(data)
	}
}

func BenchmarkProcessPayloadOptimized(b *testing.B) {
	data := []byte("  Hello World\r\nHello Go\r\n ")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkResult = processPayloadOptimized(data)
	}
}
