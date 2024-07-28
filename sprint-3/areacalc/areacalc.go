package areacalc

import (
	"strings"
)

const pi = 3.14159

type Shape interface {
	Area() float64
	Type() string
}

type Rectangle struct {
	a float64
	b float64

	shape string
}

func NewRectangle(a float64, b float64, shape string) *Rectangle {
	return &Rectangle{
		a:     a,
		b:     b,
		shape: shape,
	}
}

func (r Rectangle) Area() float64 {
	return r.a * r.b
}

func (r Rectangle) Type() string {
	return r.shape
}

type Circle struct {
	r float64

	shape string
}

func NewCircle(r float64, shape string) *Circle {
	return &Circle{
		r:     r,
		shape: shape,
	}
}

func (c Circle) Area() float64 {
	return pi * c.r * c.r
}

func (c Circle) Type() string {
	return c.shape
}

func AreaCalculator(figures []Shape) (string, float64) {
	var (
		messages  []string
		totalArea float64
	)

	for _, figure := range figures {
		messages = append(messages, figure.Type())
		totalArea += figure.Area()
	}

	return strings.Join(messages, "-"), totalArea
}
