package main

import "fmt"

type Counter struct{ n int }

func (c Counter) Get() int { return c.n }
func (c *Counter) Inc()    { c.n++ }

// Rectangle and method value bonus exercise.
type Rectangle struct{ Width, Height float64 }

func (r Rectangle) Area() float64 { return r.Width * r.Height }

func main() {
	var c Counter
	c.Inc()    // compiles: c is addressable; Go rewrites to (&c).Inc().
	(&c).Inc() // compiles.
	fmt.Println(c.Get())

	r := Rectangle{3, 4}
	f := r.Area // method value captures a copy of r at this point.
	r.Width = 10
	fmt.Println("Method value, then current method:", f(), r.Area()) // 12 40
}
