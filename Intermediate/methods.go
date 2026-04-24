package main

import "fmt"

// ============================================================
// 🔹 1. Struct for Method Example
// ============================================================

type rectangle struct {
	length   float64
	breadth  float64
}


// ============================================================
// 🔹 2. Value Receiver Method
// ============================================================

// Value receiver gets a copy of the struct
// Used when method only reads data

func (r rectangle) Area() float64 {
	return r.length * r.breadth
}


// ============================================================
// 🔹 3. Pointer Receiver Method
// ============================================================

// Pointer receiver gets address of original struct
// Used when method modifies original values

func (r *rectangle) Scale(factor float64) {
	r.length = r.length * factor
	r.breadth = r.breadth * factor
}


// ============================================================
// 🔹 4. Custom Type + Methods
// ============================================================

// User-defined type
type myInt int

// Methods can be attached to custom types too

func (a myInt) isPositive() bool {
	if a < 0 {
		return false
	}
	return true
}

// No instance data used here
// Method belongs to type, but still called using value

func (myInt) welcomeMessage() string {
	return "Welcome to myInt type"
}


// ============================================================
// 🔹 5. Method Promotion via Embedding
// ============================================================

type Shape struct {
	rectangle
}


// ============================================================
// 🔹 6. Main Function
// ============================================================

func methods() {

	rec := rectangle{
		length:  10,
		breadth: 20,
	}

	fmt.Println(rec.Area())

	rec.Scale(2)
	// Pointer receiver modifies original values

	fmt.Println(rec.Area())
	// After scaling by 2:
	// length = 20
	// breadth = 40
	// Area becomes 800


	// ============================================================
	// Custom Type Method Calls
	// ============================================================

	num := myInt(10)

	fmt.Println(num.welcomeMessage())
	fmt.Println(num.isPositive())

	num = myInt(-10)

	fmt.Println(num.isPositive())


	// ============================================================
	// Embedded Method Access
	// ============================================================

	S := Shape{
		rectangle: rectangle{
			length:  10,
			breadth: 10,
		},
	}

	fmt.Println(S.Area())
	// Method promotion:
	// Since rectangle is embedded,
	// Shape can directly use rectangle methods
}


// ============================================================
// 🔹 7. Important Concepts
// ============================================================

// Methods = functions with receiver

// Syntax:
// func (receiver Type) methodName() { }

// Value Receiver:
// - gets copy
// - cannot modify original

// Pointer Receiver:
// - gets original address
// - can modify original

// Methods can be attached to:
// - structs
// - custom types

// Methods cannot be attached to:
// - built-in types directly (int, string, etc.)

// Embedded types promote methods to outer struct