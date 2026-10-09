package main

import "fmt"

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// Scale must change the original rectangle, so it uses a pointer receiver.
func (r *Rectangle) Scale(k float64) {
	r.Width *= k
	r.Height *= k
}

type Circle struct{ Radius float64 }

func (c *Circle) Area() float64 {
	return 3.141592653589793 * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * 3.141592653589793 * c.Radius
}

func main() {
	rect := Rectangle{Width: 3, Height: 4}
	rect.Scale(2)
	fmt.Printf("Scaled rectangle: %+v, area %.2f, perimeter %.2f\n", rect, rect.Area(), rect.Perimeter())

	circle := &Circle{Radius: 1}
	shapes := []Shape{rect, circle}

	for _, shape := range shapes {
		fmt.Printf("%T: area=%.2f perimeter=%.2f\n", shape, shape.Area(), shape.Perimeter())
	}

	// Circle{Radius: 1} does NOT implement Shape because Area and Perimeter
	// have pointer receivers. Use &Circle{Radius: 1} instead.
}
