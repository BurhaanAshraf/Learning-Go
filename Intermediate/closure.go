package main

import "fmt"

func closure() {

	// ============================================================
	// Closure Example 1
	// ============================================================

	subtractor := func() func(int) int {

		countdown := 99

		// The returned function forms a closure.
		// It captures the variable "countdown" from the outer function.
		// Even after the outer function finishes, the inner function
		// still has access to this variable.

		return func(x int) int {
			countdown -= x
			return countdown
		}

	}()

	fmt.Println(subtractor(1))
	fmt.Println(subtractor(5))

	// countdown persists between calls.
	// The closure keeps updating the same captured variable.


	// ============================================================
	// Closure Example 2
	// ============================================================

	sequence := adder()

	fmt.Println(sequence())
	fmt.Println(sequence())
	fmt.Println(sequence())
	fmt.Println(sequence())

	sequence2 := adder()
	fmt.Println(sequence2())
	fmt.Println(sequence2())


	// Each call to adder() creates a NEW closure with its own state.
	// sequence and sequence2 maintain independent values of i.

}

func adder() func() int {

	i := 0

	fmt.Println("Previous Value of i is", i)

	return func() int {

		i++

		fmt.Println("We are adding 1 to i")

		return i
	}
}


// ============================================================
// Closure Definition
// ============================================================

// A closure is a function that captures and remembers variables
// from its surrounding scope even after that scope has finished executing.

// In Go, functions are first-class values, meaning they can be:
// stored in variables, passed as arguments, and returned from functions.


// ============================================================
// Lexical Scoping
// ============================================================

// Closures work because of lexical scoping.
// A function can access variables defined in the scope where it was created.


// ============================================================
// Key Properties of Closures
// ============================================================

// Closures capture variables from outer functions.
// They preserve state between multiple calls.
// Each closure instance maintains its own independent state.


// ============================================================
// Example Flow (Subtractor)
// ============================================================

// countdown starts at 99

// subtractor(1)
// countdown becomes 98

// subtractor(5)
// countdown becomes 93

// countdown is not reset because the closure remembers it.


// ============================================================
// Real Uses of Closures
// ============================================================

// Closures are commonly used for:
// stateful counters
// middleware patterns
// function factories
// callbacks
// encapsulating internal state