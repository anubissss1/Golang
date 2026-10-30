package main

import "testing"

var ids = func() []int {
	a := make([]int, 200)
	for i := range a {
		a[i] = i
	}
	return a
}()
var sinkLength int

func TestRender(t *testing.T) {
	if RenderNew(ids) != RenderPool(ids) {
		t.Fatal("different output lengths")
	}
}
func BenchmarkRenderNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkLength = RenderNew(ids)
	}
}
func BenchmarkRenderPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkLength = RenderPool(ids)
	}
}
