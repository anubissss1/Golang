package task7

import (
	"fmt"
	"testing"
)

const expected = N * (N - 1) / 2

func TestPassN(t *testing.T) {
	for _, n := range []int{0, 1, 100, 10000} {
		t.Run(fmt.Sprintf("buffer_%d", n), func(t *testing.T) {
			if got := PassN(n); got != expected {
				t.Fatalf("got %d want %d", got, expected)
			}
		})
	}
}
func TestPassBatches(t *testing.T) {
	if got := PassBatches(); got != expected {
		t.Fatalf("got %d want %d", got, expected)
	}
}

var sink int

func BenchmarkPassN(b *testing.B) {
	for _, n := range []int{0, 1, 100, 10000} {
		b.Run(fmt.Sprintf("buffer_%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sink = PassN(n)
			}
		})
	}
}
func BenchmarkPassBatches(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = PassBatches()
	}
}
