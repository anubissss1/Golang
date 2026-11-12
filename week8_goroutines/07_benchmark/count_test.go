package main

import (
	"fmt"
	"testing"
)

func TestCountParMatchesSequential(t *testing.T) {
	const n = 200_000
	want := CountSeq(n)
	for _, parts := range []int{1, 2, 8, 1000, 200_000} {
		if got := CountPar(n, parts); got != want {
			t.Errorf("parts=%d: got %d, want %d", parts, got, want)
		}
	}
}

var resultSink int

func BenchmarkCountSeq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultSink = CountSeq(200_000)
	}
}

func BenchmarkCountPar(b *testing.B) {
	for _, parts := range []int{2, 8, 1000, 200_000} {
		b.Run(fmt.Sprintf("parts=%d", parts), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				resultSink = CountPar(200_000, parts)
			}
		})
	}
}
