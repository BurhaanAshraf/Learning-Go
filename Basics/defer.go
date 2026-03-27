package main

import "fmt"

func Defer() {

	// ============================================================
	// 🔹 1. Defer Basics
	// ============================================================

	// defer postpones execution of a function
	// It runs AFTER the surrounding function returns

	process(10)
}

func process(i int) {

	// ============================================================
	// 🔹 2. Deferred Statements
	// ============================================================

	defer fmt.Println("Deferred i value is", i)
	// Argument 'i' is evaluated immediately (current value)

	defer fmt.Println("This is first deferred statement")
	defer fmt.Println("This is second deferred statement")
	defer fmt.Println("This is third deferred statement")

	// NOTE:
	// - Multiple defer statements execute in LIFO order
	//   (Last In, First Out)


	// ============================================================
	// 🔹 3. Normal Execution
	// ============================================================

	i++

	fmt.Println("This is normal execution statement")
	fmt.Println(i)


	// ============================================================
	// 🔹 4. Key Concepts
	// ============================================================

	// - defer executes at function end
	// - Arguments are evaluated immediately
	// - Execution happens later (at return)
	// - Order: LIFO (stack behavior)
}