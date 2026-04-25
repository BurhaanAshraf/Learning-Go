package main

import (
	"fmt"
	"math"
)

// ============================================================
// 🔹 1. Interface Definition
// ============================================================

// Interface defines behavior using method signatures
// Any type that implements all methods automatically
// satisfies the interface (implicit implementation)

type geometry interface {
	area() float64
	perimeter() float64
}


// ============================================================
// 🔹 2. Rectangle Implementation
// ============================================================

type Rectangle struct {
	length float64
	width  float64
}

func (r Rectangle) area() float64 {
	return r.length * r.width
}

func (r Rectangle) perimeter() float64 {
	return 2 * (r.length + r.width)
}


// ============================================================
// 🔹 3. Circle Implementation
// ============================================================

type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func (c circle) diameter() float64 {
	return 2 * c.radius
}
// diameter() is extra method
// interface does not care about extra methods


// ============================================================
// 🔹 4. Incomplete Implementation Example
// ============================================================

type rect struct {
	length float64
	width  float64
}

func (r rect) area() float64 {
	return r.length * r.width
}
// rect does NOT implement perimeter()
// so it does NOT satisfy geometry interface


// ============================================================
// 🔹 5. Function Accepting Interface
// ============================================================

func measuring(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perimeter())
}

// Any type implementing geometry can be passed here
// This is polymorphism


// ============================================================
// 🔹 6. Empty Interface + Variadic Example
// ============================================================

// interface{} means any type
// (modern Go uses "any", but interface{} is still common)

func myPrinter(i ...interface{}) {

	for _, v := range i {
		fmt.Println(v)
	}
}

// Variadic + empty interface
// allows multiple values of different types


// ============================================================
// 🔹 7. Type Switch
// ============================================================

func printType(i interface{}) {

	switch i.(type) {
	case int:
		fmt.Println("Integer")

	case bool:
		fmt.Println("Boolean")

	default:
		fmt.Println("Unknown")
	}
}

// i.(type) works only inside switch
// Used to check actual runtime type


// ============================================================
// 🔹 8. Main Function
// ============================================================

func Interface() {

	rec := Rectangle{
		length: 10,
		width:  5,
	}

	measuring(rec)


	cir := circle{
		radius: 5,
	}

	measuring(cir)


	rect := rect{
		length: 10,
		width:  20,
	}

	fmt.Println(rect.area())

	// measuring(rect)
	// ERROR:
	// rect does not implement perimeter()
	// so it does not satisfy geometry


	myPrinter(1, "John", true, 45.9)
}