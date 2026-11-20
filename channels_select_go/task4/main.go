package main

import (
	"fmt"
	"runtime"
	"time"
)

// Self-contained equivalent of the Go Tour's tree exercise (no network dependency).
type Tree struct {
	Value       int
	Left, Right *Tree
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{Value: v}
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}
func build(k int) *Tree {
	var t *Tree
	for _, v := range []int{5, 3, 7, 1, 4, 6, 9, 2, 8, 10} {
		t = insert(t, v*k)
	}
	return t
}

// Walk closes once, in its OUTER goroutine, never inside the recursive visitor.
func Walk(t *Tree, ch chan int) {
	defer close(ch)
	var visit func(*Tree)
	visit = func(n *Tree) {
		if n == nil {
			return
		}
		visit(n.Left)
		ch <- n.Value
		visit(n.Right)
	}
	visit(t)
}
func Same(t1, t2 *Tree) bool {
	a, b := make(chan int), make(chan int)
	// Cancellation prevents Walk goroutines from blocking when values differ.
	done := make(chan struct{})
	defer close(done)
	walk := func(root *Tree, out chan int) {
		defer close(out)
		var visit func(*Tree) bool
		visit = func(n *Tree) bool {
			if n == nil {
				return true
			}
			if !visit(n.Left) {
				return false
			}
			select {
			case out <- n.Value:
			case <-done:
				return false
			}
			return visit(n.Right)
		}
		visit(root)
	}
	go walk(t1, a)
	go walk(t2, b)
	for {
		va, oka := <-a
		vb, okb := <-b
		if oka != okb {
			return false
		}
		if !oka {
			return true
		}
		if va != vb {
			return false
		}
	}
}
func main() {
	c := make(chan int)
	go Walk(build(1), c)
	fmt.Print("Walk: ")
	for v := range c {
		fmt.Print(v, " ")
	}
	fmt.Println()
	base := runtime.NumGoroutine()
	fmt.Println("Same equal:", Same(build(1), build(1)))
	fmt.Println("Same different:", Same(build(1), build(2)))
	for i := 0; i < 100; i++ {
		_ = Same(build(1), build(2))
	}
	time.Sleep(20 * time.Millisecond)
	fmt.Printf("goroutines before=%d after=%d (should be equal)\n", base, runtime.NumGoroutine())
}
