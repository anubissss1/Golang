package main

import "fmt"

type Point struct{ X, Y int }

//go:noinline
func NewPointValue(x, y int) Point { return Point{x, y} }

//go:noinline
func NewPointPointer(x, y int) *Point { return &Point{x, y} }

//go:noinline
func CollectNoCap(n int) []int {
	var s []int
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}

//go:noinline
func CollectWithCap(n int) []int {
	s := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}
func main() {
	fmt.Println(NewPointValue(1, 2), NewPointPointer(1, 2), len(CollectNoCap(1000)), len(CollectWithCap(1000)))
}
