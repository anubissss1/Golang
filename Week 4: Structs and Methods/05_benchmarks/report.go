package main

type Report struct {
	Title  string
	Values [1024]float64 // the array alone occupies 8 KiB
}

//go:noinline
func (r Report) FirstValue() float64 {
	return r.Values[0]
}

//go:noinline
func (r *Report) FirstPointer() float64 {
	return r.Values[0]
}

type Point struct {
	X, Y float64
}

//go:noinline
func (p Point) FirstValue() float64 {
	return p.X
}

//go:noinline
func (p *Point) FirstPointer() float64 {
	return p.X
}
