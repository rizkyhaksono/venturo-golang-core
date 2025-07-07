package main

import (
	"fmt"
)

type Shape interface {
	Area() float64
}

type Circle struct {
	Radius float32
}

type Rectangle struct {
	Width  float32
	Height float32
}

func (c Circle) Area() float32 {
	return 3.14 * float32(c.Radius*c.Radius)
}

func (r Rectangle) Area() float32 {
	return float32(r.Width * r.Height)
}

func main() {
	circle := Circle{Radius: 5}
	rectangle := Rectangle{Width: 10, Height: 5}

	fmt.Printf("Circle Area: %f\n", circle.Area())       // Expected: ~78.54
	fmt.Printf("Rectangle Area: %f\n", rectangle.Area()) // Expected: 50.0
}
