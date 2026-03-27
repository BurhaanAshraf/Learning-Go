package main

import "fmt"

func functions() {

	// ============================================================
	// 🔹 1. Basic Function Call
	// ============================================================

	a := 5
	b := 7

	sum := addition(a, b)
	fmt.Println(sum)

	// NOTE:
	// - Function syntax:
	//   func name(parameters) returnType { ... }
	// - Parameters act like local variables inside function


	// ============================================================
	// 🔹 2. Function Concepts
	// ============================================================

	// - Go supports multiple return values
	//   func example() (int, string)

	// - If nothing is returned → zero values are returned by default

	// - Naming:
	//   Uppercase → Public (exported)
	//   Lowercase → Private (package-level)


	// ============================================================
	// 🔹 3. Anonymous Functions
	// ============================================================

	// Function without a name
	// func() {
	//     fmt.Println("Hello World")
	// }()

	// Storing function in a variable
	greet := func() {
		fmt.Println("Hello World")
	}
	_ = greet

	// NOTE:
	// - Functions are first-class citizens (can be stored, passed, returned)


	// ============================================================
	// 🔹 4. Passing Function as Argument
	// ============================================================

	add := applyOperation(5, 5, addition)
	fmt.Println(add)

	// NOTE:
	// - Functions can be passed as arguments
	// - operation func(int, int) int → function type


	// ============================================================
	// 🔹 5. Returning Function (Closure)
	// ============================================================

	multiply := multiplyOperation(5) // returns a function
	answer := multiply(4)            // calls returned function

	fmt.Println(answer)

	// NOTE:
	// - This is a closure
	// - Inner function remembers outer variable (factor)
}


// ============================================================
// 🔹 6. Basic Function
// ============================================================

// Parameters are passed by value (copied)
// Changes inside function do NOT affect original variables

func addition(a int, b int) int {

	// a = 10 → only changes local copy, not original

	c := a + b
	return c
}


// ============================================================
// 🔹 7. Function as Parameter
// ============================================================

func applyOperation(a int, b int, operation func(int, int) int) int {
	return operation(a, b)
}


// ============================================================
// 🔹 8. Returning Function (Closure)
// ============================================================

func multiplyOperation(factor int) func(x int) int {

	return func(x int) int {
	// Closure: 'factor' is remembered even after function returns		
	return factor * x
	}
}