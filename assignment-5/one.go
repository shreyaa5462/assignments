package main

import (
	"fmt"
	"math"
)

type Shape interface {
	area() float64
}
type Rectangle struct {
	len     float64
	breadth float64
}
type Circle struct {
	radius float64
}
type Square struct {
	side float64
}

func (r *Rectangle) area() float64 {
	return r.len * r.breadth
}
func (c *Circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (s *Square) area() float64 {
	return s.side * s.side
}
func printarea(str string, s Shape) {
	fmt.Println("area of %s is %v", str, s.area())
}
func main() {
	rectangle := Rectangle{10, 20}
	square := Square{20}
	circle := Circle{7}
	printarea("rectangle", &rectangle)
	printarea("square", &square)
	printarea("circle", &circle)

}
